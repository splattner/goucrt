# Set Shell to bash, otherwise some targets fail with dash/zsh etc.
SHELL := /bin/bash

# Disable built-in rules
MAKEFLAGS += --no-builtin-rules
MAKEFLAGS += --no-builtin-variables
.SUFFIXES:
.SECONDARY:
.DEFAULT_GOAL := help

.PHONY: help
help: ## Show this help
	@grep -E -h '\s##\s' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-20s\033[0m %s\n", $$1, $$2}'

# extensible array of targets. Modules can add target to this variable for the all-in-one target.
clean_targets := build-clean

PROJECT_ROOT_DIR = .
include Makefile.vars.mk

go_build ?= go build -o $(BIN_FILENAME) $(GOUCRT_MAIN_GO)

.PHONY: test
test: ## Run tests
	go test ./... -coverprofile cover.out

.PHONY: build
build: generate fmt vet $(BIN_FILENAME) docs-update-usage ## Build manager binary


.PHONY: fmt
fmt: ## Run go fmt against code
	go fmt ./cmd/ucrt

.PHONY: vet
vet: ## Run go vet against code
	go vet ./cmd/ucrt

.PHONY: lint
lint: fmt vet golangci-lint ## Invokes all linting targets
	@echo 'Check for uncommitted changes ...'
	git diff --exit-code

.PHONY: golangci-lint
golangci-lint: $(golangci_bin) ## Run golangci linters
	$(golangci_bin) run --timeout 5m --out-format colored-line-number ./...

.PHONY: docker-build
docker-build: $(BIN_FILENAME) ## Build the docker image
	docker build . \
	    -f build/Dockerfile \
		--tag $(GOUCRT_GHCR_IMG) \


.PHONY: docker-push
docker-push: ## Push the docker image
	docker push $(GOUCR_GHCR_IMG)

build-clean:
	rm -rf dist/ bin/ cover.out $(BIN_FILENAME) $(WORK_DIR)

clean: $(clean_targets) ## Cleans up all the locally generated resources


###
### Assets
###

# Build the binary without running generators
.PHONY: $(BIN_FILENAME)
$(BIN_FILENAME): export CGO_ENABLED = 0
$(BIN_FILENAME): export GOOS = $(GOUCRT_GOOS)
$(BIN_FILENAME): export GOARCH = $(GOUCRT_GOARCH)
$(BIN_FILENAME):
	$(go_build)

$(golangci_bin): | $(go_bin)
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b "$(go_bin)"