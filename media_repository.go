package undercast

import (
	"context"
	"go.mongodb.org/mongo-driver/mongo"
)

type mediaRepository struct {
	db *mongo.Database
}

type dbMedia struct {
	ID         string   `bson:"id"`
	DownloadId string   `bson:"download_id"`
	Files      []string `bson:"files"`
}

func (repo *mediaRepository) Save(ctx context.Context, media *Media) error {
	dbObj := dbMedia(*media)
	if _, err := repo.db.Collection("media").InsertOne(ctx, dbObj); err != nil {
		return err
	}
	return nil
}
