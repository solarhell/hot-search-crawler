.PHONY: build test run fmt vet

build:
	go build -o bin/hot-search-crawler ./cmd/hot-search-crawler

test:
	go test ./...

run:
	go run ./cmd/hot-search-crawler crawl --update-readme

fmt:
	gofmt -w .

vet:
	go vet ./...
