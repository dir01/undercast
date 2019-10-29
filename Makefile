test:
	DB_URL=postgres://undercast:undercast@localhost/undercast go test -v ./...

runserver:
	DB_URL=postgres://undercast:undercast@localhost/undercast UI_DEV_SERVER_URL=$(UI_DEV_SERVER_URL) DB_NAME=undercast DATA_DIR=./data go run ./cmd/server

runserver-dev: UI_DEV_SERVER_URL=http://localhost:4200
runserver-dev: runserver

build-ui:
	time bash -c "cd ui; npm i; npm run build"