package server_test

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/tidwall/gjson"
	"testing"
	"undercast"
	"undercast/mocks"
)

func TestBootstrapServer(t *testing.T) {
	mongoURI, err := getMongoURI()
	require := require.New(t)
	assert := assert.New(t)
	require.NoError(err)

	require.NoError(insertDownload(mongoURI, &downloadOpts{Source: "some://source", IsDownloadComplete: true}))
	require.NoError(insertDownload(mongoURI, &downloadOpts{Source: "some://other-source", IsDownloadComplete: false}))
	require.NoError(insertDownload(mongoURI, &downloadOpts{Source: "some://yet-another-source", IsDownloadComplete: false}))

	fakeTorrentsDownloader := &mocks.Downloader{}

	var onProgress func(id string, di *undercast.DownloadInfo)
	fakeTorrentsDownloader.On("OnProgress", mock.Anything).Run(func(args mock.Arguments) {
		onProgress = args[0].(func(id string, di *undercast.DownloadInfo))
	})

	downloadedMagnets := make([][]string, 0, 0)
	fakeTorrentsDownloader.On("Download", mock.AnythingOfType("string"), mock.AnythingOfType("string")).Run(func(args mock.Arguments) {
		id := args[0].(string)
		source := args[1].(string)
		downloadedMagnets = append(downloadedMagnets, []string{id, source})
	}).Return(nil)

	_, err = undercast.Bootstrap(undercast.Options{
		MongoURI:           mongoURI,
		TorrentsDownloader: fakeTorrentsDownloader,
	})
	require.NoError(err)

	require.Len(downloadedMagnets, 2)
	require.Equal("some://other-source", downloadedMagnets[0][1])
	otherSourceId := downloadedMagnets[0][0]
	require.Equal("some://yet-another-source", downloadedMagnets[1][1])
	yetAnotherSourceId := downloadedMagnets[1][0]

	onProgress(otherSourceId, &undercast.DownloadInfo{
		TotalBytes:    int64(100),
		CompleteBytes: int64(1),
	})

	onProgress(yetAnotherSourceId, &undercast.DownloadInfo{
		TotalBytes:    int64(200),
		CompleteBytes: int64(2),
	})

	otherDownload, err := findOneAsJSON(mongoURI, "downloads", map[string]string{"_id": otherSourceId})
	require.NoError(err)
	assert.Equal("some://other-source", gjson.Get(otherDownload, "source").Value())
	assert.EqualValues(100, gjson.Get(otherDownload, "totalBytes").Value())
	assert.EqualValues(1, gjson.Get(otherDownload, "completeBytes").Value())

	yetAnotherDownload, err := findOneAsJSON(mongoURI, "downloads", map[string]string{"_id": yetAnotherSourceId})
	require.NoError(err)
	assert.Equal("some://yet-another-source", gjson.Get(yetAnotherDownload, "source").Value())
	assert.EqualValues(200, gjson.Get(yetAnotherDownload, "totalBytes").Value())
	assert.EqualValues(2, gjson.Get(yetAnotherDownload, "completeBytes").Value())

}
