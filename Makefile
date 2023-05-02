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

# Database variables
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_NAME ?= taskapp
DB_USERNAME ?= taskapp_user
DB_PASSWORD ?= password

ROOT_PACKAGE := github.com/gihyodocker/taskapp
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

.PHONY: setup-db-tools
setup-db-tools:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.15.2
	go install github.com/volatiletech/sqlboiler/v4@v4.14.2
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.14.2

.PHONY: migrate-schema-up
migrate-schema-up: setup-db-tools
	migrate -source file://./database -database mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp\($(DB_HOST):$(DB_PORT)\)/$(DB_NAME) up

migrate-schema-down: setup-db-tools
	migrate -source file://./database -database mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp\($(DB_HOST):$(DB_PORT)\)/$(DB_NAME) down

define SQLBOILER_CONFIG
pkgname="model"
output="pkg/model"
[mysql]
  dbname = "$(DB_NAME)"
  host   = "$(DB_HOST)"
  port   = $(DB_PORT)
  user   = "$(DB_USERNAME)"
  pass   = "$(DB_PASSWORD)"
  sslmode = "false"
  blacklist = ["schema_migrations"]
endef
export SQLBOILER_CONFIG

.PHONY: sqlboiler.toml
sqlboiler.toml:
	@echo "$$SQLBOILER_CONFIG" > $@

.PHONY: generate-db-model
generate-db-model: sqlboiler.toml
	@sqlboiler mysql --no-tests
	@rm sqlboiler.toml

.PHONY: build
build:
	$(eval GIT_COMMIT := $(shell git describe --tags --always))
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) \
		go build -ldflags "-s -w -X $(LDFLAG_VERSION)=$(GIT_COMMIT)" \
		-o ./bin/$* -mod=vendor main.go
