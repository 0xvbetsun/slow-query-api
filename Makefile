cnf ?= ./deployments/.env
include $(cnf)
export $(shell sed 's/=.*//' $(cnf))

.PHONY: help

help:
	@awk 'BEGIN {FS = ":.*?## "} /^[a-zA-Z_-]+:.*?## / {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help

app: ## Run application
	go run ./cmd/app/main.go

test: ## Runs test for app
	go test ./... -coverprofile=coverage.out -covermode=atomic

cover: ## Gets percents of code coverage
	go tool cover -func cover.out | grep total:
