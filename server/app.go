package server

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// TorrentState describes state of a single torrent download
type TorrentState struct {
	Name           string
	FileNames      []string
	BytesCompleted int64
	BytesMissing   int64
	Done           bool
}

// TorrentClient allows to download torrents and to subscribe on its state changes
type TorrentClient interface {
	AddTorrent(id int, magnetOrTorrentOrLink string) error
	OnTorrentChanged(func(id int, info TorrentState))
}

// App is dealing with podcast torrents CRUD API, scheduling torrents processing task and publishing resulting files as torrents once processing is finished
type App struct {
	Router  *mux.Router
	DB      *sql.DB
	Torrent TorrentClient
}

// Initialize sets up database connection and routes
func (a *App) Initialize(dbHost, dbPort, dbUser, dbPassword, dbName string) {
	log.Println("Initializing app")
	a.initializeDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run makes app serve requests
func (a *App) Run(addr string) {
	log.Println("Serving at address " + addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/api/torrents", a.getTorrentsList()).Methods("GET")
	a.Router.HandleFunc("/api/torrents", a.createTorrent()).Methods("POST")
	a.Router.HandleFunc("/api/torrents/{id:[0-9]+}", a.deleteTorrent()).Methods("DELETE")
	a.Router.PathPrefix("/").Handler(http.FileServer(http.Dir("./ui/dist/tcaster/")))
}

func (a *App) getTorrentsList() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		torrents, err := getTorrentsList(a.DB, 0, 10)
		if err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}
		respondWithJSON(w, http.StatusOK, torrents)
	}
}

func (a *App) createTorrent() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var t torrent
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&t); err != nil {
			respondWithError(w, http.StatusBadRequest, "Invalid request payload")
			return
		}
		defer r.Body.Close()

		if err := t.createTorrent(a.DB); err != nil {
			respondWithError(w, http.StatusInternalServerError, err.Error())
			return
		}

		if t.Magnet != "" {
			a.Torrent.AddTorrent(t.ID, t.Magnet)
		} else if t.URL != "" {
			a.Torrent.AddTorrent(t.ID, t.URL)
		}
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

		t := torrent{ID: id}
		if err := t.deleteTorrent(a.DB); err == nil {
			respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
		} else if err.Error() == "Not found" {
			respondWithError(w, http.StatusNotFound, err.Error())
		} else {
			respondWithError(w, http.StatusInternalServerError, err.Error())
		}
	}
}

func (a *App) initializeDatabase(host, port, user, password, dbName string) {
	connectionString :=
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)
	log.Println("Initializing DB: ", connectionString)

	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS torrents (
	id SERIAL PRIMARY KEY,
	name VARCHAR(500) NOT NULL,
	magnet VARCHAR(500),
	url VARCHAR(500),
	CONSTRAINT require_magnet_or_url CHECK (
		(case when magnet is null or length(magnet) = 0 then 0 else 1 end)
		<> 
		(case when url is null or length(url) = 0 then 0 else 1 end)
	)
)`
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func respondWithError(w http.ResponseWriter, code int, message string) {
	respondWithJSON(w, code, map[string]string{"error": message})
}

func respondWithJSON(w http.ResponseWriter, code int, payload interface{}) {
	response, _ := json.Marshal(payload)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	w.Write(response)
}
