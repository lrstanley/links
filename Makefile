.DEFAULT_GOAL := build

GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)

BINARY=links
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null | sed -r "s:^v::g")
COMPRESS_CONC ?= $(shell nproc)
RSRC=README_TPL.md
ROUT=README.md

help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-12s\033[0m %s\n", $$1, $$2}'

update-deps: fetch ## Adds any missing dependencies, removes unused deps, etc.
	$(GOPATH)/bin/govendor add -v +external
	$(GOPATH)/bin/govendor remove -v +unused
	$(GOPATH)/bin/govendor update -v +vendor

upgrade-deps: update-deps ## Upgrades all dependencies to the latest available versions and saves them.
	$(GOPATH)/bin/govendor fetch -v +vendor

readme-gen: ## Generates readme from template file.
	cp -av "${RSRC}" "${ROUT}"
	sed -ri -e "s:\[\[tag\]\]:${VERSION}:g" -e "s:\[\[os\]\]:linux:g" -e "s:\[\[arch\]\]:amd64:g" "${ROUT}"

snapshot: clean fetch generate ## Generate a snapshot release.
	$(GOPATH)/bin/goreleaser --snapshot --skip-validate --skip-publish

release: clean fetch generate ## Generate a release, but don't publish to GitHub.
	$(GOPATH)/bin/goreleaser --skip-validate --skip-publish

publish: clean fetch generate ## Generate a release, and publish to GitHub.
	$(GOPATH)/bin/goreleaser

fetch: ## Fetches the necessary dependencies to build.
	test -f $(GOPATH)/bin/govendor || go get -u -v github.com/kardianos/govendor
	test -f $(GOPATH)/bin/goreleaser || go get -u -v github.com/goreleaser/goreleaser
	test -f $(GOPATH)/bin/rice || go get -u -v github.com/GeertJohan/go.rice/rice
	$(GOPATH)/bin/govendor sync

clean: ## Cleans up generated files/folders from the build.
	/bin/rm -rfv "dist/" "${BINARY}" rice-box.go

generate: ## Generates the Go files that allow assets to be embedded.
	$(GOPATH)/bin/rice -v embed-go

compress: ## Uses upx to compress release binaries (if installed, uses all cores/parallel comp.)
	(which upx > /dev/null && find dist/*/* | xargs -I{} -n1 -P ${COMPRESS_CONC} upx --best "{}") || echo "not using upx for binary compression"

build: fetch generate ## Compile and generate a binary with static assets embedded.
	go build -ldflags '-d -s -w' -tags netgo -installsuffix netgo -v -x -o "${BINARY}"

debug:
	go run -x -v *.go --site-name "http://localhost:8080" --debug --http ":8080"
