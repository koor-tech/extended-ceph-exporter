SHELL := /usr/bin/env bash

PROJECTNAME ?= extended-ceph-exporter

GO111MODULE  ?= on
GO           ?= go
PREFIX       ?= $(shell pwd)
BIN_DIR      ?= $(PREFIX)/.bin
TARBALL_DIR  ?= $(PREFIX)/.tarball
PACKAGE_DIR  ?= $(PREFIX)/.package
ARCH         ?= amd64
PACKAGE_ARCH ?= linux-amd64

VERSION      := $(shell cat VERSION)
TOPDIR       := $(shell pwd)

# The GOHOSTARM and PROMU parts have been taken from the prometheus/promu repository
# which is licensed under Apache License 2.0 Copyright 2018 The Prometheus Authors
FIRST_GOPATH := $(firstword $(subst :, ,$(shell $(GO) env GOPATH)))

GOHOSTOS     ?= $(shell $(GO) env GOHOSTOS)
GOHOSTARCH   ?= $(shell $(GO) env GOHOSTARCH)

ifeq (arm, $(GOHOSTARCH))
	GOHOSTARM ?= $(shell GOARM= $(GO) env GOARM)
	GO_BUILD_PLATFORM ?= $(GOHOSTOS)-$(GOHOSTARCH)v$(GOHOSTARM)
else
	GO_BUILD_PLATFORM ?= $(GOHOSTOS)-$(GOHOSTARCH)
endif

PROMU_VERSION ?= 0.13.0
PROMU_URL     := https://github.com/prometheus/promu/releases/download/v$(PROMU_VERSION)/promu-$(PROMU_VERSION).$(GO_BUILD_PLATFORM).tar.gz

PROMU := $(FIRST_GOPATH)/bin/promu
# END copied code

pkgs = $(shell go list ./... | grep -v /vendor/ | grep -v /test/)

CONTAINER_IMAGE_NAME ?= docker.io/koorinc/extended-ceph-exporter
CONTAINER_IMAGE_TAG  ?= $(subst /,-,$(shell git rev-parse --abbrev-ref HEAD))

all: format style vet test build

build: promu
	@echo ">> building binaries"
	GO111MODULE=$(GO111MODULE) $(PROMU) build --prefix $(PREFIX)

check_license:
	@OUTPUT="$$($(PROMU) check licenses)"; \
	if [[ $$OUTPUT ]]; then \
		echo "Found go files without license header:"; \
		echo "$$OUTPUT"; \
		exit 1; \
	else \
		echo "All files with license header"; \
	fi

crossbuild: promu
	$(PROMU) crossbuild
	$(PROMU) crossbuild tarballs

docker: crossbuild
	$(MAKE) docker-build

docker-build:
	@echo ">> building docker image"
	docker build \
		--build-arg BUILD_DATE="$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')" \
		--build-arg REVISION="$(shell git rev-parse HEAD)" \
		-t "$(CONTAINER_IMAGE_NAME):$(CONTAINER_IMAGE_TAG)" \
		.
	docker tag "$(CONTAINER_IMAGE_NAME):$(CONTAINER_IMAGE_TAG)" "$(CONTAINER_IMAGE_NAME):latest"

docker-publish:
	docker push "$(CONTAINER_IMAGE_NAME):$(CONTAINER_IMAGE_TAG)"
	docker push "$(CONTAINER_IMAGE_NAME):latest"

format:
	go fmt $(pkgs)

helm-doc:
	GO111MODULE=on go install github.com/norwoodj/helm-docs/cmd/helm-docs@v1.11.0
	helm-docs --chart-search-root=./charts

promu:
	$(eval PROMU_TMP := $(shell mktemp -d))
	curl -s -L $(PROMU_URL) | tar -xvzf - -C $(PROMU_TMP)
	mkdir -p $(FIRST_GOPATH)/bin
	cp $(PROMU_TMP)/promu-$(PROMU_VERSION).$(GO_BUILD_PLATFORM)/promu $(FIRST_GOPATH)/bin/promu
	rm -r $(PROMU_TMP)

style:
	@echo ">> checking code style"
	@! gofmt -d $(shell find . -path ./vendor -prune -o -name '*.go' -print) | grep '^'

tarball: tree                                                                                                                                       
	@echo ">> building release tarball"
	@$(PROMU) tarball --prefix $(TARBALL_DIR) $(BIN_DIR)

clean:
	rm -rf $(PROJECTNAME) $(PROJECTNAME).spec $(PROJECTNAME)-$(VERSION).tar.gz 
	
test:
	@$(GO) test $(pkgs)

test-short:
	@echo ">> running short tests"
	@$(GO) test -short $(pkgs)

vet:
	@echo ">> vetting code"
	@$(GO) vet $(pkgs)

.PHONY: all build crossbuild docker docker-publish format promu style tarball test test-short vet
