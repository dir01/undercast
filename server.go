package undercast

import (
	"encoding/json"
	"github.com/gorilla/mux"
	"io"
	"net/http"
	"net/http/httputil"
	"net/url"
)

type Options struct {
	MongoURI       string
	MongoDbName    string
	UIDevServerURL string
}

func Bootstrap(options Options) (*Server, error) {
	db, err := getDb(options.MongoURI, options.MongoDbName)
	if err != nil {
		return nil, err
	}
	downloadsService := &downloadsService{repository: &downloadsRepository{db}}
	server := &Server{downloadsService: downloadsService, uiDevServerURL: options.UIDevServerURL}
	return server, nil
}

type Server struct {
	downloadsService *downloadsService
	uiDevServerURL   string
	router           *mux.Router
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
	s.router.PathPrefix("/").Handler(s.getUIHandler())
}

func (s *Server) createDownload() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		req := &downloadRequest{}
		if err := s.decodeRequest(req, r.Body); err != nil {
			s.respond(w, nil, err)
			return
		}
		download, err := s.downloadsService.Add(r.Context(), req.Source)
		s.respond(w, download, err)
	}
}

func (s *Server) getUIHandler() http.Handler {
	if s.uiDevServerURL == "" {
		return http.FileServer(http.Dir("./ui/build/"))
	}

	if url, err := url.ParseRequestURI(s.uiDevServerURL); err != nil {
		panic("Failed to parse provided uiDevServerURL " + s.uiDevServerURL)
	} else {
		proxy := httputil.NewSingleHostReverseProxy(url)
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

func (s *Server) respond(w http.ResponseWriter, data interface{}, err error) {
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
		w.Write([]byte(`{"error": "Failed to serialize response"}`))
	} else {
		w.Write(bytes)
	}
}

type response struct {
	Error   string      `json:"error"`
	Payload interface{} `json:"payload"`
	Status  string      `json:"status"`
}
