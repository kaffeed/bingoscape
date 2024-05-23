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

## test: run unit tests
.PHONY: test
test:
	go test -race -cover $(PACKAGES)

## build: build a binary
.PHONY: build
build: test
	go build  -o ./tmp/main -v ./cmd 

## docker-build: build project into a docker container image
.PHONY: docker-build
docker-build: test
	GOPROXY=direct docker buildx build -t ${name} .

## docker-run: run project in a container
.PHONY: docker-run
docker-run:
	docker run -it --rm -p 8080:8080 ${name}

## start: build and run local project
.PHONY: start
start: build
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
	migrate -path db/migrations/ -database ${DB} -verbose up

## migration-down: run migrate database-down
.PHONY: migration-down
migration-down: 
	migrate -path db/migrations/ -database ${DB} -verbose down

## migration-fix: force db version
.PHONY: migration-fix
migration-fix: 
	migrate -path db/migrations/ -database "${DB}" force ${version}
