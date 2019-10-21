package main

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
	"undercast/server"
)

var a *server.App

func TestMain(m *testing.M) {
	a = &server.App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	code := m.Run()

	dropTable()
	os.Exit(code)
}

func TestListTorrents(t *testing.T) {
	t.Run("empty table results in empty array", func(t *testing.T) {
		clearTable()

		req, _ := http.NewRequest("GET", "/api/torrents", nil)
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusOK)
	})
}

func TestCreateTorrent(t *testing.T) {
	t.Run("with magnet", func(t *testing.T) {
		clearTable()
		tor := setupTorrentMock(a)

		payload := []byte(`{
		"name": "Around the world in 80 days",
		"magnet": "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c"
	}`)
		req, _ := http.NewRequest("POST", "/api/torrents", bytes.NewBuffer(payload))
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusCreated)
		checkResponseBody(t, response, `{"id":1,"name":"Around the world in 80 days","magnet":"magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c","url":""}`)
		if tor.id != 1 || tor.source != "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c" {
			t.Errorf("Magnet link not added to torrent client")
		}
	})

	t.Run("with url", func(t *testing.T) {
		clearTable()
		tor := setupTorrentMock(a)

		payload := []byte(`{
		"name": "Around the world in 80 days",
		"url": "http://legittorrents.info/download.php?id=1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c"
		}`)
		req, _ := http.NewRequest("POST", "/api/torrents", bytes.NewBuffer(payload))
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusCreated)
		checkResponseBody(t, response, `{"id":1,"name":"Around the world in 80 days","magnet":"","url":"http://legittorrents.info/download.php?id=1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c"}`)
		if tor.id != 1 || tor.source != "http://legittorrents.info/download.php?id=1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c" {
			t.Errorf("Torrent URL not added to torrent client")
		}

	})

	t.Run("fails to create without source or url", func(t *testing.T) {
		payload := []byte(`{ "name": "Around the world in 80 days" }`)
		req, _ := http.NewRequest("POST", "/api/torrents", bytes.NewBuffer(payload))
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusInternalServerError)
		checkResponseBody(t, response, `{"error":"pq: new row for relation \"torrents\" violates check constraint \"require_magnet_or_url\""}`)
	})
}

func TestDeleteTorrent(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		clearTable()
		id := insertTorrent()

		req, _ := http.NewRequest("DELETE", fmt.Sprintf("/api/torrents/%d", id), nil)
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusOK)

		req, _ = http.NewRequest("GET", "/api/torrents", nil)
		response = executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusOK)
		checkResponseBody(t, response, `[]`)
	})

	t.Run("fails if no torrent", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/torrents/100", nil)
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusNotFound)
	})

	t.Run("fails if wrong id", func(t *testing.T) {
		req, _ := http.NewRequest("DELETE", "/api/torrents/99999999999999999999999999999", nil)
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusBadRequest)
	})

}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseStatusCode(t *testing.T, response *httptest.ResponseRecorder, expectedCode int) {
	actualCode := response.Code
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n%s", expectedCode, actualCode, response.Body)
	}
}

func checkResponseBody(t *testing.T, response *httptest.ResponseRecorder, expectedBody string) {
	actualBody := response.Body.String()
	if expectedBody != actualBody {
		t.Errorf("Unexpected response body\nEXPECTED:\n%s\nACTUAL:\n%s", expectedBody, actualBody)
	}
}

func insertTorrent() int {
	var id int
	err := a.DB.QueryRow("INSERT INTO torrents (name, magnet) VALUES ($1, $2) RETURNING id", "Some name", "Some magnet").Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}

func clearTable() {
	a.DB.Exec("DELETE FROM torrents")
	a.DB.Exec("ALTER SEQUENCE torrents_id_seq RESTART WITH 1")
}

func dropTable() {
	a.DB.Exec("DROP TABLE torrents")
}

type torrentMock struct {
	id     int
	source string
}

func setupTorrentMock(a *server.App) *torrentMock {
	t := &torrentMock{}
	a.Torrent = t
	return t
}

func (t *torrentMock) AddTorrent(id int, source string) error {
	t.id = id
	t.source = source
	return nil
}

func (t *torrentMock) OnTorrentChanged(callback func(id int, state server.TorrentState)) {

}
