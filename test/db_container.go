package main

import (
	"context"
	"fmt"

	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type dbInfo struct {
	url       string
	terminate func()
}

func startPostgresContainer() (result dbInfo, err error) {
	result = dbInfo{}
	req := testcontainers.ContainerRequest{
		Image:        "labianchin/docker-postgres-for-testing",
		ExposedPorts: []string{"5432"},
		WaitingFor:   wait.ForListeningPort("5432/tcp"),
		Env: map[string]string{
			"POSTGRES_USER":     "test",
			"POSTGRES_PASSWORD": "test",
			"POSTGRES_DB":       "test",
		},
	}
	ctx := context.Background()
	dbContainer, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})

	if err != nil {
		return
	}

	result.terminate = func() { dbContainer.Terminate(ctx) }

	defer func() {
		if err != nil {
			result.terminate()
		}
	}()

	ip, err := dbContainer.Host(ctx)
	if err != nil {
		return
	}

	port, err := dbContainer.MappedPort(ctx, "5432")
	if err != nil {
		return
	}

	result.url = fmt.Sprintf("postgres://test:test@%s:%s/test?sslmode=disable", ip, port.Port())

	return
}
