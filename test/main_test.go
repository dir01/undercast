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
type Episode = server.Episode

var a *server.App

func TestMain(m *testing.M) {
	a = &server.App{}
	setupTorrentMock(a)
	a.Initialize(os.Getenv("DB_URL"), "")

	code := m.Run()

	dropTables()
	os.Exit(code)
}

func TestTorrentDownload(t *testing.T) {
	id1 := insertTorrent("source 1", "DOWNLOADING")
	id2 := insertTorrent("source 2", "DOWNLOADING")
	id3 := insertTorrent("source 3", "ENCODING")

	app := &server.App{}
	tm := setupTorrentMock(app)
	app.Initialize(os.Getenv("DB_URL"), "")

	t.Run("it resumes on boot", func(t *testing.T) {
		tm.assertTorrentAdded(t, id1, "source 1")
		tm.assertTorrentAdded(t, id2, "source 2")
		tm.assertTorrentNotAdded(t, id3, "source 3")
	})

	t.Run("on each torrent client update torrent is updated in db", func(t *testing.T) {
		clearTables()
		torrent := &Torrent{Source: "foo", State: "DOWNLOADING"}
		app.Repository.SaveTorrent(torrent)

		tm.callback(torrent.ID, server.TorrentState{
			Name:           "Around the world in 80 days",
			FileNames:      []string{"Chapter 1/all.mp3", "Chapter 2/all.mp3"},
			BytesCompleted: 300,
			BytesMissing:   9000,
			Done:           false,
		})

		reloaded, err := app.Repository.GetTorrent(torrent.ID)
		if err != nil {
			t.Error(err)
		}
		assertDeepEquals(t, &Torrent{
			ID:             torrent.ID,
			State:          "DOWNLOADING",
			Source:         "foo",
			Name:           "Around the world in 80 days",
			FileNames:      []string{"Chapter 1/all.mp3", "Chapter 2/all.mp3"},
			BytesCompleted: 300,
			BytesMissing:   9000,
			Episodes: []Episode{
				Episode{
					ID:        1,
					Name:      "Chapter 1",
					FileNames: []string{"Chapter 1/all.mp3"},
				},
				Episode{
					ID:        2,
					Name:      "Chapter 2",
					FileNames: []string{"Chapter 2/all.mp3"},
				},
			},
		}, reloaded)

		tm.callback(torrent.ID, server.TorrentState{
			Name:           "Around the world in 80 days",
			FileNames:      []string{"Chapter 1/all.mp3", "Chapter 2/all.mp3"},
			BytesCompleted: 1300,
			BytesMissing:   8000,
			Done:           false,
		})

		reloaded, err = app.Repository.GetTorrent(torrent.ID)
		if err != nil {
			t.Error(err)
		}
		assertDeepEquals(t, &Torrent{
			ID:             torrent.ID,
			State:          "DOWNLOADING",
			Source:         "foo",
			Name:           "Around the world in 80 days",
			FileNames:      []string{"Chapter 1/all.mp3", "Chapter 2/all.mp3"},
			BytesCompleted: 1300,
			BytesMissing:   8000,
			Episodes: []Episode{
				Episode{
					ID:        1,
					Name:      "Chapter 1",
					FileNames: []string{"Chapter 1/all.mp3"},
				},
				Episode{
					ID:        2,
					Name:      "Chapter 2",
					FileNames: []string{"Chapter 2/all.mp3"},
				},
			},
		}, reloaded)
	})

	t.Run("when torrent client finishes download, torrent state is DOWNLOAD_COMPLETE", func(t *testing.T) {
		torrent := &Torrent{Source: "foo", State: "DOWNLOADING"}
		app.Repository.SaveTorrent(torrent)
		tm.callback(torrent.ID, server.TorrentState{
			Name:           "Around the world in 80 days",
			BytesCompleted: 9300,
			BytesMissing:   0,
			Done:           true,
		})

		reloaded, err := app.Repository.GetTorrent(torrent.ID)
		if err != nil {
			t.Error(err)
		}
		assertDeepEquals(t, &Torrent{
			ID:             torrent.ID,
			State:          "ENCODING",
			Source:         "foo",
			Name:           "Around the world in 80 days",
			BytesCompleted: 9300,
			BytesMissing:   0,
		}, reloaded)
	})

}

func TestCreateTorrent(t *testing.T) {
	t.Run("from source field", func(t *testing.T) {
		clearTables()
		tor := setupTorrentMock(a)

		payload := []byte(`{ "source": "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c" }`)
		response := getResponse("POST", "/api/torrents", bytes.NewBuffer(payload))

		checkResponse(t, response, http.StatusCreated,
			`{"id":1,"state":"DOWNLOADING","name":"","source":"magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c","filenames":null,"bytesCompleted":0,"bytesMissing":0,"episodes":null}`,
		)
		tor.assertTorrentAdded(t, 1, "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c")
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
		clearTables()

		a.Repository.SaveTorrent(&Torrent{Source: "a"})
		a.Repository.SaveTorrent(&Torrent{Source: "b"})
		a.Repository.SaveTorrent(&Torrent{Source: "c"})

		response := getResponse("GET", "/api/torrents", nil)
		checkResponse(t, response, http.StatusOK, []Torrent{
			Torrent{ID: 1, Source: "a"},
			Torrent{ID: 2, Source: "b"},
			Torrent{ID: 3, Source: "c"},
		})
	})

	t.Run("empty table results in empty array", func(t *testing.T) {
		clearTables()

		response := getResponse("GET", "/api/torrents", nil)

		checkResponse(t, response, http.StatusOK, nil)
	})
}

func TestDeleteTorrent(t *testing.T) {
	t.Run("successful deletion", func(t *testing.T) {
		clearTables()
		torrent := &Torrent{Source: "something"}
		a.Repository.SaveTorrent(torrent)

		response := getResponse("DELETE", fmt.Sprintf("/api/torrents/%d", torrent.ID), nil)
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

func insertTorrent(source, state string) int {
	db := getDB(os.Getenv("DB_URL"))
	var id int
	err := db.QueryRow("INSERT INTO torrents (source, state) VALUES ($1, $2) RETURNING id", source, state).Scan(&id)
	if err != nil {
		log.Fatal(err)
	}
	return id
}
