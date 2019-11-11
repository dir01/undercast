package bittorrent

import (
	"fmt"
	"path"
	"strings"
	"time"

	"undercast/server"

	"github.com/anacrolix/torrent"
	anacrolix "github.com/anacrolix/torrent"
)

// NewClient creates new Client
func NewClient(dataDir string) (server.TorrentClient, error) {
	c := &torrentClient{}
	cfg := anacrolix.NewDefaultClientConfig()
	cfg.DataDir = dataDir
	c.dataDir = dataDir
	tc, err := torrent.NewClient(cfg)
	c.client = tc
	c.torrentsMap = make(map[int]*anacrolix.Torrent)
	return c, err
}

type torrentClient struct {
	client      *anacrolix.Client
	torrentsMap map[int]*anacrolix.Torrent
	callback    func(id int, state server.TorrentState)
	dataDir     string
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
		fmt.Println(t.Files()[0].Path())
		i := t.Info()
		state := &server.TorrentState{
			Name:           i.Name,
			FilePaths:      copyFilePaths(t, tc.dataDir),
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

func copyFilePaths(t *anacrolix.Torrent, root string) (filepaths []string) {
	for _, f := range t.Files() {
		filepaths = append(filepaths, path.Join(root, f.Path()))
	}
	return
}
