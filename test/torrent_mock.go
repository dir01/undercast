package main

import "undercast/server"

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
