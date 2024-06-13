.DEFAULT_GOAL := help
.PHONY: help
help: ## help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.PHONY: lint
lint: ## run golang linter
	golangci-lint run

.PHONY: templ
templ: ## templ generate
	go run github.com/a-h/templ/cmd/templ@latest generate

.PHONY: parcel-build
parcel-build: ## build with parcel
	parcel build

.PHONY: npm-build
npm-build: ## build with npm
	npm install && npm run build

.PHONY: build
build: ## build binary
	go build -o ./bin/blog -v

.PHONY: run
run: lint templ parcel-build build ## run binary
	./bin/blog

.PHONY: build-all
build-all: npm-build templ ## build for prod
	go mod tidy -v && go build -o ./bin/blog -v