.DEFAULT_GOAL := build
BINARY=links

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

fetch: ## Fetches the necessary dependencies to build.
	which $(BIN)/rice 2>&1 > /dev/null || go install github.com/GeertJohan/go.rice/rice@latest
	go mod download
	go mod tidy

upgrade-deps: ## Upgrade all dependencies to the latest version.
	go get -u ./...

upgrade-deps-patch: ## Upgrade all dependencies to the latest patch release.
	go get -u=patch ./...

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "dist/" "${BINARY}" rice-box.go

generate: ## Generates the Go files that allow assets to be embedded.
	$(BIN)/rice -v embed-go

prepare: fetch clean generate ## Prepare the dependencies needed for a build.
	@echo

build: prepare ## Compile and generate a binary with static assets embedded.
	CGO_ENABLED=0 go build -ldflags '-d -s -w -extldflags=-static' -tags=netgo,osusergo,static_build -installsuffix netgo -buildvcs=false -trimpath  -o "${BINARY}"

debug: clean
	go run -v *.go --site-name "http://localhost:8080" --debug --http ":8080" --prom.enabled
