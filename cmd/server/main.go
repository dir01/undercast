package main

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
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
		TorrentsDownloader: undercast.NewTorrentsDownloader("./data"),
		S3Config: &aws.Config{
			Credentials: credentials.NewStaticCredentials(
				os.Getenv("AWS_ACCESS_KEY_ID"),
				os.Getenv("AWS_SECRET_ACCESS_KEY"), "",
			),
			Region: aws.String(os.Getenv("AWS_REGION")),
		},
		S3BucketName: os.Getenv("S3_BUCKET"),
		TempDir:      "./data/tmp",
	})
	if err != nil {
		log.Fatal(err)
	}
	addr := ":4242"
	log.Println("Serving at address " + addr)
	log.Fatal(server.ListenAndServe(addr))
}
