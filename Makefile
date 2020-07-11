.PHONY: test
test:
	go test -v ./...

runserver:
	MONGO_URI=mongodb://localhost/ \
	MONGO_DB_NAME=test \
	UI_DEV_SERVER_URL=$(UI_DEV_SERVER_URL) \
	go run ./cmd/server



runserver-dev: UI_DEV_SERVER_URL=http://localhost:4200
runserver-dev: runserver

build-ui:
	time bash -c "cd ui; npm i; npm run build"

run-ui:
	cd ui; npm start

build:
	time go build -o ./bin/server ./cmd/server/main.go