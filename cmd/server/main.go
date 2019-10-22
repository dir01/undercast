package main

import (
	"os"
	"undercast/server"
)

func main() {
	a := server.App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
		os.Getenv("UI_DEV_SERVER_URL"),
	)
	a.Run(":8080")
}
