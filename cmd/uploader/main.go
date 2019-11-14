package main

import (
	"fmt"
	"os"
	"undercast/server"
)

func main() {
	filename := os.Args[1]
	u := server.NewUploader("AKIA3LO6ZJ4SHU5ONFTE", "Ir+lELP1Msix7LlJoN18nQeOBGhd8m4Cj0wnDtgX", "", "eu-west-1", "unit35173")
	url, err := u.UploadFile(filename)
	fmt.Println(url, err)
}
