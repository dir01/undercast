package server_test

import (
	"github.com/tidwall/gjson"
	"go.mongodb.org/mongo-driver/bson"
	"net/http"
)

func (s *ServerSuite) TestCreateDownload() {
	resp := s.requestAPI("POST", "/api/downloads", map[string]string{"source": "magnet://whatever"})
	s.Assert().Equal(http.StatusOK, resp.Code)
	s.Assert().Equal(`magnet://whatever`, gjson.Get(resp.Body.String(), "payload.source").Value())
	dbResultStr := s.findOneAsJSON("downloads", bson.M{})
	s.Assert().Equal("magnet://whatever", gjson.Get(dbResultStr, "source").Value())
}
