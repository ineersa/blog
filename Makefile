.DEFAULT_GOAL := help
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## run golang linter
	golangci-lint run

.PHONY: templ
templ:
	templ generate

.PHONY: parcel-build
parcel-build:
	parcel build

.PHONY: build
build:
	go build -o ./bin/blog -v

.PHONY: run
run: lint templ parcel-build build
	./bin/blog