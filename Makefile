# TODO: Remove when this repository is published
ifndef GOPRIVATE
	GOPRIVATE="github.com/gihyodocker"
endif

.PHONY: tidy
tidy:
	GO111MODULE=on go mod tidy

.PHONY: vendor
vendor:
	GOPRIVATE=$(GOPRIVATE) GO111MODULE=on go mod vendor

.PHONY: mod
mod:
	GOPRIVATE=$(GOPRIVATE) GO111MODULE=on go mod download 
