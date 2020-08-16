package undercast

import (
	"encoding/gob"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/gorilla/sessions"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Options struct {
	MongoURI       string
	MongoDbName    string
	UIDevServerURL string
	SessionSecret  string
	GlobalPassword string
}

func Bootstrap(options Options) (*Server, error) {
	db, err := getDb(options.MongoURI, options.MongoDbName)
	if err != nil {
		return nil, err
	}

	downloadsService := &downloadsService{repository: &downloadsRepository{db}}

	store := sessions.NewCookieStore([]byte(options.SessionSecret))
	gob.Register(map[string]interface{}{})

	server := &Server{
		downloadsService: downloadsService,
		uiDevServerURL:   options.UIDevServerURL,
		sessionStore:     store,
		globalPassword:   options.GlobalPassword,
	}
	return server, nil
}

type Server struct {
	downloadsService *downloadsService
	uiDevServerURL   string
	router           *mux.Router
	sessionStore     sessions.Store
	globalPassword   string
}

func (s *Server) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	s.initRoutes()
	s.router.ServeHTTP(w, r)
}

func (s *Server) ListenAndServe(addr string) error {
	s.initRoutes()
	return http.ListenAndServe(addr, s.router)
}

func (s *Server) initRoutes() {
	s.router = mux.NewRouter()
	s.router.HandleFunc("/api/downloads", s.createDownload()).Methods("POST")
	s.router.HandleFunc("/api/downloads", s.listDownloads()).Methods("GET")
	s.router.HandleFunc("/api/auth/login", s.login()).Methods("POST")
	s.router.HandleFunc("/api/auth/logout", s.logout()).Methods("POST")
	s.router.HandleFunc("/api/auth/profile", s.getProfile()).Methods("GET")
	s.router.PathPrefix("/").Handler(s.getUIHandler())
}

func (s *Server) createDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &downloadRequest{}
		if err := s.decodeRequest(req, r.Body); err != nil {
			s.respond(w, http.StatusBadRequest, nil, err)
			return
		}
		download, err := s.downloadsService.Add(r.Context(), req.Source)
		if err == nil {
			s.respond(w, http.StatusOK, download, nil)
		} else {
			s.respond(w, http.StatusBadRequest, nil, err)
		}
	}
}

func (s *Server) login() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &loginRequest{}
		if err := s.decodeRequest(req, r.Body); err != nil {
			s.respond(w, http.StatusBadRequest, nil, err)
			return
		}
		if req.Password != s.globalPassword {
			s.respond(w, http.StatusBadRequest, "", fmt.Errorf("wrong_password"))
			return
		}
		session, _ := s.sessionStore.Get(r, "auth-session")
		session.Values["profile"] = map[string]interface{}{"isActive": true}
		session.Save(r, w)
		s.respond(w, http.StatusOK, "OK", nil)
	}
}

func (s *Server) logout() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.sessionStore.Get(r, "auth-session")
		session.Values = map[interface{}]interface{}{}
		session.Save(r, w)
		s.respond(w, http.StatusOK, "OK", nil)
	}
}

func (s *Server) getProfile() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		session, _ := s.sessionStore.Get(r, "auth-session")
		profile, ok := session.Values["profile"]
		if !ok {
			s.respond(w, http.StatusNotFound, profile, fmt.Errorf("no_profile"))
			return
		}
		s.respond(w, http.StatusOK, profile, nil)
	}
}

func (s *Server) listDownloads() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		downloads, err := s.downloadsService.List(r.Context())
		s.respond(w, http.StatusOK, downloads, err)
	}
}

func (s *Server) getUIHandler() http.Handler {
	if s.uiDevServerURL == "" {
		return http.FileServer(http.Dir("./ui/build/"))
	}

	if parsed, err := url.ParseRequestURI(s.uiDevServerURL); err != nil {
		panic("Failed to parse provided uiDevServerURL " + s.uiDevServerURL)
	} else {
		proxy := httputil.NewSingleHostReverseProxy(parsed)
		return proxy
	}
}

func (s *Server) decodeRequest(req interface{}, body io.ReadCloser) error {
	decoder := json.NewDecoder(body)
	if err := decoder.Decode(req); err != nil {
		return err
	}
	defer body.Close()
	return nil
}

func (s *Server) respond(w http.ResponseWriter, status int, data interface{}, err error) {
	w.Header().Set("Content-Type", "application/json")
	var resp response
	if err != nil {
		resp.Status = "error"
		resp.Error = err.Error()
	} else {
		resp.Status = "success"
		resp.Payload = data
	}
	if bytes, err := json.Marshal(resp); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(`{"status":"error", "error":"Failed to serialize response"}`))
	} else {
		w.WriteHeader(status)
		w.Write(bytes)
	}
}

type response struct {
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
	Status  string      `json:"status"`
}
