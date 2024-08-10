build:
	@go build -o ./bin/medspace

run: build
	@./bin/medspace


test:
	@go test ./...