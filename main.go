package main

import (
	"fmt"
	"os"
)

func main() {
	fmt.Println("Hellp")
	a := App{}
	a.Initialize(
		os.Getenv("DB_HOST"),
		os.Getenv("DB_PORT"),
		os.Getenv("DB_USER"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)
	a.Run(":8080")
}
