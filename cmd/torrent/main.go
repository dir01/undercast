package main

import (
	"fmt"
	"time"

	"undercast/bittorrent"
	"undercast/server"
)

func main() {
	c, _ := bittorrent.NewClient("/tmp")

	downloading := true

	c.AddTorrent(11, "/Users/dir01/Projects/undercast/bittorrent/harrypotter.torrent")
	c.OnTorrentChanged(func(id int, info server.TorrentState) {
		fmt.Println("Got info", id, info)
		if info.Done {
			downloading = false
		}
	})

	for downloading {
		select {
		case <-time.After(1 * time.Second):
			fmt.Println("tick")
		}
	}

}
