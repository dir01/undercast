package bittorrent

import (
	"strings"
	"time"

	"undercast/server"

	"github.com/anacrolix/torrent"
	anacrolix "github.com/anacrolix/torrent"
	anacrolixMetainfo "github.com/anacrolix/torrent/metainfo"
)

// NewClient creates new Client
func NewClient(dataDir string) (server.TorrentClient, error) {
	c := &torrentClient{}
	cfg := anacrolix.NewDefaultClientConfig()
	cfg.DataDir = dataDir
	tc, err := torrent.NewClient(cfg)
	c.client = tc
	c.torrentsMap = make(map[int]*anacrolix.Torrent)
	return c, err
}

type torrentClient struct {
	client      *anacrolix.Client
	torrentsMap map[int]*anacrolix.Torrent
	callback    func(id int, state server.TorrentState)
}

func (tc *torrentClient) AddTorrent(id int, magnetOrURLOrTorrent string) error {
	var t *anacrolix.Torrent
	var e error

	if strings.HasPrefix(magnetOrURLOrTorrent, "magnet:") {
		t, e = tc.client.AddMagnet(magnetOrURLOrTorrent)
	} else {
		t, e = tc.client.AddTorrentFromFile(magnetOrURLOrTorrent)
	}
	if e != nil {
		return e
	}

	tc.torrentsMap[id] = t

	go func() {
		<-t.GotInfo()
		i := t.Info()
		state := &server.TorrentState{
			Name:           i.Name,
			FileNames:      copyFileNames(i),
			BytesCompleted: t.BytesCompleted(),
			BytesMissing:   t.BytesMissing(),
			Done:           false,
		}

		tc.callback(id, *state)
		t.DownloadAll()

		done := false
		for !done {
			time.Sleep(5 * time.Second)
			state.BytesCompleted = t.BytesCompleted()
			state.BytesMissing = t.BytesMissing()
			if t.BytesMissing() == 0 {
				state.Done = true
				done = true
			}
			tc.callback(id, *state)
		}
	}()
	return nil
}

func (tc *torrentClient) OnTorrentChanged(callback func(id int, state server.TorrentState)) {
	tc.callback = callback
}

func copyFileNames(i *anacrolixMetainfo.Info) (filenames []string) {
	filenames = make([]string, len(i.Files))
	for _, f := range i.Files {
		filenames = append(filenames, f.DisplayPath(i))
	}
	return
}
