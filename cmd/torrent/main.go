package main

import (
	"fmt"
	"os"
	"time"

	"undercast/bittorrent"
	"undercast/server"
)

func main() {
	c, _ := bittorrent.NewClient()

	downloading := true

	c.AddTorrent(11, os.Args[1])
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
