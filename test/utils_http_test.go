package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (suite *ServerSuite) requestAPI(method string, url string, body interface{}) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	suite.Require().NoError(err)
	req, _ := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	req.Header["Cookie"] = suite.tempCookies
	rr := httptest.NewRecorder()
	suite.server.ServeHTTP(rr, req)
	suite.tempCookies = append(suite.tempCookies, rr.HeaderMap["Set-Cookie"]...)
	return rr
}
