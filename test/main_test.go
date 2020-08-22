package server_test

import (
	"context"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"undercast"
	"undercast/mocks"
)

func TestServer(t *testing.T) {
	s := &ServerSuite{
		globalPassword:     "qwerty",
		torrentsDownloader: &mocks.Downloader{},
	}

	if mongoURI, err := getMongoURI(); err == nil {
		s.mongoURI = mongoURI
	} else {
		t.Error(err)
	}

	s.torrentsDownloader.On("OnProgress", mock.AnythingOfType("func(string, *undercast.DownloadInfo)")).Return()

	if server, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           s.mongoURI,
		SessionSecret:      "some-secret",
		GlobalPassword:     s.globalPassword,
		TorrentsDownloader: s.torrentsDownloader,
	}); err == nil {
		s.server = server
	} else {
		t.Error(err)
	}

	if db, err := getDatabase(s.mongoURI); err == nil {
		s.db = db
	} else {
		t.Error(err)
	}

	suite.Run(t, s)
}

type ServerSuite struct {
	suite.Suite
	mongoURI           string
	server             *undercast.Server
	db                 *mongo.Database
	containers         []testcontainers.Container
	globalPassword     string
	tempCookies        []string
	torrentsDownloader *mocks.Downloader
}

func (s *ServerSuite) TearDownSuite() {
	ctx := context.Background()
	for _, c := range s.containers {
		_ = c.Terminate(ctx)
	}
}

func (s *ServerSuite) SetupTest() {
	s.tempCookies = []string{}
	err := dropDb(s.mongoURI)
	if err != nil {
		panic(err)
	}
}
