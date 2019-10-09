package main

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

var a App

func TestMain(m *testing.M) {
	a = App{}
	a.Initialize(
		os.Getenv("TEST_DB_HOST"),
		os.Getenv("TEST_DB_PORT"),
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))
	createTable()

	code := m.Run()

	dropTable()

	os.Exit(code)
}

func TestEmptyTable(t *testing.T) {
	clearTable()

	req, _ := http.NewRequest("GET", "/episodes", nil)
	response := executeRequest(req)

	checkResponseCode(t, response, http.StatusOK)
}

func TestCreateEpisodeWithMagnet(t *testing.T) {
	clearTable()

	payload := []byte(`{
		"name": "Around the world in 80 days",
		"magnet": "magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c"
	}`)
	req, _ := http.NewRequest("POST", "/episodes", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, response, http.StatusCreated)
	checkResponseBody(t, response, `{"id":1,"name":"Around the world in 80 days","magnet":"magnet:?xt=urn:btih:1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c","url":""}`)
}

func TestCreateEpisodeWithUrl(t *testing.T) {
	clearTable()

	payload := []byte(`{
		"name": "Around the world in 80 days",
		"url": "http://legittorrents.info/download.php?id=1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c"
	}`)
	req, _ := http.NewRequest("POST", "/episodes", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, response, http.StatusCreated)
	checkResponseBody(t, response, `{"id":1,"name":"Around the world in 80 days","magnet":"","url":"http://legittorrents.info/download.php?id=1ce53bc6bd5d16b4f92f9cd40bc35e10724f355c"}`)
}

func TestFailToCreateEpisode(t *testing.T) {
	clearTable()

	payload := []byte(`{
		"name": "Around the world in 80 days"
	}`)
	req, _ := http.NewRequest("POST", "/episodes", bytes.NewBuffer(payload))
	response := executeRequest(req)

	checkResponseCode(t, response, http.StatusInternalServerError)
	checkResponseBody(t, response, `{"error":"pq: new row for relation \"episodes\" violates check constraint \"require_magnet_or_url\""}`)
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	rr := httptest.NewRecorder()
	a.Router.ServeHTTP(rr, req)
	return rr
}

func checkResponseCode(t *testing.T, response *httptest.ResponseRecorder, expectedCode int) {
	actualCode := response.Code
	if expectedCode != actualCode {
		t.Errorf("Expected response code %d. Got %d\n%s", expectedCode, actualCode, response.Body)
	}
}

func checkResponseBody(t *testing.T, response *httptest.ResponseRecorder, expectedBody string) {
	actualBody := response.Body.String()
	if expectedBody != actualBody {
		t.Errorf("Unexpected response body\nEXPECTED:\n%s\nACTUAL:\n%s", expectedBody, actualBody)
	}
}

func createTable() {
	if _, err := a.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	a.DB.Exec("DELETE FROM episodes")
	a.DB.Exec("ALTER SEQUENCE episodes_id_seq RESTART WITH 1")
}

func dropTable() {
	a.DB.Exec("DROP TABLE episodes")
}

const tableCreationQuery = `
CREATE TABLE IF NOT EXISTS episodes (
	id SERIAL PRIMARY KEY,
	name VARCHAR(500) NOT NULL,
	magnet VARCHAR(500),
	url VARCHAR(500),
	CONSTRAINT require_magnet_or_url CHECK (
		(case when magnet is null or length(magnet) = 0 then 0 else 1 end)
		<> 
		(case when url is null or length(url) = 0 then 0 else 1 end)
	)
)`
