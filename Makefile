.DEFAULT_GOAL := build
THIS_FILE := $(lastword $(MAKEFILE_LIST))

GOPATH := $(shell go env | grep GOPATH | sed 's/GOPATH="\(.*\)"/\1/')
PATH := $(GOPATH)/bin:$(PATH)
export $(PATH)

BINARY=links
LD_FLAGS += -s -w

update-deps: fetch
	@echo -e "\n\033[0;36m [ Updating dependencies ]\033[0;m"
	$(GOPATH)/bin/govendor add +external
	$(GOPATH)/bin/govendor remove +unused
	$(GOPATH)/bin/govendor update +external

fetch:
	@echo -e "\n\033[0;36m [ Fetching dependencies ]\033[0;m"
	test -f $(GOPATH)/bin/govendor || go get -v github.com/kardianos/govendor
	test -f $(GOPATH)/bin/rice || go get -v github.com/GeertJohan/go.rice/rice
	$(GOPATH)/bin/govendor sync

clean:
	@echo -e "\n\033[0;36m [ Removing previously compiled binaries, and cleaning up ]\033[0;m"
	/bin/rm -fv "${BINARY}"
	/bin/rm -fv "rice-box.go"

build: clean fetch
	@echo -e "\n\033[0;36m [ Building ${BINARY} ]\033[0;m"
	$(GOPATH)/bin/rice -v embed-go
	go build -ldflags "${LD_FLAGS}" -i -x -v -o ${BINARY}
	# @$(MAKE) -f $(THIS_FILE) clean
