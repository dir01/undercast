package main

import (
	"log"
	"os"
	"undercast"
)

func main() {
	server, err := undercast.Bootstrap(undercast.Options{
		MongoURI:    os.Getenv("MONGO_URI"),
		MongoDbName: os.Getenv("MONGO_DB_NAME"),
	})
	if err != nil {
		log.Fatal(err)
	}
	addr := ":8080"
	log.Println("Serving at address " + addr)
	log.Fatal(server.ListenAndServe(addr))
}
