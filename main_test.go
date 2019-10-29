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
	setupTorrentMock(a)
	a.Initialize( os.Getenv("DB_URL"), "", )

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
	t.Run("from source field", func(t *testing.T) {
		clearTable()
		tor := setupTorrentMock(a)

		payload := []byte(`{
		"source": "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c"
	}`)
		req, _ := http.NewRequest("POST", "/api/torrents", bytes.NewBuffer(payload))
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusCreated)
		checkResponseBody(t, response, `{"id":1,"state":"","name":"","source":"magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c","filenames":null,"bytesCompleted":0,"bytesMissing":0}`)
		if tor.id != 1 || tor.source != "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c" {
			t.Errorf("Magnet link not added to torrent client")
		}
	})

	t.Run("fails to create torrent without source", func(t *testing.T) {
		payload := []byte(`{}`)
		req, _ := http.NewRequest("POST", "/api/torrents", bytes.NewBuffer(payload))
		response := executeRequest(req)

		checkResponseStatusCode(t, response, http.StatusInternalServerError)
		checkResponseBody(t, response, `{"error":"pq: new row for relation \"torrents\" violates check constraint \"require_source\""}`)
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
	err := a.DB.QueryRow("INSERT INTO torrents (source) VALUES ($1) RETURNING id", "magnet url or something").Scan(&id)
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
