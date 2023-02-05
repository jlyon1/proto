TIMESTAMP ?= $(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
SHA ?= $(shell git rev-parse --short HEAD)
TAG ?= $(shell git describe --tags 2> /dev/null)

CURRENT_DIR = $(shell pwd -P)
DIST_DIR = ./dist

.PHONY: help

help: ## this list
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

build: ## Build the output binary
	@go build -o proto