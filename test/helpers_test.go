package server_test

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"net/http"
	"net/http/httptest"
	"time"
)

func (s *ServerSuite) requestAPI(method string, url string, body interface{}) *httptest.ResponseRecorder {
	bodyBytes, err := json.Marshal(body)
	s.Require().NoError(err)
	req, _ := http.NewRequest(method, url, bytes.NewReader(bodyBytes))
	rr := httptest.NewRecorder()
	s.server.ServeHTTP(rr, req)
	return rr
}

func (s *ServerSuite) findOneAsJSON(collectionName string, filter interface{}) string {
	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	var result map[string]interface{}
	err := s.db.Collection(collectionName).FindOne(ctx, filter).Decode(&result)
	s.Assert().NoError(err)
	b, err := json.Marshal(result)
	s.Assert().NoError(err)
	return string(b)
}

func (s *ServerSuite) getDatabase(dbName string) (*mongo.Database, error) {
	client, err := mongo.NewClient(options.Client().ApplyURI(s.mongoURI))
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	if err = client.Connect(ctx); err != nil {
		return nil, err
	}

	return client.Database(dbName), nil
}

func (s *ServerSuite) getMongoURI() (string, error) {
	ctx := context.Background()

	req := testcontainers.ContainerRequest{
		Image:        "mongo",
		ExposedPorts: []string{"27017/tcp"},
		WaitingFor:   wait.ForListeningPort("27017").WithStartupTimeout(5 * time.Minute),
	}

	mongoC, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return "", err
	}

	s.containers = append(s.containers, mongoC)

	ip, err := mongoC.Host(ctx)
	if err != nil {
		return "", err
	}

	port, err := mongoC.MappedPort(ctx, "27017")
	if err != nil {
		return "", err
	}

	return fmt.Sprintf("mongodb://%s:%s", ip, port.Port()), nil
}
