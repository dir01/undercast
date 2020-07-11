package undercast

import (
	"context"

	"go.mongodb.org/mongo-driver/mongo"
)

type dbDownload struct {
	ID     string `bson:"_id"`
	Source string `bson:"source"`
}

type downloadsRepository struct {
	db *mongo.Database
}

func (r *downloadsRepository) Save(ctx context.Context, download *Download) error {
	dbObj := dbDownload(*download)
	_, err := r.db.Collection("downloads").InsertOne(ctx, dbObj)
	if err != nil {
		return err
	}
	return nil
}
