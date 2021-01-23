VERSION     = v0.6.x
LDFLAGS     = -ldflags "-X uhppote.VERSION=$(VERSION)" 
DIST       ?= development
CLI         = ./bin/uhppoted-app-wild-apricot

DATETIME  = $(shell date "+%Y-%m-%d %H:%M:%S")
DEBUG    ?= --debug

.PHONY: bump

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -o bin ./...

test: build
	go test ./...

vet: build
	go vet ./...

lint: build
	golint ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

build-all: test vet
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64         go build -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 go build -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64         go build -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64         go build -o dist/$(DIST)/windows ./...

release: build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)
	cd dist; zip --recurse-paths $(DIST).zip $(DIST)

bump:
	go get -u github.com/uhppoted/uhppote-core
	go get -u github.com/uhppoted/uhppoted-api

debug: build
	$(CLI) --debug get

# GENERAL COMMANDS

usage: build
	$(CLI)

help: build
	$(CLI) help

version: build
	$(CLI) version

# ACL COMMANDS

get: build
	$(CLI) --debug get
