SHELL := /bin/bash
.SHELLFLAGS = -e -c
.DEFAULT_GOAL := help
.ONESHELL:

db_data: ## Generate the db directory
	mkdir -p $(shell pwd)/db_data

.PHONY: db
db:
	docker compose up -d

.PHONY: run
run: db_data db ## Run the server application
	go run $(shell pwd)/cmd/main.go

UNAME_S = $(shell uname -s)

.PHONY: serve/client
serve: ## Serve the client application and open default browser
	python3 -m http.server -d wwwroot/ &
ifeq ($(UNAME_S), Linux)
	xdg-open http://0.0.0.0:8000
else ifeq ($(UNAME_S), Darwin)
	open http://0.0.0.0:8000
endif

.PHONY: help
help: ## Displays help info
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 } /^##@/ { printf "\n\033[1m%s\033[0m\n", substr($$0, 5) } ' $(MAKEFILE_LIST)
