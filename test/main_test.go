package server_test

import (
	"context"
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
		torrentsDownloader: &mocks.DownloaderMock{},
	}

	mongoURI, err := getMongoURI()
	if err != nil {
		t.Error(err)
	}

	s.torrentsDownloader.OnProgressFunc = func(fn func(string, *undercast.DownloadProgress)) {}

	if server, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           mongoURI,
		SessionSecret:      "some-secret",
		GlobalPassword:     s.globalPassword,
		TorrentsDownloader: s.torrentsDownloader,
	}); err == nil {
		s.server = server
	} else {
		t.Error(err)
	}

	if db, err := getDatabase(mongoURI); err == nil {
		s.db = db
	} else {
		t.Error(err)
	}

	suite.Run(t, s)
}

type ServerSuite struct {
	suite.Suite
	server             *undercast.Server
	db                 *mongo.Database
	containers         []testcontainers.Container
	globalPassword     string
	tempCookies        []string
	torrentsDownloader *mocks.DownloaderMock
}

func (s *ServerSuite) TearDownSuite() {
	ctx := context.Background()
	for _, c := range s.containers {
		_ = c.Terminate(ctx)
	}
}

func (s *ServerSuite) SetupTest() {
	s.tempCookies = []string{}
	err := dropDb(s.db)
	if err != nil {
		panic(err)
	}
}
