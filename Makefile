default: help

.PHONY: help
help: ## list makefile targets
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

PHONY: test
test: ## run go tests
	go test -v ./...

PHONY: cover
cover: ## display test coverage
	go test -v -race $(shell go list ./... | grep -v /vendor/) -v -coverprofile=coverage.out
	go tool cover -func=coverage.out

PHONY: fmt
fmt: ## format go files
	gofmt -w -s  .

PHONY: lint
lint: ## lint go files
	golangci-lint run -c .golang-ci.yml

PHONY: build
build: ## build go binary
	go build -o apastat .

PHONY: build-linux
build-linux: ## build go binary for linux
	GOOS=linux GOARCH=amd64 go build -o apastat-linux .
