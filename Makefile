.PHONY: build clean test help default ci

BIN_NAME=cpma
SOURCES:=$(shell find . -name '*.go')
SOURCE_DIRS=cmd pkg
DATE:=`date -u +%Y/%m/%d.%H:%M:%S`
VERSION:=`git describe --tags --always --long --dirty`
LDFLAGS=-ldflags "-X=github.com/fusor/cpma/cmd.BuildVersion=$(VERSION) -X=github.com/fusor/cpma/cmd.BuildTime=$(DATE)"

default: build

help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

build: 	bundle ## Compile the project
	@echo "GOPATH=${GOPATH}"
	GO111MODULE=on go build $(LDFLAGS) -o bin/${BIN_NAME}

clean: ## Clean the directory tree
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

ci: bundle lint fmtcheck vet build test

cover: ## Project test coverage and generate covergate html file
	GO111MODULE=on go test -cover -covermode=count -coverprofile=coverage.out ./pkg/... ./cmd/... \
	&& go tool cover -html=coverage.out -o coverage.html

test: ## Test the project
	GO111MODULE=on go test ./pkg/... ./cmd/...

lint: ## Run golint
	@golint -set_exit_status $(addsuffix /... , $(SOURCE_DIRS))

fmt: ## Run go fmt
	@gofmt -d $(SOURCES)

fmtcheck: ## Check go formatting
	@gofmt -l $(SOURCES) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

vet: ## Run go vet
	@GO111MODULE=on go vet $(addsuffix /..., $(addprefix ./, $(SOURCE_DIRS)))

e2e: ## Execute e2e test
	GO111MODULE=on go test ./test/e2e/...

bundle: # Bundle files for html reports
	cd pkg/transform/reportoutput/ && go generate && cd ../../..

