.DEFAULT_GOAL := test

DOCKER_IMAGE = olytix-server
DOCKER_TAG = latest
PORT = 8080


.PHONY: fmt
fmt:
	@gofmt -s -w -l .
	@goimports -w -l .

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@go test -v ./... | { grep -v 'no test files'; true; }

.PHONY: build
build:
	docker build -t $(DOCKER_IMAGE):$(DOCKER_TAG) .

.PHONY: run
run:
	docker run -p $(PORT):$(PORT) $(DOCKER_IMAGE):$(DOCKER_TAG)

.PHONY: help
help:
	docker run $(DOCKER_IMAGE):$(DOCKER_TAG) -h

.PHONY: swag
swag:
	@swag fmt
	@swag init -g internal/app/app.go

.PHONY: mockery
mockery:
	@mockery
