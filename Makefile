GOCMD=go
GOBUILD=$(GOCMD) build
GOGET=$(GOCMD) get
GODEP=dep
BINARY_NAME=ocinventory
BINARY_UNIX=$(BINARY_NAME)_unix
DOCKER_WORK_DIR=/go/src/github.com/literalice/openshift-inventory-utils

.PHONY: deps build docker-build

deps:
	$(GOGET) -u github.com/golang/dep/cmd/dep
	$(GODEP) ensure
build: deps
	mkdir -p $(CURDIR)/.bin
	$(GOBUILD) -o $(CURDIR)/.bin/$(BINARY_NAME) -v ./cmd/ocinventory
docker-build:
	docker run --rm -v "$(CURDIR)":$(DOCKER_WORK_DIR) -w $(DOCKER_WORK_DIR) golang:1.11.2 make build BINARY_NAME=$(BINARY_UNIX)