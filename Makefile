DIST       ?= development
CLI         = ./bin/uhppoted-app-wild-apricot
WORKDIR     = ../runtime/wild-apricot
CREDENTIALS = $(WORKDIR)/.credentials.json
RULES       = $(WORKDIR)/debug.grl
# RULES       = $(WORKDIR)/wild-apricot.grl
RULES_WITH_PIN = $(WORKDIR)/wild-apricot-with-pin.grl

DATETIME  = $(shell date "+%Y-%m-%d %H:%M:%S")
DEBUG    ?= --debug

.PHONY: clean
.PHONY: bump
.PHONY: bump-release

all: test      \
	 benchmark \
     coverage

clean:
	go clean
	rm -rf bin

update:
	go get -u github.com/uhppoted/uhppote-core@master
	go get -u github.com/uhppoted/uhppoted-lib@master
	go get -u github.com/hyperjumptech/grule-rule-engine
	go get -u golang.org/x/sys
	go mod tidy

update-release:
	go get -u github.com/uhppoted/uhppote-core
	go get -u github.com/uhppoted/uhppoted-lib
	go mod tidy

update-all:
	go get -u github.com/uhppoted/uhppote-core
	go get -u github.com/uhppoted/uhppoted-lib
	go get -u github.com/hyperjumptech/grule-rule-engine
	go get -u golang.org/x/sys
	go mod tidy

format: 
	go fmt ./...

build: format
	mkdir -p bin
	go build -trimpath -o bin ./...

test: build
	go test ./...

benchmark: build
	go test -bench ./...

coverage: build
	go test -cover ./...

vet: build
	go vet ./...

lint: build
	env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	env GOOS=linux   GOARCH=amd64 staticcheck ./...
	env GOOS=windows GOARCH=amd64 staticcheck ./...

vuln:
	govulncheck ./...

build-all: test vet lint
	mkdir -p dist/$(DIST)/windows
	mkdir -p dist/$(DIST)/darwin
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm
	mkdir -p dist/$(DIST)/arm7
	env GOOS=linux   GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/linux   ./...
	env GOOS=linux   GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/arm     ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7    ./...
	env GOOS=darwin  GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin  ./...
	env GOOS=windows GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/windows ./...

release: update-release build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist --exclude=".DS_Store" -cvzf dist/$(DIST).tar.gz $(DIST)
	cd dist;  zip --recurse-paths $(DIST).zip $(DIST)

publish: release
	echo "Releasing version $(VERSION)"
	gh release create "$(VERSION)" \
	   "./dist/uhppoted-app-wild-apricot_$(VERSION).tar.gz" \
	   "./dist/uhppoted-app-wild-apricot_$(VERSION).zip" \
	   --draft --prerelease --title "$(VERSION)-beta" --notes-file release-notes.md

debug: build
	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --dry-run  --lockfile ".lock.me"
	# $(CLI) load-acl    --credentials $(CREDENTIALS) --rules $(RULES) --force --dry-run
	# $(CLI) load-acl    --credentials $(CREDENTIALS) --rules $(RULES) --force

godoc:
	godoc -http=:80	-index_interval=60s

# GENERAL COMMANDS

usage: build
	$(CLI)

help: build
	$(CLI) help

version: build
	$(CLI) version

# ACL COMMANDS

get-members: build
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf get-members --credentials $(CREDENTIALS)
#	$(CLI) --debug get-members --credentials $(CREDENTIALS) --file "$(WORKDIR)/members.tsv"
#	cat "$(WORKDIR)/members.tsv"

get-members-with-pin: build
	$(CLI) get-members --credentials $(CREDENTIALS) --with-pin
	$(CLI) get-members --credentials $(CREDENTIALS) --with-pin --file "$(WORKDIR)/members.tsv"
	cat "$(WORKDIR)/members.tsv"

get-groups: build
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf  get-groups --credentials $(CREDENTIALS)
#	$(CLI) --debug get-groups --credentials $(CREDENTIALS) --file "$(WORKDIR)/groups.tsv"
#	cat "$(WORKDIR)/groups.tsv"

get-doors: build
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf get-doors
	# $(CLI) --debug get-doors --file "$(WORKDIR)/doors.tsv"
	# cat "$(WORKDIR)/doors.tsv"

get-acl: build
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf get-acl --credentials $(CREDENTIALS) --rules $(RULES)

get-acl-with-pin: build
	$(CLI) get-acl --credentials $(CREDENTIALS) --rules $(RULES_WITH_PIN) --with-pin
	# $(CLI) --debug get-acl --credentials $(CREDENTIALS) --rules $(RULES_WITH_PIN) --with-pin
	# $(CLI) get-acl --credentials $(CREDENTIALS) --rules $(RULES_WITH_PIN) --with-pin --file "$(WORKDIR)/ACL.tsv"
	# cat "$(WORKDIR)/ACL.tsv"

get-acl-file: build
	$(CLI) get-acl --credentials $(CREDENTIALS) --rules "file://../runtime/wild-apricot/wild-apricot.grl" --file "$(WORKDIR)/ACL.tsv"

get-acl-drive: build
	$(CLI) get-acl --credentials $(CREDENTIALS) --rules "https://drive.google.com/uc?export=download&id=1dwc9HFCbjCf4YB2siexk--coI_xOAtul"

compare-acl: build
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf compare-acl --credentials $(CREDENTIALS) --rules $(RULES)
	# $(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --report "$(WORKDIR)/ACL.rpt"
	# cat "$(WORKDIR)/ACL.rpt"

compare-acl-with-pin: build
#	$(CLI) --debug compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin --report "$(WORKDIR)/ACL.rpt"
	cat "$(WORKDIR)/ACL.rpt"
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin --summary
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin --summary --report "$(WORKDIR)/ACL.rpt"
	cat "$(WORKDIR)/ACL.rpt"

compare-acl-summary: build
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --summary
	$(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --summary --report "$(WORKDIR)/ACL.rpt"
	cat "$(WORKDIR)/ACL.rpt"

load-acl: build
#	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --dry-run --force --log ../runtime/wild-apricot/ACL.log --report ../runtime/wild-apricot/ACL.report
#	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --dry-run --force --log ../runtime/wild-apricot/ACL.log --report ../runtime/wild-apricot/ACL.report.tsv
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf load-acl --credentials $(CREDENTIALS) --rules $(RULES)

load-acl-with-pin: build
	# $(CLI) --debug load-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin
	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin


