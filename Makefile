.PHONY: build clean test help default ci

BIN_NAME=cpma

default: build

help:
	@echo 'Management commands for cpma:'
	@echo
	@echo 'Usage:'
	@echo '    make build           Compile the project.'
	
	@echo '    make clean           Clean the directory tree.'
	@echo

build:
	@echo "GOPATH=${GOPATH}"
	go build -o bin/${BIN_NAME}

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

ci: build test

test:
	go test ./...

