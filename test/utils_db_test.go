package server_test

import (
	"context"
	"encoding/json"
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

func (s *ServerSuite) findOneAsJSON(collectionName string, filter interface{}) string {
	str, err := findOneAsJSON(s.db, collectionName, filter)
	s.Require().NoError(err)
	return str
}

func findOneAsJSON(db *mongo.Database, collectionName string, filter interface{}) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var result map[string]interface{}
	err := db.Collection(collectionName).FindOne(ctx, filter).Decode(&result)
	if err != nil {
		return "", err
	}
	b, err := json.Marshal(result)
	if err != nil {
		return "", err
	}
	return string(b), nil
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
	ID                 string `bson:"_id"`
	Source             string `bson:"source"`
	IsDownloadComplete bool   `bson:"isDownloadComplete"`
}

func (s *ServerSuite) insertDownload(opts *downloadOpts) {
	err := insertDownload(s.db, opts)
	s.Require().NoError(err)
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
