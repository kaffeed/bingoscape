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


# run templ generation in watch mode to detect all .templ files and 
# re-create _templ.txt files on change, then send reload event to browser. 
# Default url: http://localhost:7331
#
templ:
	@go run github.com/a-h/templ/cmd/templ@latest generate --watch --proxy="http://localhost$(HTTP_LISTEN_ADDR)" --open-browser=false -v

# run air to detect any go file changes to re-build and re-run the server.
server:
	@go run github.com/air-verse/air@latest \
	--build.cmd "go build --tags dev -o tmp/bin/main ./cmd/app/" --build.bin "tmp/bin/main" --build.delay "100" \
	--build.exclude_dir "node_modules" \
	--build.include_ext "go" \
	--build.stop_on_error "false" \
	--misc.clean_on_exit true

# run tailwindcss to generate the styles.css bundle in watch mode.
watch-assets:
	@npx tailwindcss -i app/assets/app.css -o ./public/assets/styles.css --watch   

# # run esbuild to generate the index.js bundle in watch mode.
# watch-esbuild:
# 	npx esbuild app/assets/index.js --bundle --outdir=public/assets --watch

# watch for any js or css change in the assets/ folder, then reload the browser via templ proxy.
sync_assets:
	@go run github.com/air-verse/air@latest \
	--build.cmd "go run github.com/a-h/templ/cmd/templ@latest generate --notify-proxy" \
	--build.bin "true" \
	--build.delay "100" \
	--build.exclude_dir "" \
	--build.include_dir "public" \
	--build.include_ext "js,css"

# build the application for production. This will compile your app
# to a single binary with all its assets embedded.
#
generate-code:
	@go run github.com/a-h/templ/cmd/templ@latest generate -v
	@npx tailwindcss -i app/assets/app.css -o ./public/assets/styles.css

build:
	@make generate-code
	@go build -o bin/bingoscape cmd/app/main.go
	@echo "compiled you application with all its assets to a single binary => bin/bingoscape"

# start the application in development
dev:
	@make -j5 templ server watch-assets sync_assets

## test: run unit tests
.PHONY: test
test:
	go test -race -cover $(PACKAGES)

## db-status: run database-status
.PHONY: db-status
db-status:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_URL) go run github.com/pressly/goose/v3/cmd/goose@latest status

## db-reset: run database-reset
.PHONY: db-reset
db-reset:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_URL) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) reset

## db-down: run migrate database-down
.PHONY: db-down
db-down:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_URL) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) down

## db-up: migrate database-up
.PHONY: db-up
db-up:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_URL) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) up

## db-mig-create: create new migration
.PHONY: db-mig-create
db-mig-create:
	@GOOSE_DRIVER=$(DB_DRIVER) GOOSE_DBSTRING=$(DB_URL) go run github.com/pressly/goose/v3/cmd/goose@latest -dir=$(MIGRATION_DIR) -s create $(filter-out $@,$(MAKECMDGOALS)) sql

