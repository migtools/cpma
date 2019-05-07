.PHONY: build clean test help default ci

BIN_NAME=cpma
SOURCES:=$(shell find . -name '*.go')
SOURCE_DIRS=cmd pkg internal

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
	GO111MODULE=on go build -o bin/${BIN_NAME}

clean:
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

ci: build test

test:
	GO111MODULE=on go test ./...

lint: ## Run golint
	@golint -set_exit_status $(addsuffix /... , $(SOURCE_DIRS))

fmt: ## Run go fmt
	@gofmt -d $(SOURCES)

fmtcheck: ## Check go formatting
	@gofmt -l $(SOURCES) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi
