package main

import (
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
)

func getResponse(method string, url string, body io.Reader) *httptest.ResponseRecorder {
	req, _ := http.NewRequest(method, url, body)
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponse(
	t *testing.T,
	response *httptest.ResponseRecorder,
	expectedCode int,
	expectedPayload interface{},
) {
	actualCode := response.Code
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n%s", expectedCode, actualCode, response.Body)
	}

	if expectedPayload == nil {
		return
	}

	var expectedBody string
	switch expectedPayload.(type) {
	case string:
		expectedBody = expectedPayload.(string)
	default:
		expectedBytes, _ := json.Marshal(expectedPayload)
		expectedBody = string(expectedBytes)
	}

	actualBody := response.Body.String()
	if expectedBody != actualBody {
		t.Errorf("Unexpected response body\nEXPECTED:\n%s\nACTUAL:\n%s", expectedBody, actualBody)
	}
}
