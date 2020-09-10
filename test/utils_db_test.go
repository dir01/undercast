package server_test

import (
	"context"
	"fmt"
	uuid "github.com/satori/go.uuid"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"time"
	"undercast"
)

func getMongoURI() (string, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017").WithStartupTimeout(5 * time.Minute),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", err
	}

	ip, err := mongoC.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		return "", err
	}

	mongoURI := fmt.Sprintf("mongodb://%s:%s/%s", ip, port.Port(), uuid.NewV4().String())
	return mongoURI, nil
}

func getDatabase(mongoURI string) (*mongo.Database, error) {
	return undercast.GetDb(mongoURI)
}

func findOne(db *mongo.Database, collectionName string, filter interface{}) (map[string]interface{}, error) {
	results, err := find(db, collectionName, filter)
	if err != nil {
		return nil, err
	}
	if len(results) != 1 {
		return nil, fmt.Errorf("Expected exactly one result, got %d instead", len(results))
	}
	return results[0], nil
}

func find(db *mongo.Database, collectionName string, filter interface{}) ([]map[string]interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var results []map[string]interface{}
	cursor, err := db.Collection(collectionName).Find(ctx, filter)
	if err != nil {
		return nil, err
	}
	if err := cursor.All(ctx, &results); err != nil {
		return nil, err
	}
	return results, nil
}

func dropDb(db *mongo.Database) error {
	ctx := context.Background()
	err := db.Drop(ctx)
	if err != nil {
		return err
	}
	return nil
}

type downloadOpts struct {
	ID                 string   `bson:"_id"`
	Source             string   `bson:"source"`
	IsDownloadComplete bool     `bson:"isDownloadComplete"`
	RootDir            string   `bson:"rootDir"`
	Files              []string `bson:"files"`
}

func insertDownload(db *mongo.Database, opts *downloadOpts) error {
	if opts.ID == "" {
		opts.ID = uuid.NewV4().String()
	}
	ctx := context.Background()
	_, err := db.Collection("downloads").InsertOne(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}

type mediaOpts struct {
	ID         string   `bson:"_id"`
	DownloadID string   `bson:"downloadId"`
	Status     string   `bson:"status"`
	Files      []string `bson:"files"`
}

func insertMedia(db *mongo.Database, opts *mediaOpts) error {
	if opts.ID == "" {
		opts.ID = uuid.NewV4().String()
	}
	ctx := context.Background()
	_, err := db.Collection("media").InsertOne(ctx, opts)
	if err != nil {
		return err
	}
	return nil
}
