package server_test

import (
	"context"
	"github.com/stretchr/testify/suite"
	"github.com/testcontainers/testcontainers-go"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"testing"
	"undercast"
	"undercast/mocks"
)

func TestServer(t *testing.T) {
	log.SetOutput(ioutil.Discard)
	s := &ServerSuite{
		globalPassword:     "qwerty",
		torrentsDownloader: &mocks.DownloaderMock{},
	}

	mongoURI, err := getMongoURI()
	if err != nil {
		t.Error(err)
	}

	s.fakeS3 = getFakeS3("test-bucket")

	if server, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           mongoURI,
		SessionSecret:      "some-secret",
		GlobalPassword:     s.globalPassword,
		TorrentsDownloader: s.torrentsDownloader,
		S3Config:           s.fakeS3.Config,
		S3BucketName:       "test-bucket",
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
	fakeS3             fakeS3
}

func (suite *ServerSuite) TearDownSuite() {
	ctx := context.Background()
	for _, c := range suite.containers {
		_ = c.Terminate(ctx)
	}
	suite.fakeS3.Server.Close()
}

func (suite *ServerSuite) SetupTest() {
	suite.tempCookies = []string{}
	err := dropDb(suite.db)
	if err != nil {
		panic(err)
	}
}
