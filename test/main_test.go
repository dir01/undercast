package server_test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"undercast"
)

func TestServer(t *testing.T) {
	s := &ServerSuite{
		globalPassword: "qwerty",
	}

	if mongoURI, err := s.getMongoURI(); err == nil {
		s.mongoURI = mongoURI
	} else {
		t.Error(err)
	}

	dbName := "test"

	if server, err := undercast.Bootstrap(undercast.Options{
		MongoURI:       s.mongoURI,
		MongoDbName:    dbName,
		SessionSecret:  "some-secret",
		GlobalPassword: s.globalPassword,
	}); err == nil {
		s.server = server
	} else {
		t.Error(err)
	}

	if db, err := s.getDatabase(dbName); err == nil {
		s.db = db
	} else {
		t.Error(err)
	}

	suite.Run(t, s)
}

type ServerSuite struct {
	suite.Suite
	mongoURI       string
	server         *undercast.Server
	db             *mongo.Database
	containers     []testcontainers.Container
	globalPassword string
	tempCookies    []string
}

func (s *ServerSuite) TearDownSuite() {
	ctx := context.Background()
	for _, c := range s.containers {
		_ = c.Terminate(ctx)
	}
}

func (s *ServerSuite) SetupTest() {
	s.tempCookies = []string{}
}
