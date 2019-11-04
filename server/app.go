package server

import (
	"database/sql"
	"encoding/json"
	"log"
	"net/http"
	"net/http/httputil"
	"net/url"
	"strconv"
	"sync"

	"github.com/gorilla/mux"
	"github.com/gorilla/websocket"
	_ "github.com/lib/pq" //
)

// TorrentState describes state of a single torrent download
type TorrentState struct {
	Name           string   `json:"name"`
	FileNames      []string `json:"filenames"`
	BytesCompleted int64    `json:"bytesCompleted"`
	BytesMissing   int64    `json:"bytesMissing"`
	Done           bool     `json:"done"`
}

// TorrentClient allows to download torrents and to subscribe on its state changes
type TorrentClient interface {
	AddTorrent(id int, magnetOrTorrentOrLink string) error
	OnTorrentChanged(func(id int, info TorrentState))
}

// App is dealing with podcast torrents CRUD API, scheduling torrents processing tasks,
// and publishing resulting files as torrents once processing is finished
type App struct {
	Router     *mux.Router
	DB         *sql.DB
	Repository *repository
	Torrent    TorrentClient

	uiDevServerURL     string
	wsConnections      []*websocket.Conn
	wsConnectionsMutex *sync.Mutex
}

// Initialize sets up database connection and routes
func (a *App) Initialize(dbURL, uiDevServerURL string) {
	a.uiDevServerURL = uiDevServerURL
	a.wsConnectionsMutex = &sync.Mutex{}

	log.Println("Initializing app")
	a.initializeDatabase(dbURL)
	a.Repository = newRepository(a.DB)
	a.uiDevServerURL = uiDevServerURL
	a.initializeRoutes()
	a.setupTorrent()
}

// Run makes app serve requests
func (a *App) Run(addr string) {
	log.Println("Serving at address " + addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router = mux.NewRouter()
	a.Router.HandleFunc("/api/torrents", a.getTorrentsList()).Methods("GET")
	a.Router.HandleFunc("/api/torrents", a.createTorrent()).Methods("POST")
	a.Router.HandleFunc("/api/torrents/{id:[0-9]+}", a.deleteTorrent()).Methods("DELETE")
	a.Router.HandleFunc("/api/ws", a.handleWebsocket())
	a.Router.PathPrefix("/").Handler(a.getUIHandler())
}

func (a *App) getTorrentsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		torrents, err := a.Repository.getTorrentsList(0, 10)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, torrents)
	}
}

func (a *App) createTorrent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t Torrent
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&t); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		t.State = "ADDED"
		if err := a.Repository.SaveTorrent(&t); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		a.Torrent.AddTorrent(t.ID, t.Source)

		respondWithJSON(w, http.StatusCreated, t)
	}
}

func (a *App) deleteTorrent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		vars := mux.Vars(r)
		id, err := strconv.Atoi(vars["id"])
		if err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid torrent id")
			return
		}

		if err := a.Repository.deleteTorrent(id); err == nil {
			respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
		} else if err.Error() == "Not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (a *App) handleWebsocket() http.HandlerFunc {
	var upgrader = websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
	}
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }
	return func(w http.ResponseWriter, r *http.Request) {
		ws, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			log.Println(err)
		}
		a.wsConnectionsMutex.Lock()
		a.wsConnections = append(a.wsConnections, ws)
		a.wsConnectionsMutex.Unlock()
	}
}

func (a *App) setupTorrent() {
	a.Torrent.OnTorrentChanged(func(id int, state TorrentState) {
		torrent, err := a.Repository.GetTorrent(id)
		if err != nil {
			log.Print("Failed to get torrent from repository", id, err)
			return
		}
		torrent.Name = state.Name
		torrent.FileNames = state.FileNames
		torrent.BytesCompleted = state.BytesCompleted
		torrent.BytesMissing = state.BytesMissing
		if state.Done {
			torrent.markAsDownloaded()
		}

		torrent.maybeSetDefaultEpisodes()
		a.Repository.SaveTorrent(torrent)
		a.dispatchWebsocketMessage(torrent)
	})
	if torrents, err := a.Repository.getUnfinisedTorrents(); err != nil {
		log.Fatal("Failed to get unfinished torrents\n", err)
	} else {
		for _, t := range torrents {
			a.Torrent.AddTorrent(t.ID, t.Source)
		}
	}
}

func (a *App) dispatchWebsocketMessage(message interface{}) {
	bytes, _ := json.Marshal(message)
	a.wsConnectionsMutex.Lock()
	i := 0
	for _, ws := range a.wsConnections {
		if err := ws.WriteMessage(1, bytes); err == nil {
			a.wsConnections[i] = ws
			i++
		}
	}
	a.wsConnections = a.wsConnections[:i]
	a.wsConnectionsMutex.Unlock()
}

func (a *App) initializeDatabase(dbURL string) {
	dbURL = dbURL + "?sslmode=disable"
	log.Println("Initializing DB: ", dbURL)
	var err error
	a.DB, err = sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal(err)
	}
}

func (a *App) getUIHandler() http.Handler {
	if a.uiDevServerURL == "" {
		return http.FileServer(http.Dir("./ui/dist/undercast/"))
	}

	if url, err := url.ParseRequestURI(a.uiDevServerURL); err != nil {
		panic("Failed to parse provided uiDevServerURL " + a.uiDevServerURL)
	} else {
		proxy := httputil.NewSingleHostReverseProxy(url)
		return proxy
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	resp, err := json.Marshal(payload)
	if err != nil {
		log.Println("Error while serializing payload:\n", err)
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(resp)
}
