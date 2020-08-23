package server_test

import (
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func (s *ServerSuite) TestCreateDownload() {
	magnetLink := "magnet:?xt=urn:btih:980E4184AEE6F326A9F9E2EE3E9D40ACAA90BC40"

	var downloadedId string
	s.torrentsDownloader.DownloadFunc = func(id, source string) error {
		downloadedId = id
		return nil
	}

	resp := s.requestAPI("POST", "/api/downloads", map[string]string{"source": magnetLink})
	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(magnetLink, gjson.Get(resp.Body.String(), "payload.source").Value())
	s.Assert().Equal(downloadedId, gjson.Get(resp.Body.String(), "payload.id").Value())
	dbResultStr := s.findOneAsJSON("downloads", bson.M{})
	s.Assert().Equal(magnetLink, gjson.Get(dbResultStr, "source").Value())
}

func (s *ServerSuite) TestListDownloads() {
	s.insertDownload(&downloadOpts{Source: "some://source"})
	resp := s.requestAPI("GET", "/api/downloads", nil)
	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal("some://source", gjson.Get(resp.Body.String(), "payload.0.source").Value())
}
