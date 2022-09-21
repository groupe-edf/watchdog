.DEFAULT_GOAL := help
.PHONY: \
	build \
	changelog \
	clean \
	compile \
	docs \
	fmt \
	help \
	lint \
	publish \
	run \
	serve \
	test \
	todo

SHELL := /bin/bash

BLUE   := $(shell tput -Txterm setaf 6)
GREEN  := $(shell tput -Txterm setaf 2)
RED    := $(shell tput -Txterm setaf 1)
RESET  := $(shell tput -Txterm sgr0)
YELLOW := $(shell tput -Txterm setaf 3)

export CGO_ENABLED=1
export COMPOSE_CONVERT_WINDOWS_PATHS=1

GO ?= go
GO_BUILD=$(GO) build
GO_CLEAN=$(GO) clean
GO_GET=GO111MODULE=off $(GO) get -u
GO_INSTALL=$(GO) install
GO_RUN=$(GO) run
GO_TEST=$(GO) test -v
GOOS=$(shell $(GO) env GOOS)

GIT_COMMIT=$(shell git rev-parse HEAD)
GIT_SHA=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --tags --abbrev=0 --exact-match 2>/dev/null)
GIT_DIRTY=$(shell test -n "`git status --porcelain`" && echo "dirty" || echo "clean")

BINARY_NAME=$(notdir $(CURDIR))
PLUGIN_DIRECTORY=pkg/modules
PLUGIN_OUTPUT=plugins
PLUGIN_EXTENSION=so
TARGET=target

BINARY_VERSION ?= ${GIT_TAG}
ifeq ($(BINARY_VERSION),)
	BINARY_VERSION := $(shell cat ./VERSION)
endif
BUILD_DATE := $(shell date +'%d.%m.%Y')

# Use linker flags to provide version/build settings to the target
# List variables with go tool nm | grep ...
LDFLAGS :=-w -s
ifneq ($(BINARY_VERSION),)
	LDFLAGS += -X github.com/groupe-edf/watchdog/internal/version.Version=${BINARY_VERSION}
endif
LDFLAGS += -X github.com/groupe-edf/watchdog/internal/version.BuildDate=$(BUILD_DATE)
LDFLAGS += -X github.com/groupe-edf/watchdog/internal/version.Commit=$(GIT_COMMIT)
LDFLAGS += -X github.com/groupe-edf/watchdog/internal/version.Sha=$(GIT_SHA)
LDFLAGS_TEST := ${LDFLAGS} -X github.com/groupe-edf/watchdog/pkg/config.LogPath=watchdog.log

ifeq ($(GOOS), windows)
	export CGO_ENABLED=1
	BINARY_OUTPUT := $(TARGET)/bin/${BINARY_NAME}.exe
else
	GO_TEST := $(GO_TEST) -race
	BINARY_OUTPUT := $(TARGET)/bin/${BINARY_NAME}
endif

$(TARGET):
	@mkdir -p $@

all: lint test build

bootstrap: ## Install all development and ci tools
	$(GO_INSTALL) github.com/client9/misspell/cmd/misspell@latest
	$(GO_INSTALL) github.com/git-chglog/git-chglog/cmd/git-chglog@latest
	$(GO_INSTALL) github.com/go-swagger/go-swagger/cmd/swagger@latest
	$(GO_INSTALL) github.com/golangci/golangci-lint/cmd/golangci-lint@latest
	$(GO_INSTALL) github.com/securego/gosec/v2/cmd/gosec@latest
	$(GO_INSTALL) golang.org/x/tools/cmd/cover@latest
	$(GO_INSTALL) golang.org/x/tools/cmd/godoc@latest
	$(GO_INSTALL) golang.org/x/tools/cmd/goimports@latest
	$(GO_INSTALL) golang.org/x/lint/golint@latest
	$(GO_INSTALL) golang.org/x/vuln/cmd/govulncheck@latest

build: ## Build watchdog CLI
	$(GO_BUILD) -o $(BINARY_NAME) -v -o $(BINARY_OUTPUT) -ldflags="$(LDFLAGS)" ./cmd/watchdog-cli
	@echo "${GREEN}> Build completed successfully in $(BINARY_OUTPUT)${RESET}"

build-plugins: $(PLUGIN_DIRECTORY)/*.go
	$(GO_BUILD) -buildmode=plugin -o $(PLUGIN_OUTPUT)/$(basename $(<F)).$(PLUGIN_EXTENSION) ./$^

build-test:
	$(GO_BUILD) -o $(BINARY_NAME) -v -o $(BINARY_OUTPUT) -ldflags="$(LDFLAGS_TEST)" ./cmd/watchdog-cli

changelog:
	# https://keepachangelog.com/en/1.0.0/
	@git-chglog --config .ci/git-chglog.yml --output=CHANGELOG.md $(VERSION)

clean:
	$(GO) clean -testcache
	@rm -rf $(TARGET) bin sonar-scanner-* .scannerwork *.test

check: fmt vet lint

# Bootstrap documentation
#		git submodule update --init --recursive
#		npm install -D --save autoprefixer
#		npm install -D --save postcss-cli
docs: ## Generate watchdog documentation
	@hugo --verbose --source docs --destination ./public

docker-build:
	docker build \
		--file build/Dockerfile \
		--tag ${DOCKER_REGISTRY}/watchdog/watchdog-server \
		.

docs-serve:
	@hugo server --watch --source docs

fmt:
	@$(GO) fmt ./...

generate:
	go generate

help:
	@echo ""
	@echo "    ${BLACK}:: ${RED}Self-documenting Makefile${RESET} ${BLACK}::${RESET}"
	@echo ""
	@grep -E '^[a-zA-Z_0-9%-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "${BLUE}%-30s${RESET} %s\n", $$1, $$2}'

info:
	@echo "Version: ${BINARY_VERSION}"
	@echo "Git Tag: ${GIT_TAG}"
	@echo "Git Commit: ${GIT_COMMIT}"
	@echo "Git Tree State: ${GIT_DIRTY}"

install: ## Install watchdog in local
	$(GO) install -ldflags="$(LDFLAGS)"

lint:
	@golint -set_exit_status ./...
	@golangci-lint run ./...
	@misspell -error docs/content/**/* internal pkg test

release:
	@goreleaser release --skip-publish --rm-dist

release-snapshot:
	@goreleaser release --skip-publish --snapshot --rm-dist

URI=
run: ## Run watchdog locally to analyze repostiory `make run URI="https://github.com/groupe-edf/watchdog"`
	@$(GO_RUN) -ldflags="$(LDFLAGS)" ./cmd/watchdog-cli analyze \
		--config config/config.yml \
		--logs-path /var/log/watchdog/watchdog.log \
		--policies-file .watchdog.yml \
		--uri "$(URI)"
	@echo "${GREEN}> Repository successfully analyzed${RESET}"

serve: ## Serve watchdog
	export BUILD=production
	$(GO_RUN) -ldflags="$(LDFLAGS)" ./cmd/watchdog-server serve \
		--config config/config.yml

swagger-generate:
	export SWAGGER_GENERATE_EXTENSION=false
	@swagger generate spec --output swagger.json --scan-models

swagger: swagger-generate  ## Serve a spec and swagger or redoc documentation ui
	@swagger serve swagger.json --no-open --port=4444 --flavor=swagger


test: test-unit test-integration test-security clean

TAGS=integration
RUN=.
test-integration:
	$(GO_TEST) --tags=$(TAGS) ./test/integration/... -run $(RUN)

test-security:
	@gosec -exclude=G101,G104,G203,G204,G302,G306,G307 -fmt=json ./...
	@govulncheck ./...

COVERAGE_PROFILE=$(TARGET)/coverage.txt
RUN=.
test-unit:
	@mkdir -p $(TARGET)
	$(GO_TEST) -coverprofile=$(COVERAGE_PROFILE) -covermode=atomic ./... -run $(RUN)

tidy:
	@$(GO) mod tidy

# Show to-do items per file.
todo:
	@grep \
		--exclude-dir=vendor \
		--exclude=Makefile \
		--text \
		--color \
		-nRo -E ' (TODO|FIXME):.*|SkipNow' .

uninstall:
	@rm -f $$(which ${BINARY_NAME})

vet:
	@$(GO) vet ./...
