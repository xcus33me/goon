build:
	@go build -o bin/goon cmd/api/main.go

run: build
	@./bin/goon

test:
	@go test -v ./...