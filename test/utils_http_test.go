package server_test

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
)

func (s *ServerSuite) requestAPI(method string, url string, body interface{}) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	s.Require().NoError(err)
	req, _ := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	req.Header["Cookie"] = s.tempCookies
	rr := httptest.NewRecorder()
	s.server.ServeHTTP(rr, req)
	s.tempCookies = append(s.tempCookies, rr.HeaderMap["Set-Cookie"]...)
	return rr
}
