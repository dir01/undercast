test:
	DB_HOST=localhost DB_PORT=5432 DB_USER=undercast DB_PASSWORD=undercast DB_NAME=undercast go test -v ./...

runserver:
	DB_HOST=localhost DB_PORT=5432 DB_USER=undercast DB_PASSWORD=undercast UI_DEV_SERVER_URL=$(UI_DEV_SERVER_URL) DB_NAME=undercast DATA_DIR=./data go run ./cmd/server

runserver-dev: UI_DEV_SERVER_URL=http://localhost:4200
runserver-dev: runserver
