package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"os"
	"testing"
	"undercast/server"
)

type Torrent = server.Torrent

var a *server.App

func TestMain(m *testing.M) {
	a = &server.App{}
	setupTorrentMock(a)
	a.Initialize(os.Getenv("DB_URL"), "")

	code := m.Run()

	dropTable()
	os.Exit(code)
}

func TestCreateTorrent(t *testing.T) {
	t.Run("from source field", func(t *testing.T) {
		clearTable()
		tor := setupTorrentMock(a)

		payload := []byte(`{ "source": "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c" }`)
		response := getResponse("POST", "/api/torrents", bytes.NewBuffer(payload))

		checkResponse(t, response, http.StatusCreated,
			`{"id":1,"state":"","name":"","source":"magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c","filenames":null,"bytesCompleted":0,"bytesMissing":0}`,
		)

		if tor.id != 1 || tor.source != "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c" {
			t.Errorf("Magnet link not added to torrent client")
		}
	})

	t.Run("fails to create torrent without source", func(t *testing.T) {
		payload := []byte(`{}`)
		response := getResponse("POST", "/api/torrents", bytes.NewBuffer(payload))
		checkResponse(
			t, response, http.StatusInternalServerError,
			`{"error":"pq: new row for relation \"torrents\" violates check constraint \"require_source\""}`,
		)
	})
}

func TestListTorrents(t *testing.T) {
	t.Run("paginated queries", func(t *testing.T) {
		clearTable()

		a.Repository.CreateTorrent(&Torrent{Source: "a"})
		a.Repository.CreateTorrent(&Torrent{Source: "b"})
		a.Repository.CreateTorrent(&Torrent{Source: "c"})

		response := getResponse("GET", "/api/torrents", nil)
		checkResponse(t, response, http.StatusOK, []Torrent{
			Torrent{ID: 1, Source: "a"},
			Torrent{ID: 2, Source: "b"},
			Torrent{ID: 3, Source: "c"},
		})
	})

	t.Run("empty table results in empty array", func(t *testing.T) {
		clearTable()

		response := getResponse("GET", "/api/torrents", nil)

		checkResponse(t, response, http.StatusOK, nil)
	})
}

func TestDeleteTorrent(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		clearTable()
		id := insertTorrent()

		response := getResponse("DELETE", fmt.Sprintf("/api/torrents/%d", id), nil)
		checkResponse(t, response, http.StatusOK, `{"result":"success"}`)

		response = getResponse("GET", "/api/torrents", nil)
		checkResponse(t, response, http.StatusOK, `[]`)
	})

	t.Run("fails if no torrent", func(t *testing.T) {
		response := getResponse("DELETE", "/api/torrents/100", nil)
		checkResponse(t, response, http.StatusNotFound, `{"error":"Not found"}`)
	})

	t.Run("fails if wrong id", func(t *testing.T) {
		response := getResponse("DELETE", "/api/torrents/99999999999999999999999999999", nil)
		checkResponse(t, response, http.StatusBadRequest, `{"error":"Invalid torrent id"}`)
	})
}

func insertTorrent() int {
	var id int
	err := a.DB.QueryRow("INSERT INTO torrents (source) VALUES ($1) RETURNING id", "magnet url or something").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
