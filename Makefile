build:
	@go build -o bin/api-blog ./cmd/main.go

run: build
	@./bin/api-blog
