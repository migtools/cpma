.PHONY: build clean test help default ci

BIN_NAME=cpma
SOURCES:=$(shell find . -name '*.go')
SOURCE_DIRS=cmd pkg

default: build

help: ## Show this help screen
	@echo 'Usage: make <OPTIONS> ... <TARGETS>'
	@echo ''
	@echo 'Available targets are:'
	@echo ''
	@grep -E '^[ a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'
	@echo ''

build: ## Compile the project
	@echo "GOPATH=${GOPATH}"
	GO111MODULE=on go build -o bin/${BIN_NAME}

clean: ## Clean the directory tree
	@test ! -e bin/${BIN_NAME} || rm bin/${BIN_NAME}

ci: build test

test: ## Test the project
	GO111MODULE=on go test ./...

lint: ## Run golint
	@golint -set_exit_status $(addsuffix /... , $(SOURCE_DIRS))

fmt: ## Run go fmt
	@gofmt -d $(SOURCES)

fmtcheck: ## Check go formatting
	@gofmt -l $(SOURCES) | grep ".*\.go"; if [ "$$?" = "0" ]; then exit 1; fi

vet: ## Run go vet
	@GO111MODULE=on go vet $(addsuffix /..., $(addprefix ./, $(SOURCE_DIRS)))
