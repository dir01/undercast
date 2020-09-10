package undercast

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type mediaRepository struct {
	db *mongo.Database
}

type dbMedia struct {
	ID         string     `bson:"_id"`
	DownloadId string     `bson:"downloadId"`
	Files      []string   `bson:"files"`
	Url        string     `bson:"url"`
	State      mediaState `bson:"state"`
}

func (repo *mediaRepository) Save(ctx context.Context, media *Media) error {
	dbObj := dbMedia(*media)
	opts := options.Update().SetUpsert(true)
	filter := bson.M{"_id": media.ID}
	update := bson.M{"$set": dbObj}
	if _, err := repo.db.Collection("media").UpdateOne(ctx, filter, update, opts); err != nil {
		return err
	}
	return nil
}

func (repo *mediaRepository) ListByDownloadId(ctx context.Context, downloadId string) ([]Media, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"createdAt", -1}})

	cursor, err := repo.db.Collection("media").Find(ctx, bson.D{{"downloadId", downloadId}}, findOptions)
	if err != nil {
		return nil, err
	}

	dbObjs := make([]dbMedia, 0)
	err = cursor.All(ctx, &dbObjs)
	if err != nil {
		return nil, err
	}

	medias := make([]Media, 0, len(dbObjs))
	for _, dbO := range dbObjs {
		medias = append(medias, Media(dbO))
	}

	return medias, nil
}

func (repo *mediaRepository) GetMedia(ctx context.Context, mediaId string) (*Media, error) {
	result := repo.db.Collection("media").FindOne(ctx, bson.D{{"_id", mediaId}})
	if err := result.Err(); err != nil {
		return nil, err
	}
	dbObj := &dbMedia{}
	if err := result.Decode(dbObj); err != nil {
		return nil, err
	}
	media := Media(*dbObj)
	return &media, nil
}
