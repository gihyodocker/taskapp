# TODO: Remove when this repository is published
ifndef GOPRIVATE
	GOPRIVATE="github.com/gihyodocker"
endif

ifndef GOARCH
	GOARCH=$(shell go env GOARCH)
endif

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif

ROOT_PACKAGE := github.com/gihyodocker/todoapp
VERSION_PACKAGE := $(ROOT_PACKAGE)/pkg/version
LDFLAG_VERSION := $(VERSION_PACKAGE).version

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GOPRIVATE=$(GOPRIVATE) GO111MODULE=on go mod vendor

.PHONY: mod
mod:
	GOPRIVATE=$(GOPRIVATE) GO111MODULE=on go mod download 

.PHONY: build
build:
	$(eval GIT_COMMIT := $(shell git describe --tags --always))
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "-s -w -X $(LDFLAG_VERSION)=$(GIT_COMMIT)" \
		-o ./bin/$* -mod=vendor main.go