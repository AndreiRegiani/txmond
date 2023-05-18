.PHONY: build clean lint ci-lint build-docker

BINARY = txmond
DOCKER_IMAGE = txmond:latest

build:
	go build -o $(BINARY) cmd/txmond.go

clean:
	go clean
	rm -f $(BINARY)

lint:
	golint  ./cmd/...

ci-lint:
	golangci-lint run ./cmd/...

build-docker:
	docker build -t $(DOCKER_IMAGE) .
