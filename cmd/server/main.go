package main

import (
	"os"
	"undercast/bittorrent"
	"undercast/server"
)

func main() {
	t, err := bittorrent.NewClient(os.Getenv("DATA_DIR"))
	if err != nil {
		panic("Failed to initialize bittorrent client:\n" + err.Error())
	}
	a := server.App{Torrent: t}
	a.Initialize(
		os.Getenv("DB_URL"),
		os.Getenv("UI_DEV_SERVER_URL"),
	)
	a.Run(":8080")
}
