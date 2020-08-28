package main

import (
	"log"
	"os"
	"undercast"
)

func main() {
	server, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           os.Getenv("MONGO_URI"),
		UIDevServerURL:     os.Getenv("UI_DEV_SERVER_URL"),
		SessionSecret:      os.Getenv("SESSION_SECRET"),
		GlobalPassword:     os.Getenv("GLOBAL_PASSWORD"),
		TorrentsDownloader: undercast.NewTorrentsDownloader(),
	})
	if err != nil {
		log.Fatal(err)
	}
	addr := ":4242"
	log.Println("Serving at address " + addr)
	log.Fatal(server.ListenAndServe(addr))
}
