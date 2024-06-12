ifneq (,$(wildcard ./.env))
    include .env
    export
endif

PACKAGES := $(shell go list ./...)
name := $(shell basename ${PWD})

all: help

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a make command to run"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

## init: initialize project (make init module=github.com/user/project)
.PHONY: init
init:
	go mod init ${module}
	go install github.com/cosmtrek/air@latest
	go install github.com/pressly/goose/v3/cmd/goose@latest
	go install github.com/a-h/templ/cmd/templ@latest
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

## test: run unit tests
.PHONY: test
test:
	go test -race -cover $(PACKAGES)

## build-api: build a binary
.PHONY: build-api
build-api: test
	go build  -o ./tmp/main -v ./cmd/api

## docker-build: build project into a docker container image
.PHONY: docker-build
docker-build: test
	GOPROXY=direct docker buildx build -t ${name} .


## build-release: build a binary
.PHONY: build-release
build-release: test css templ
	go build  -o ./tmp/main -v ./cmd/api

## docker-run: run project in a container
.PHONY: docker-run
docker-run:
	docker run -it --rm -p 8080:8080 ${name}

## start-api: build and run local project
.PHONY: start-api
start-api: build-api
	air

## css: build tailwindcss
.PHONY: css
css:
	npx tailwindcss -i assets/css/input.css -o assets/css/styles.css --minify

## css-watch: watch build tailwindcss
.PHONY: css-watch
css-watch:
	npx tailwindcss -i assets/css/input.css -o assets/css/styles.css --watch

## templ: generate templ
.PHONY: templ
templ:
	templ generate

## templ-watch: generate templ and watch for changes
.PHONY: templ-watch
templ-watch:
	templ generate --watch

## migration-up: migrate database-up
.PHONY: migration-up
migration-up: 
	goose -dir db/migrations/ postgres "${DB}" up

## migration-down: run migrate database-down
.PHONY: migration-down
migration-down: 
	goose -dir db/migrations/ postgres "${DB}" down

## migration-new: force db version
.PHONY: migration-new
migration-new: 
	goose -dir db/migrations/ -s create "${migration}" sql

## db-version: check database version
.PHONY: db-version
db-version: 
	goose -dir db/migrations/ postgres "${DB}" version
