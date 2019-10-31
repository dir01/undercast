package main

import (
	"testing"
	"undercast/server"
)

type addTorrentCall struct {
	id     int
	source string
}

type torrentMock struct {
	addTorrentCalls []addTorrentCall
	callback        func(id int, state server.TorrentState)
}

func setupTorrentMock(a *server.App) *torrentMock {
	tm := &torrentMock{}
	a.Torrent = tm
	return tm
}

func (tm *torrentMock) AddTorrent(id int, source string) error {
	call := addTorrentCall{id: id, source: source}
	tm.addTorrentCalls = append(tm.addTorrentCalls, call)
	return nil
}

func (tm *torrentMock) OnTorrentChanged(callback func(id int, state server.TorrentState)) {
	tm.callback = callback
}

func (tm *torrentMock) assertTorrentAdded(t *testing.T, id int, source string) {
	if !tm.isTorrentAdded(id, source) {
		t.Errorf("Torrent was not added to client, but it was expected to: \nid=\"%d\" and source=\"%s\"", id, source)
	}
}

func (tm *torrentMock) assertTorrentNotAdded(t *testing.T, id int, source string) {
	if tm.isTorrentAdded(id, source) {
		t.Errorf("Torrent was added to client, but it was expected not to: \nid=\"%d\" and source=\"%s\"", id, source)
	}
}

func (tm *torrentMock) isTorrentAdded(id int, source string) bool {
	for _, c := range tm.addTorrentCalls {
		if c.id == id && c.source == source {
			return true
		}
	}
	return false
}
