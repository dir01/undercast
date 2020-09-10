package server_test

import (
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func (suite *ServerSuite) TestCreateDownload() {
	magnetLink := "magnet:?xt=urn:btih:980E4184AEE6F326A9F9E2EE3E9D40ACAA90BC40"
	var downloadId string
	suite.torrentsDownloader.DownloadFunc = func(id, source string) error {
		downloadId = id
		return nil
	}

	resp := suite.requestAPI("POST", "/api/downloads", map[string]string{"source": magnetLink})
	suite.Assert().Equal(http.StatusOK, resp.Code)
	suite.Assert().Equal(magnetLink, gjson.Get(resp.Body.String(), "payload.source").Value())
	suite.Assert().Equal(downloadId, gjson.Get(resp.Body.String(), "payload.id").Value())
	download, err := findOne(suite.db, "downloads", bson.M{})
	suite.Require().NoError(err)
	suite.Assert().Equal(magnetLink, download["source"])
}

func (suite *ServerSuite) TestListDownloads() {
	err := insertDownload(suite.db, &downloadOpts{Source: "some://source"})
	suite.Require().NoError(err)
	resp := suite.requestAPI("GET", "/api/downloads", nil)
	suite.Assert().Equal(http.StatusOK, resp.Code)
	suite.Assert().Equal("some://source", gjson.Get(resp.Body.String(), "payload.0.source").Value())
}
