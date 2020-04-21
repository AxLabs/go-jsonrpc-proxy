ORG_NAME := AxLabs
REPO_NAME := go-jsonrpc-proxy
PKG_ROOT := github.com/${ORG_NAME}/$(REPO_NAME)
PKG_LIST := go list ${PKG_ROOT}/...
PKG_GO_JSONRPC_PROXY := ${PKG_ROOT}/cmd
CMD_DIR := ./cmd

.PHONY: all lint vet test go-jsonrpc-proxy

.EXPORT_ALL_VARIABLES:

GO111MODULE=on

all: lint vet test go-jsonrpc-proxy

vet:
	@go vet $(shell $(PKG_LIST))

# Lint the files
lint:
	@golint -set_exit_status $(shell $(PKG_LIST))

# Run unit tests
test:
	@go test -v -short -count=1 $(shell $(PKG_LIST))

go-jsonrpc-proxy: $(CMD_DIR)/go-jsonrpc-proxy

$(CMD_DIR)/go-jsonrpc-proxy:
	@echo "Building $@..."
	@go build -i -o $(CMD_DIR)/go-jsonrpc-proxy -v $(PKG_GO_JSONRPC_PROXY)
	@chmod u+x $(CMD_DIR)/go-jsonrpc-proxy