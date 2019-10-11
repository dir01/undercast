package main

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

// App is dealing with podcast episodes CRUD API, scheduling episodes processing task and publishing resulting files as episodes once processing is finished
type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// Initialize sets up database connection and routes
func (a *App) Initialize(dbHost, dbPort, dbUser, dbPassword, dbName string) {
	a.initializeDatabase(dbHost, dbPort, dbUser, dbPassword, dbName)
	a.Router = mux.NewRouter()
	a.initializeRoutes()
}

// Run makes app serve requests
func (a *App) Run(addr string) {
	log.Fatal(http.ListenAndServe(addr, a.Router))
}

func (a *App) initializeRoutes() {
	a.Router.HandleFunc("/episodes", a.getEpisodesList).Methods("GET")
	a.Router.HandleFunc("/episodes", a.createEpisode).Methods("POST")
	a.Router.HandleFunc("/episodes/{id:[0-9]+}", a.deleteEpisode).Methods("DELETE")
}

func (a *App) getEpisodesList(w http.ResponseWriter, r *http.Request) {
	episodes, err := getEpisodesList(a.DB, 0, 10)
	if err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusOK, episodes)
}

func (a *App) createEpisode(w http.ResponseWriter, r *http.Request) {
	var e episode
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&e); err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid request payload")
		return
	}
	defer r.Body.Close()

	if err := e.createEpisode(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, err.Error())
		return
	}

	respondWithJSON(w, http.StatusCreated, e)
}

func (a *App) deleteEpisode(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, err := strconv.Atoi(vars["id"])
	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid episode id")
		return
	}

	e := episode{ID: id}
	if err := e.deleteEpisode(a.DB); err == nil {
		respondWithJSON(w, http.StatusOK, map[string]string{"result": "success"})
	} else if err.Error() == "Not found" {
		respondWithError(w, http.StatusNotFound, err.Error())
	} else {
		respondWithError(w, http.StatusInternalServerError, err.Error())
	}
}

func (a *App) initializeDatabase(host, port, user, password, dbName string) {
	connectionString :=
		fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", user, password, host, dbName)
	fmt.Println(connectionString)
	var err error
	a.DB, err = sql.Open("postgres", connectionString)
	if err != nil {
		log.Fatal(err)
	}

	const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS episodes (
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
