.DEFAULT_GOAL := build
THIS_FILE := $(lastword $(MAKEFILE_LIST))

GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)

BINARY=links
LD_FLAGS += -s -w
VERSION=$(shell git describe --tags --abbrev=0 2>/dev/null | sed -r "s:^v::g")
RSRC=README_TPL.md
ROUT=README.md

readme-gen:
	cp -av "${RSRC}" "${ROUT}"
	sed -ri -e "s:\[\[tag\]\]:${VERSION}:g" -e "s:\[\[os\]\]:linux:g" -e "s:\[\[arch\]\]:amd64:g" "${ROUT}"

release: clean fetch generate
	$(GOPATH)/bin/goreleaser --skip-publish

publish: clean fetch generate
	$(GOPATH)/bin/goreleaser

snapshot: clean fetch generate
	$(GOPATH)/bin/goreleaser --snapshot --skip-validate --skip-publish

update-deps: fetch
	$(GOPATH)/bin/govendor add +external
	$(GOPATH)/bin/govendor remove +unused
	$(GOPATH)/bin/govendor update +external

fetch:
	test -f $(GOPATH)/bin/govendor || go get -u -v github.com/kardianos/govendor
	test -f $(GOPATH)/bin/goreleaser || go get -u -v github.com/goreleaser/goreleaser
	test -f $(GOPATH)/bin/rice || go get -u -v github.com/GeertJohan/go.rice/rice
	$(GOPATH)/bin/govendor sync

clean:
	/bin/rm -rfv "dist/" ${BINARY} rice-box.go

generate:
	$(GOPATH)/bin/rice -v embed-go

compress:
	(which /usr/bin/upx > /dev/null && find dist/*/* | xargs -I{} -n1 -P 4 /usr/bin/upx --best "{}") || echo "not using upx for binary compression"

build: fetch generate
	go build -ldflags "${LD_FLAGS}" -i -v -o ${BINARY}
