VERSION="0.0.1"
GIT_COMMIT=$(git rev-parse --short HEAD)
BUILD_TIME=$(date -u +'%Y-%m-%dT%H:%M:%SZ')

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
	go build -o cmd/structcopy-gen/structcopy-gen -ldflags "-X 'github.com/bookweb/structcopy-gen/config.Version=${VERSION}'" cmd/structcopy-gen/main.go

install:
	cd cmd/structcopy-gen && go install

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
