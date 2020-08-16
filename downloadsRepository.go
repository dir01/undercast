package undercast

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type dbDownload struct {
	ID        string    `bson:"_id"`
	Source    string    `bson:"source"`
	CreatedAt time.Time `bson:"createdAt"`
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

func (r *downloadsRepository) List(ctx context.Context) ([]Download, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := r.db.Collection("downloads").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	dbObjs := make([]dbDownload, 0)
	err = cursor.All(ctx, &dbObjs)
	if err != nil {
		return nil, err
	}

	downloads := make([]Download, 0, len(dbObjs))
	for _, dbO := range dbObjs {
		downloads = append(downloads, Download(dbO))
	}

	return downloads, nil
}
