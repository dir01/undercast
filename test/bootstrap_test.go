package server_test

import (
	"fmt"
	"github.com/aws/aws-sdk-go/service/s3"
	uuid "github.com/satori/go.uuid"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"io/ioutil"
	"log"
	"os"
	"path"
	"testing"
	"time"
	"undercast"
	"undercast/mocks"
)

func TestBootstrapServer(t *testing.T) {
	log.SetOutput(ioutil.Discard)
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
		Name:    "Some-torrent-name",
		Files:   []string{"foo/bar_1", "foo/bar_2"},
		RootDir: "/some/directory",
	})

	download, err := findOne(suite.db, "downloads", map[string]string{"_id": id})
	suite.Require().NoError(err)
	suite.Assert().Equal("Some-torrent-name", download["name"])
	suite.Assert().Equal(primitive.A{"foo/bar_1", "foo/bar_2"}, download["files"])
	suite.Assert().Equal("/some/directory", download["rootDir"])

	media, err := findOne(suite.db, "media", map[string]string{"downloadId": id})
	suite.Require().NoError(err)
	suite.Assert().Equal("waiting", media["state"])
}

func (suite *BootstrapSuite) TestDownloadComplete() {
	downloadId := uuid.NewV4().String()
	mediaId := uuid.NewV4().String()
	testDir := getCurrentTestDir()
	tempDir := os.TempDir()
	defer os.Remove(path.Join(tempDir, mediaId+".mp3"))

	err := insertDownload(suite.db, &downloadOpts{
		ID:      downloadId,
		Source:  "some://whatever-source",
		RootDir: testDir,
		Files:   []string{"one.mp3", "two.mp3"},
	})
	suite.Require().NoError(err)

	err = insertMedia(suite.db, &mediaOpts{
		ID:         mediaId,
		DownloadID: downloadId,
		Status:     "waiting",
		Files:      []string{"one.mp3", "two.mp3"},
	})
	suite.Require().NoError(err)

	var onProgress func(id string, p *undercast.DownloadProgress)
	suite.downloaderMock.OnProgressFunc = func(fn func(id string, p *undercast.DownloadProgress)) {
		onProgress = fn
	}

	fakeS3 := getFakeS3("test-bucket")
	defer fakeS3.Server.Close()

	_, err = undercast.Bootstrap(undercast.Options{
		MongoURI:           suite.mongoURI,
		TorrentsDownloader: suite.downloaderMock,
		TempDir:            tempDir,
		S3BucketName:       fakeS3.Bucket,
		S3Config:           fakeS3.Config,
	})
	onProgress(downloadId, &undercast.DownloadProgress{
		TotalBytes:         int64(100),
		CompleteBytes:      int64(100),
		IsDownloadComplete: true,
	})

	download, err := findOne(suite.db, "downloads", map[string]string{"_id": downloadId})
	suite.Require().NoError(err)
	suite.Assert().EqualValues(100, download["totalBytes"])
	suite.Assert().EqualValues(100, download["completeBytes"])
	suite.Assert().Equal(true, download["isDownloadComplete"])

	suite.Assert().Eventually(func() bool {
		media, err := findOne(suite.db, "media", map[string]string{"downloadId": downloadId})
		if err != nil {
			return false
		}
		expectedUrl := fmt.Sprintf("https://test-bucket.s3.eu-central-1.amazonaws.com/media/%s.mp3", media["_id"])
		return media["url"] == expectedUrl && media["state"] == "uploaded"
	}, 1*time.Second, 100*time.Millisecond)

	suite.Assert().Eventually(func() bool {
		_, err := findOne(suite.db, "episodes", map[string]string{"mediaId": mediaId})
		return err == nil
	}, 1*time.Second, 100*time.Millisecond)

	list, err := (fakeS3.Client.ListObjectsV2(&s3.ListObjectsV2Input{Bucket: &fakeS3.Bucket}))
	suite.Require().NoError(err)
	suite.Require().Len(list.Contents, 2)

	// This etag corresponds to that of a correctly converted file
	suite.Assert().Equal(`"b7ead2f1bbbbecb345870f7f1ceb9ef3"`, *list.Contents[1].ETag)
	suite.Assert().Equal(fmt.Sprintf("media/%s.mp3", mediaId), *list.Contents[1].Key)

	suite.Assert().Equal("feeds/feed.xml", *list.Contents[0].Key)
}
