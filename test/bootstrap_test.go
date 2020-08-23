package server_test

import (
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/mongo"
	"testing"
	"undercast"
	"undercast/mocks"
)

func TestBootstrapServer(t *testing.T) {
	require := require.New(t)
	mongoURI, err := getMongoURI()
	db, err := getDatabase(mongoURI)
	require.NoError(err)
	require.NoError(insertDownload(db, &downloadOpts{Source: "some://source", IsDownloadComplete: true}))
	require.NoError(insertDownload(db, &downloadOpts{Source: "some://other-source", IsDownloadComplete: false}))
	require.NoError(insertDownload(db, &downloadOpts{Source: "some://yet-another-source", IsDownloadComplete: false}))
	s := &BootstrapSuite{db: db, mongoURI: mongoURI}
	suite.Run(t, s)
}

type BootstrapSuite struct {
	suite.Suite
	mongoURI       string
	db             *mongo.Database
	downloaderMock *mocks.DownloaderMock
}

func (suite *BootstrapSuite) SetupTest() {
	suite.downloaderMock = &mocks.DownloaderMock{
		DownloadFunc:   func(id, source string) error { return nil },
		IsMatchingFunc: nil,
		OnProgressFunc: func(fn func(id string, di *undercast.DownloadInfo)) {},
	}
}

func (suite *BootstrapSuite) TestIncompleteDownloadsResumed() {
	_, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           suite.mongoURI,
		TorrentsDownloader: suite.downloaderMock,
	})

	suite.Require().NoError(err)
	suite.Require().Len(suite.downloaderMock.DownloadCalls(), 2)
	suite.Require().Equal("some://other-source", suite.downloaderMock.DownloadCalls()[0].Source)
	suite.Require().Equal("some://yet-another-source", suite.downloaderMock.DownloadCalls()[1].Source)
}

func (suite *BootstrapSuite) TestProgressUpdate() {
	var onProgress func(id string, di *undercast.DownloadInfo)
	suite.downloaderMock.OnProgressFunc = func(fn func(id string, di *undercast.DownloadInfo)) {
		onProgress = fn
	}

	_, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           suite.mongoURI,
		TorrentsDownloader: suite.downloaderMock,
	})
	id := suite.downloaderMock.DownloadCalls()[0].ID
	onProgress(id, &undercast.DownloadInfo{
		TotalBytes:    int64(100),
		CompleteBytes: int64(1),
	})

	otherDownload, err := findOneAsJSON(suite.db, "downloads", map[string]string{"_id": id})
	suite.Require().NoError(err)
	suite.Assert().EqualValues(100, gjson.Get(otherDownload, "totalBytes").Value())
	suite.Assert().EqualValues(1, gjson.Get(otherDownload, "completeBytes").Value())
}
