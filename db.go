package undercast

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"strings"
	"time"
)

func GetDb(mongoURI string) (*mongo.Database, error) {
	lastSlashIndex := strings.LastIndex(mongoURI, "/")
	dbName := mongoURI[lastSlashIndex+1:]
	mongoURI = mongoURI[:lastSlashIndex]

	client, err := mongo.NewClient(options.Client().ApplyURI(mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	err = client.Connect(ctx)
	if err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}
