ifneq (,$(wildcard ./.env))
    include .env
    export
endif

.PHONY: test
test:
	go test -v ./...

test-cover:
	go test -coverpkg ./... -v -coverprofile=coverage.out -covermode=count ./...

runserver:
	MONGO_URI=mongodb://localhost/test \
	UI_DEV_SERVER_URL=$(UI_DEV_SERVER_URL) \
	SESSION_SECRET=seekreet \
	GLOBAL_PASSWORD=batman42 \
	go run ./cmd/server



runserver-dev: UI_DEV_SERVER_URL="http://localhost:8080"
runserver-dev:
	go get github.com/cespare/reflex
	reflex -s -r '\.go$$' -R '^ui/' make runserver

build-ui:
	time bash -c "cd ui; npm i; npm run build"

run-ui:
	cd ui; npm run dev

build:
	time go build -o ./bin/server ./cmd/server/main.go