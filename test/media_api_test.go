package server_test

import (
	uuid "github.com/satori/go.uuid"
	"net/http"
	"path"
	"runtime"
)

func (s *ServerSuite) TestCreateMedia() {
	testDir := getCurrentTestDir()
	downloadId := uuid.NewV4().String()
	s.insertDownload(&downloadOpts{
		ID:      downloadId,
		Source:  "some-source",
		RootDir: testDir,
		Files:   []string{"one.mp3", "two.mp3", "three.mp3"},
	})

	resp := s.requestAPI("POST", "/api/media", map[string]interface{}{
		"id":         "some-media-id",
		"downloadId": downloadId,
		"files":      []string{"one.mp3", "two.mp3"},
	})

	s.Require().Equal(http.StatusOK, resp.Code)
}

func getCurrentTestDir() string {
	_, filename, _, _ := runtime.Caller(0)
	return path.Dir(filename)
}
