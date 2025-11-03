VERSION := $(shell cat .VERSION)
GIT_COMMIT := $(shell git rev-parse --short HEAD)
BUILD_TIME := $(shell date -u +'%Y-%m-%dT%H:%M:%S%NZ')

.PHONY: help
help: ## Show this help message.
	@echo 'usage: make [target] ...'
	@echo
	@echo 'targets:'
	@egrep '^(.+)\:\ ##\ (.+)' ${MAKEFILE_LIST} | column -t -c 2 -s ':#'

.PHONY: install-deps build install
install-deps:
# 	run `make install-tool` on pm-go-tools

mod:
	go mod tidy
	go mod vendor

generate:
	go generate ./...

build:
	go build -ldflags "-X github.com/bookweb/structcopy-gen/config.Version=$(VERSION) \
		-X github.com/bookweb/structcopy-gen/config.CommitHash=$(GIT_COMMIT) \
		-X github.com/bookweb/structcopy-gen/config.BuildTime=${BUILD_TIME}" \
		-o ./cmd/structcopy-gen/structcopy-gen ./cmd/structcopy-gen/main.go

install:
	go install -ldflags "-X github.com/bookweb/structcopy-gen/config.Version=$(VERSION) \
		-X github.com/bookweb/structcopy-gen/config.CommitHash=$(GIT_COMMIT) \
		-X github.com/bookweb/structcopy-gen/config.BuildTime=${BUILD_TIME}" \
		./cmd/structcopy-gen

generate-standalone:
	go generate ./examples/internal/standalone
	
generate-example:
	go generate ./examples/internal/example

.PHONY: lint test coverage
lint:
	golangci-lint run

test:
	go test github.com/bookweb/structcopy-gen/tests && \
	go test github.com/bookweb/structcopy-gen/internal/gen/...

coverage:
	@go test -v -cover ./... -coverprofile coverage.out -coverpkg ./... 2>&1 >/dev/null && \
	go tool cover -func coverage.out -o coverage.out 2>&1 >/dev/null && \
	cat coverage.out

.PHONY: tag release
tag:
	autotag -b master > .VERSION

tag-dev:
	autotag -p dev -b develop > .VERSION

tag-stg:
	autotag -p next -b release-next > .VERSION

tag-first:
	git tag v0.0.1 -m'create project'

release-init:
	goreleaser init

release-snapshot:
	goreleaser release --snapshot --clean

release:
	goreleaser release --clean
