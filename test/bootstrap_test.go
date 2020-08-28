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
	require.NoError(err)
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
	suite.downloaderMock = &mocks.DownloaderMock{}
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

func (suite *BootstrapSuite) TestInfoUpdate() {
	var onInfo func(id string, di *undercast.DownloadInfo)
	suite.downloaderMock.OnInfoFunc = func(fn func(id string, di *undercast.DownloadInfo)) {
		onInfo = fn
	}

	_, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           suite.mongoURI,
		TorrentsDownloader: suite.downloaderMock,
	})
	id := suite.downloaderMock.DownloadCalls()[0].ID
	onInfo(id, &undercast.DownloadInfo{
		Name:  "Some-torrent-name",
		Files: []string{"foo/bar_1", "foo/bar_2"},
	})

	download, err := findOneAsJSON(suite.db, "downloads", map[string]string{"_id": id})
	suite.Require().NoError(err)
	suite.Assert().Equal("Some-torrent-name", gjson.Get(download, "name").Value())
	suite.Assert().Equal(`["foo/bar_1","foo/bar_2"]`, gjson.Get(download, "files").String())
}

func (suite *BootstrapSuite) TestProgressUpdate() {
	var onProgress func(id string, p *undercast.DownloadProgress)
	suite.downloaderMock.OnProgressFunc = func(fn func(id string, p *undercast.DownloadProgress)) {
		onProgress = fn
	}

	_, err := undercast.Bootstrap(undercast.Options{
		MongoURI:           suite.mongoURI,
		TorrentsDownloader: suite.downloaderMock,
	})
	id := suite.downloaderMock.DownloadCalls()[0].ID
	onProgress(id, &undercast.DownloadProgress{
		TotalBytes:    int64(100),
		CompleteBytes: int64(1),
	})

	download, err := findOneAsJSON(suite.db, "downloads", map[string]string{"_id": id})
	suite.Require().NoError(err)
	suite.Assert().EqualValues(100, gjson.Get(download, "totalBytes").Value())
	suite.Assert().EqualValues(1, gjson.Get(download, "completeBytes").Value())
}
