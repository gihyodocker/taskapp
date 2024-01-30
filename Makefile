ifndef GOARCH
	GOARCH=$(shell go env GOARCH)
endif

ifndef GOOS
	GOOS := $(shell go env GOOS)
endif

COMMAND_DIRS := $(wildcard cmd/*)
BUILD_TARGETS := $(addprefix build-,$(notdir $(COMMAND_DIRS)))
MYSQL_PASSWORD_PATH := ./secrets/mysql_user_password

# Database variables
DB_HOST ?= localhost
DB_PORT ?= 3306
DB_NAME ?= taskapp
DB_USERNAME ?= taskapp_user
DB_PASSWORD := $(shell if [ -f $(MYSQL_PASSWORD_PATH) ]; then cat $(MYSQL_PASSWORD_PATH); else echo "password"; fi )

ROOT_PACKAGE := github.com/gihyodocker/taskapp

.PHONY: install-tools
install-tools:
	@sh hack/install-tools.sh

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GO111MODULE=on go mod vendor

.PHONY: mod
mod:
	GO111MODULE=on go mod download 

.PHONY: setup-db-tools
setup-db-tools:
	go install -tags 'mysql' github.com/golang-migrate/migrate/v4/cmd/migrate@v4.17.0
	go install github.com/volatiletech/sqlboiler/v4@v4.15.0
	go install github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-mysql@v4.15.0

.PHONY: migrate-schema-up
migrate-schema-up: setup-db-tools
	migrate -path ./containers/migrator/history -database mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp\($(DB_HOST):$(DB_PORT)\)/$(DB_NAME) up

migrate-schema-down: setup-db-tools
	migrate -path ./containers/migrator/history -database mysql://$(DB_USERNAME):$(DB_PASSWORD)@tcp\($(DB_HOST):$(DB_PORT)\)/$(DB_NAME) down

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

.PHONY: $(BUILD_TARGETS)
$(BUILD_TARGETS): build-%:
	CGO_ENABLED=0 GO111MODULE=on GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o ./bin/$* -mod=vendor cmd/$*/main.go

.PHONY: make-mysql-passwords
make-mysql-passwords:
	@mkdir -p ./secrets
	@go run cmd/tools/main.go mysql generate-password
	@cp -R ./secrets ./k8s/kustomize/base/mysql/
	@cp -R ./secrets ./k8s/kustomize/base/migrator/

.PHONY: api-config-local.yaml
api-config-local.yaml:
	@go run cmd/api/main.go config \
		--database-password $(DB_PASSWORD) \
		--output-file ./api-config-local.yaml

.PHONY: api-config.yaml
api-config.yaml:
	@go run cmd/api/main.go config \
		--database-host mysql \
		--database-password $(DB_PASSWORD) \
		--output-file ./api-config.yaml
	@cp ./api-config.yaml ./k8s/kustomize/base/api/

.PHONY: serve-api
serve-api:
	@go run cmd/api/main.go server \
		--config-file ./api-config-local.yaml

.PHONY: serve-web
serve-web:
	@go run cmd/web/main.go server \
		--assets-dir $(PWD)/assets \

.PHONY: make-k8s-mysql-secret
make-k8s-mysql-secret:
	@kubectl create secret generic mysql --dry-run=client -o yaml \
		--from-literal=root_password=$(shell cat ./secrets/mysql_root_password) \
		--from-literal=user_password=$(DB_PASSWORD) > ./k8s/plain/local/mysql-secret.yaml
	@cp ./k8s/plain/local/mysql-secret.yaml ./k8s/plain/aks/mysql-secret.yaml

.PHONY: make-k8s-api-config
make-k8s-api-config:
	@kubectl create secret generic api-config --dry-run=client -o yaml \
		--from-file=api-config.yaml=./api-config.yaml > ./k8s/plain/local/api-config-secret.yaml
	@cp ./k8s/plain/local/api-config-secret.yaml ./k8s/plain/aks/api-config-secret.yaml
