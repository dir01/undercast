test:
	DB_HOST=localhost DB_PORT=5432 DB_USER=undercast DB_PASSWORD=undercast DB_NAME=undercast go test -v ./...
runserver:
	DB_HOST=localhost DB_PORT=5432 DB_USER=undercast DB_PASSWORD=undercast DB_NAME=undercast go run ./cmd/server
