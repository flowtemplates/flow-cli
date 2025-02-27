.DEFAULT_GOAL := test

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
	go build -o .out/bin cmd/cli/main.go 

.PHONY: run
run:
	go run cmd/cli/main.go
