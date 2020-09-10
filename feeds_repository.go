package undercast

import (
	"context"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"time"
)

type feedsRepository struct {
	db *mongo.Database
}

type dbEpisode struct {
	ID       string    `bson:"_id"`
	MediaId  string    `bson:"mediaId"`
	MediaURL string    `bson:"mediaUrl"`
	Date     time.Time `bson:"date"`
}

func (repo *feedsRepository) InsertEpisode(ctx context.Context, episode *Episode) error {
	dbObj := dbEpisode(*episode)
	if _, err := repo.db.Collection("episodes").InsertOne(ctx, dbObj); err != nil {
		return err
	}
	return nil
}

func (repo *feedsRepository) ListEpisodes(ctx context.Context) ([]Episode, error) {
	findOptions := options.Find()
	findOptions.SetSort(bson.D{{"date", -1}})

	cursor, err := repo.db.Collection("episodes").Find(ctx, bson.D{}, findOptions)
	if err != nil {
		return nil, err
	}

	dbObjs := make([]dbEpisode, 0)
	err = cursor.All(ctx, &dbObjs)
	if err != nil {
		return nil, err
	}

	episodes := make([]Episode, 0, len(dbObjs))
	for _, dbO := range dbObjs {
		episodes = append(episodes, Episode(dbO))
	}

	return episodes, nil
}
