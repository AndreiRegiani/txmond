.PHONY: build clean lint build-docker

GO_FILES = $(wildcard *.go **/*.go)
BINARY = txmond
DOCKER_IMAGE = txmond:latest

build:
	go build -o $(BINARY) cmd/txmond.go

clean:
	go clean
	rm -f $(BINARY)

lint:
	golint $(GO_FILES)

build-docker:
	docker build -t $(DOCKER_IMAGE) .
