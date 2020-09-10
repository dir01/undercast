package server_test

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
	"path"
	"runtime"
	"time"
)

func (suite *ServerSuite) TestCreateMedia() {
	testDir := getCurrentTestDir()
	downloadId := uuid.NewV4().String()

	err := insertDownload(suite.db, &downloadOpts{
		ID:                 downloadId,
		Source:             "some-source",
		RootDir:            testDir,
		Files:              []string{"one.mp3", "two.mp3", "three.mp3"},
		IsDownloadComplete: true,
	})
	suite.Require().NoError(err)

	resp := suite.requestAPI("POST", "/api/media", map[string]interface{}{
		"id":         "some-media-id",
		"downloadId": downloadId,
		"files":      []string{"one.mp3", "two.mp3"},
	})

	suite.Assert().Equal(http.StatusOK, resp.Code)

	suite.Assert().Eventually(func() bool {
		media, err := findOne(suite.db, "media", map[string]string{"_id": "some-media-id"})
		if err != nil {
			return false
		}
		return media["state"] == "uploaded"
	}, 10*time.Second, 100*time.Millisecond)
}

func getCurrentTestDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}
