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
	go get -u github.com/uhppoted/uhppote-core@main
	go get -u github.com/uhppoted/uhppoted-lib@main
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
	echo "*** WARNING: pending update to staticcheck to Go 1.24.4"
	# env GOOS=darwin  GOARCH=amd64 staticcheck ./...
	# env GOOS=linux   GOARCH=amd64 staticcheck ./...
	# env GOOS=windows GOARCH=amd64 staticcheck ./...

vuln:
	govulncheck ./...

build-all: build test vet lint
	mkdir -p dist/$(DIST)/linux
	mkdir -p dist/$(DIST)/arm
	mkdir -p dist/$(DIST)/arm7
	mkdir -p dist/$(DIST)/arm6
	mkdir -p dist/$(DIST)/darwin-x64
	mkdir -p dist/$(DIST)/darwin-arm64
	mkdir -p dist/$(DIST)/windows
	env GOOS=linux   GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/linux        ./...
	env GOOS=linux   GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/arm          ./...
	env GOOS=linux   GOARCH=arm   GOARM=7 GOWORK=off go build -trimpath -o dist/$(DIST)/arm7         ./...
	env GOOS=linux   GOARCH=arm   GOARM=6 GOWORK=off go build -trimpath -o dist/$(DIST)/arm6         ./...
	env GOOS=darwin  GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin-x64   ./...
	env GOOS=darwin  GOARCH=arm64         GOWORK=off go build -trimpath -o dist/$(DIST)/darwin-arm64 ./...
	env GOOS=windows GOARCH=amd64         GOWORK=off go build -trimpath -o dist/$(DIST)/windows      ./...

release: update-release build-all
	find . -name ".DS_Store" -delete
	tar --directory=dist/$(DIST)/linux        --exclude=".DS_Store" -cvzf dist/$(DIST)-linux-x64.tar.gz    .
	tar --directory=dist/$(DIST)/arm          --exclude=".DS_Store" -cvzf dist/$(DIST)-arm-x64.tar.gz      .
	tar --directory=dist/$(DIST)/arm7         --exclude=".DS_Store" -cvzf dist/$(DIST)-arm7.tar.gz         .
	tar --directory=dist/$(DIST)/arm6         --exclude=".DS_Store" -cvzf dist/$(DIST)-arm6.tar.gz         .
	tar --directory=dist/$(DIST)/darwin-x64   --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin-x64.tar.gz   .
	tar --directory=dist/$(DIST)/darwin-arm64 --exclude=".DS_Store" -cvzf dist/$(DIST)-darwin-arm64.tar.gz .
	cd dist/$(DIST)/windows && zip --recurse-paths ../../$(DIST)-windows-x64.zip . -x ".DS_Store"

publish: release
	echo "Releasing version $(VERSION)"
	gh release create "$(VERSION)" "./dist/$(DIST)-arm-x64.tar.gz"      \
	                               "./dist/$(DIST)-arm7.tar.gz"         \
	                               "./dist/$(DIST)-arm6.tar.gz"         \
	                               "./dist/$(DIST)-darwin-arm64.tar.gz" \
	                               "./dist/$(DIST)-darwin-x64.tar.gz"   \
	                               "./dist/$(DIST)-linux-x64.tar.gz"    \
	                               "./dist/$(DIST)-windows-x64.zip"     \
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
	$(CLI) --config ../runtime/wild-apricot/uhppoted.conf get-members --credentials $(CREDENTIALS) --with-pin
	$(CLI) --config ../runtime/wild-apricot/uhppoted.conf get-members --credentials $(CREDENTIALS) --with-pin --file "$(WORKDIR)/members.tsv"
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
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf get-acl --credentials $(CREDENTIALS) --rules $(RULES_WITH_PIN) --with-pin
	# $(CLI) --debug get-acl --credentials $(CREDENTIALS) --rules $(RULES_WITH_PIN) --with-pin
	# $(CLI) get-acl --credentials $(CREDENTIALS) --rules $(RULES_WITH_PIN) --with-pin --file "$(WORKDIR)/ACL.tsv"
	# cat "$(WORKDIR)/ACL.tsv"

get-acl-file: build
	$(CLI)  --debug --config ../runtime/wild-apricot/uhppoted.conf get-acl --credentials $(CREDENTIALS) --rules "file://../runtime/wild-apricot/wild-apricot.grl" --file "$(WORKDIR)/ACL.tsv"

get-acl-drive: build
	$(CLI)  --debug --config ../runtime/wild-apricot/uhppoted.conf get-acl --credentials $(CREDENTIALS) --rules "https://drive.google.com/uc?export=download&id=1dwc9HFCbjCf4YB2siexk--coI_xOAtul"

compare-acl: build
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf compare-acl --credentials $(CREDENTIALS) --rules $(RULES)
	# $(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --report "$(WORKDIR)/ACL.rpt"
	# cat "$(WORKDIR)/ACL.rpt"

compare-acl-with-pin: build
	$(CLI)  --debug --config ../runtime/wild-apricot/uhppoted.conf compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin
	# $(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin
	# $(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin --report "$(WORKDIR)/ACL.rpt"
	# cat "$(WORKDIR)/ACL.rpt"
	# $(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin --summary
	# $(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin --summary --report "$(WORKDIR)/ACL.rpt"
	# cat "$(WORKDIR)/ACL.rpt"

compare-acl-summary: build
	$(CLI)  --debug --config ../runtime/wild-apricot/uhppoted.conf compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --summary
	# $(CLI) compare-acl --credentials $(CREDENTIALS) --rules $(RULES) --summary --report "$(WORKDIR)/ACL.rpt"
	# cat "$(WORKDIR)/ACL.rpt"

load-acl: build
#	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --dry-run --force --log ../runtime/wild-apricot/ACL.log --report ../runtime/wild-apricot/ACL.report
#	$(CLI) load-acl --credentials $(CREDENTIALS) --rules $(RULES) --dry-run --force --log ../runtime/wild-apricot/ACL.log --report ../runtime/wild-apricot/ACL.report.tsv
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf load-acl --credentials $(CREDENTIALS) --rules $(RULES)

load-acl-with-pin: build
	# $(CLI) --debug load-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin
	$(CLI) --debug --config ../runtime/wild-apricot/uhppoted.conf load-acl --credentials $(CREDENTIALS) --rules $(RULES) --with-pin


