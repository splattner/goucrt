IMG_TAG ?= latest

GOUCRT_MAIN_GO ?= cmd/ucrt/ucrt.go
GOUCRT__GOOS ?= linux
GOUCRT__GOARCH ?= amd64

CURDIR ?= $(shell pwd)
BIN_FILENAME ?= $(CURDIR)/$(PROJECT_ROOT_DIR)/ucrt
WORK_DIR = $(CURDIR)/.work

go_bin ?= $(PWD)/.work/bin
$(go_bin):
	@mkdir -p $@

golangci_bin = $(go_bin)/golangci-lint


# Image URL to use all building/pushing image targets
GOUCRT_GHCR_IMG ?= ghcr.io/splattner/gocurt:$(IMG_TAG)

