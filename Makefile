.DEFAULT_GOAL := build

DIRS=bin
BINARY=links

VERSION=$(shell git describe --tags --always --abbrev=0 --match=v* 2> /dev/null | sed -r "s:^v::g" || echo 0)

$(info $(shell mkdir -p $(DIRS)))
BIN=$(CURDIR)/bin
export GOBIN=$(CURDIR)/bin


help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	which $(BIN)/rice 2>&1 > /dev/null || go get -v github.com/GeertJohan/go.rice/rice
	go mod download
	go mod tidy
	go mod vendor

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "dist/" "${BINARY}" rice-box.go

generate: ## Generates the Go files that allow assets to be embedded.
	$(BIN)/rice -v embed-go

build: fetch clean generate ## Compile and generate a binary with static assets embedded.
	go build -ldflags '-d -s -w' -tags netgo -installsuffix netgo -v -o "${BINARY}"

debug: clean
	go run -v *.go --site-name "http://0.0.0.0:8080" --debug --http ":8080"
