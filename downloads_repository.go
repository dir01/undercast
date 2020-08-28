package undercast

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type dbDownload struct {
	ID                 string    `bson:"_id"`
	Source             string    `bson:"source"`
	Name               string    `bson:"name"`
	CreatedAt          time.Time `bson:"createdAt"`
	TotalBytes         int64     `bson:"totalBytes"`
	CompleteBytes      int64     `bson:"completeBytes"`
	Files              []string  `bson:"files"`
	RootDir            string    `bson:"rootDir"`
	IsDownloadComplete bool      `bson:"isDownloadComplete"`
}

type downloadsRepository struct {
	db *mongo.Database
}

func (r *downloadsRepository) Save(ctx context.Context, download *Download) error {
	dbObj := dbDownload(*download)
	opts := options.Update().SetUpsert(true)
	filter := bson.D{{"_id", dbObj.ID}}
	update := bson.D{{"$set", dbObj}}
	_, err := r.db.Collection("downloads").UpdateOne(ctx, filter, update, opts)
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

func (r *downloadsRepository) GetById(ctx context.Context, id string) (*Download, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := r.db.Collection("downloads").Find(ctx, bson.D{{"_id", id}}, findOptions)
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

	return &downloads[0], nil
}

func (r *downloadsRepository) ListIncomplete(ctx context.Context) ([]Download, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := r.db.Collection("downloads").Find(ctx, bson.D{{"isDownloadComplete", false}}, findOptions)
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
