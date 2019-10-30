package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
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

func getResponse(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponse(
	t *testing.T,
	response *httptest.ResponseRecorder,
	expectedCode int,
	expectedPayload interface{},
) {
	actualCode := response.Code
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n%s", expectedCode, actualCode, response.Body)
	}

	if expectedPayload == nil {
		return
	}

	var expectedBody string
	switch expectedPayload.(type) {
	case string:
		expectedBody = expectedPayload.(string)
	default:
		expectedBytes, _ := json.Marshal(expectedPayload)
		expectedBody = string(expectedBytes)
	}

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
