build:
	@go build -o bin/goon cmd/main.go

test:
	@go test -v ./...

run: build
	@./bin/goon
