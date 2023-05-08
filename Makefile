ADDR="http://localhost:8080"

PROJECTNAME=$(shell basename "$(PWD)")

# Go related variables
GOBASE=$(shell pwd)
GOBIN=$(GOBASE)/bin

# Redirect error output to file, so we can show it in development mode
STDERR=/tmp/.$(PROJECTNAME)-stderr.txt

# PID file will keep the process id of the server
PID=/tmp/.$(PROJECTNAME).pid

# Make is verbose in Linux. Make it silent
# MAKEFLAGS += --silient

# all: Prints out the instructions
all: help

## start: Start in development mode. Auto starts when code changes.
start:
	$(export APP_ENVIRONMENT="local") CGO_ENABLED=0 go run ./cmd/server/main.go

build:
	@GOPATH=$(GOPATH) GOBIN=$(GOBIN) CGO_ENABLED=0 go build -o $(GOBIN)/$(PROJECTNAME) ./cmd/server/main.go

# Read more here:`https://kodfabrik.com/journal/a-good-makefile-for-go/`
# Stop development mode
stop: stop-server

start-server:
	@echo "  >  $(PROJECTNAME) is available at $(ADDR)"
	@-$(GOBIN)/$(PROJECTNAME) 2>&1 & echo $$! > $(PID)
	@cat $(PID) | sed "/^/s/^/  \>  PID: /"

stop-server:
	@-touch $(PID)
	@-kill `cat $(PID)` 2> /dev/null || true
	@-rm $(PID)

restart-server: stop-server start-server

## clean: Clean build files. Runs `go clean` internally
clean:
	@echo "	> Cleaning build cache"
	go clean

## test: Run unit test
test:
	@echo "	> Run unit test"
	go test ./...

## test-coverage: Run unit test with coverage
test-coverage:
	@echo "	> Run unit test with coverage"
	go test ./... -coverprofile=coverage.out

## vet: Run go vet interally
vet:
	@echo "	> Run go vet"
	go vet

## dep: Download dependencies
dep:
	@echo "	> Download dependencies"
	go mod download

## tidy: Run go mod tidy internally
dep:
	@echo "	> Run go mod tidy"
	go mod tidy

## deploy-dev: Deploy app in dev mode with docker compose
deploy-dev: stop-deploy-dev
	@echo "	> Build docker image"
	@docker build -t $(PROJECTNAME):v0.0.1 .
	@echo "	> Launch app with docker compose"
	@docker-compose -f ./docker-compose.yml up

## stop-deploy-dev: Stop app in dev mode by running command `docker-compose down`
stop-deploy-dev:
	@echo "	> Stop all container in dev mode"
	@docker-compose -f ./docker-compose.yml down

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo

