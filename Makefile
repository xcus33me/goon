build:
	@go build -o bin/goon cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/goon

seed:
	@go build -o bin/goon cmd/seed/main.go && ./bin/goon

db-status:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) goose -dir=$(migrationPath) status

up:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) goose -dir=$(migrationPath) up

down:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) goose -dir=$(migrationPath) down

reset:
	@GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_URL) goose -dir=$(migrationPath) reset
