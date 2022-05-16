#-------------------------------------------------------------------------------
# Global variables.

GO=$(shell which go)

#-------------------------------------------------------------------------------
# Running `make` will show the list of subcommands that will run.

all: help

.PHONY: help
## help: prints this help message
help:
	@echo "Usage:"
	@echo ""
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

#-------------------------------------------------------------------------------
# Dependencies


#-------------------------------------------------------------------------------
# Compile

.PHONY: build-prep
## build-prep: [build] updates go.mod and downloads dependencies
build-prep:
	mkdir -p ./bin
	$(GO) mod tidy -go=1.17 -v
	$(GO) mod download -x
	$(GO) get -v ./...

.PHONY: build-release-prep
## build-release-prep: [build] post-development, ready to release steps
build-release-prep:
	$(GO) mod download

.PHONY: build
## build: [build] compiles the source code into a native binary
build: build-prep
	$(GO) build -ldflags="-s -w  -X main.commit=$$(git rev-parse HEAD) -X main.date=$$(date -I) -X main.version=$$(cat ./VERSION | tr -d '\n')" -o ./bin/ssm-shell *.go


.PHONY: install
## install: [build] Installs the command to ~/.bin/, which should be on your PATH.
install:
	mkdir -p ~/bin
	cp -fv bin/ssm-shell ~/bin/ssm-shell

#-------------------------------------------------------------------------------
# Clean

.PHONY: clean-go
## clean-go: [clean] clean Go's module cache
clean-go:
	$(GO) clean -i -r -x -testcache -modcache -cache

.PHONY: clean
## clean: [clean] runs ALL non-Docker cleaning tasks
clean: clean-go

#-------------------------------------------------------------------------------
# Linting

.PHONY: golint
## golint: [lint] runs `golangci-lint` (static analysis, formatting) against all Golang (*.go) tests with a standardized set of rules
golint:
	@ echo " "
	@ echo "=====> Running gofmt and golangci-lint..."
	gofmt -s -w *.go
	golangci-lint run --fix *.go

.PHONY: goupdate
## goupdate: [lint] runs `go-mod-outdated` to check for out-of-date packages
goupdate:
	@ echo " "
	@ echo "=====> Running go-mod-outdated..."
	$(GO) list -u -m -json all | go-mod-outdated -update -direct -style markdown

.PHONY: goconsistent
## goconsistent: [lint] runs `go-consistent` to verify that implementation patterns are consistent throughout the project
goconsistent:
	@ echo " "
	@ echo "=====> Running go-consistent..."
	go-consistent -v ./...

.PHONY: goimportorder
## goimportorder: [lint] runs `go-consistent` to verify that implementation patterns are consistent throughout the project
goimportorder:
	@ echo " "
	@ echo "=====> Running impi..."
	impi --local github.com/northwood-labs/ssm-shell --ignore-generated=true --scheme=stdLocalThirdParty ./...

.PHONY: goconst
## goconst: [lint] runs `goconst` to identify values that are re-used and could be constants
goconst:
	@ echo " "
	@ echo "=====> Running goconst..."
	goconst -match-constant -numbers ./...

.PHONY: markdownlint
## markdownlint: [lint] runs `markdownlint` (formatting, spelling) against all Markdown (*.md) documents with a standardized set of rules
markdownlint:
	@ echo " "
	@ echo "=====> Running Markdownlint..."
	npx markdownlint-cli --fix '*.md' --ignore 'node_modules'

.PHONY: lint
## lint: [lint] runs ALL linting/validation tasks
lint: markdownlint golint goupdate goconsistent

#-------------------------------------------------------------------------------
# Documentation and Schema

.PHONY: usage
## usage: [docs] generates the "Command Usage" Markdown documentation from the contents of `--help` options
usage:
	./dist/ssm-shell usage > markdown/usage.md

.PHONY: docs
## docs: [docs] generates the documentation for the JSON Schema
docs: usage
	mkdir -p docs/main/
	mkdocs build --clean --theme=material --site-dir docs/main/

.PHONY: deploy-docs
## deploy-docs: [deploy] Perform a production-mode build of the static artifacts, and push them up to GHE Pages.
deploy-docs:
	rm -Rf /tmp/gh-pages
	git clone git@github.com:northwood-labs/ssm-shell.git --branch gh-pages --single-branch /tmp/gh-pages
	rm -Rf /tmp/gh-pages/*
	cp -Rf ./docs/* /tmp/gh-pages/
	touch /tmp/gh-pages/.nojekyll
	find /tmp/gh-pages -type d | xargs chmod -f 0755
	find /tmp/gh-pages -type f | xargs chmod -f 0644
	cd /tmp/gh-pages/ && \
		git add . && \
		git commit -a -m "Automated commit on $$(date)" && \
		git push origin gh-pages

.PHONY: flatten-docs
## flatten: [flatten-docs] (Optional) Flattens the git history so that git clones can be faster.
flatten-docs:
	rm -Rf /tmp/gh-pages
	git clone git@github.com:northwood-labs/ssm-shell.git --branch gh-pages --single-branch /tmp/gh-pages
	cd /tmp/gh-pages && \
		git checkout --orphan flatten && \
		git add --all . && \
		git commit -a -m "Flattening commit on $$(date)" && \
		git branch -D gh-pages && \
		git branch -m gh-pages && \
		git push -f origin gh-pages

#-------------------------------------------------------------------------------
# Git Tasks

.PHONY: tag
## tag: [release] tags (and GPG-signs) the release
tag:
	@ if [ $$(git status -s -uall | wc -l) != 1 ]; then echo 'ERROR: Git workspace must be clean.'; exit 1; fi;

	@echo "This release will be tagged as: v$$(cat ./VERSION)"
	@echo "This version should match your release. If it doesn't, re-run 'make version'."
	@echo "---------------------------------------------------------------------"
	@read -p "Press any key to continue, or press Control+C to cancel. " x;

	@echo " "
	@chag update v$$(cat ./VERSION)
	@echo " "

	@echo "These are the contents of the CHANGELOG for this release. Are these correct?"
	@echo "---------------------------------------------------------------------"
	@chag contents
	@echo "---------------------------------------------------------------------"
	@echo "Are these release notes correct? If not, cancel and update CHANGELOG.md."
	@read -p "Press any key to continue, or press Control+C to cancel. " x;

	@echo " "

	git add .
	git commit -a -m "Preparing the v$$(cat ./VERSION) release."
	chag tag --sign

.PHONY: version
## version: [release] sets the version for the next release; pre-req for a release tag
version:
	@echo "Current version: $$(cat ./VERSION)"
	@read -p "Enter new version number: " nv; \
	printf "$$nv" > ./VERSION

.PHONY: release
## release: [release] compiles the source code into binaries for all supported platforms and prepares release artifacts
release:
	cd ./src/ && goreleaser release --rm-dist --skip-publish
	mv -vf dist/ssm-shell.rb Formula/ssm-shell.rb
	sha256sum ./dist/ssm-shell_darwin_amd64.zip
