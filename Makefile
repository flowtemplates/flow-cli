.DEFAULT_GOAL := test

.PHONY: fmt
fmt:
	@gofumpt -w -l .
	@goimports -w -l .

.PHONY: lint
lint:
	@golangci-lint run

.PHONY: test
test:
	@gotestsum -f testdox

.PHONY: build
build:
	go build -o .out/flow ./cmd/flow
