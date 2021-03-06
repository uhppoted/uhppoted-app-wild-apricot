DIST       ?= development
CLI         = ./bin/uhppoted-app-wild-apricot
WORKDIR     = ../runtime/wild-apricot
CREDENTIALS = $(WORKDIR)/.credentials.json
RULES       = $(WORKDIR)/wild-apricot.grl

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
	go get -u github.com/uhppoted/uhppoted-lib

debug: build
	$(CLI) load-acl    --credentials $(CREDENTIALS) --rules $(RULES) --force --dry-run
	$(CLI) load-acl    --credentials $(CREDENTIALS) --rules $(RULES) --force

# GENERAL COMMANDS

usage: build
	$(CLI)

help: build
	$(CLI) help

version: build
	$(CLI) version

# ACL COMMANDS

get-members: build
	$(CLI) --debug get-members --credentials $(CREDENTIALS)
#	$(CLI) --debug get-members --credentials $(CREDENTIALS) --file "$(WORKDIR)/members.tsv"
#	cat "$(WORKDIR)/members.tsv"

get-groups: build
	$(CLI) --debug get-groups --credentials $(CREDENTIALS)
#	$(CLI) --debug get-groups --credentials $(CREDENTIALS) --file "$(WORKDIR)/groups.tsv"
#	cat "$(WORKDIR)/groups.tsv"

get-doors: build
	$(CLI) --debug get-doors
	$(CLI) --debug get-doors --file "$(WORKDIR)/doors.tsv"
	cat "$(WORKDIR)/doors.tsv"

get-acl: build
	$(CLI) --debug get-acl --credentials $(CREDENTIALS) --rules $(RULES) --file "$(WORKDIR)/ACL.tsv"

get-acl-file: build
	$(CLI) get-acl --credentials $(CREDENTIALS) --rules "file://../runtime/wild-apricot/wild-apricot.grl" --file "$(WORKDIR)/ACL.tsv"

get-acl-drive: build
	$(CLI) get-acl --credentials $(CREDENTIALS) --rules "https://drive.google.com/uc?export=download&id=1dwc9HFCbjCf4YB2siexk--coI_xOAtul"

compare-acl: build
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES)
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --report "$(WORKDIR)/ACL.rpt"
	cat "$(WORKDIR)/ACL.rpt"

compare-acl-summary: build
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --summary
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --summary --report "$(WORKDIR)/ACL.rpt"
	cat "$(WORKDIR)/ACL.rpt"

load-acl: build
#	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --dry-run --force --log ../runtime/wild-apricot/ACL.log --report ../runtime/wild-apricot/ACL.report
#	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --dry-run --force --log ../runtime/wild-apricot/ACL.log --report ../runtime/wild-apricot/ACL.report.tsv
	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES)


