.DEFAULT_GOAL: default

VERSION ?= 0.0.0
COMMIT ?= HEAD
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GO ?= go
GOFMT ?= gofmt
GOLINT ?= golangci-lint

export CGO_ENABLED = 0
export GO111MODULE = on
export GOARCH
export GOOS

GO_MODULE := $(shell $(GO) list -m)
GO_PACKAGES := $(shell $(GO) list ./... | grep -vE '/(tools|test|vendor)')
GO_SOURCES := $(shell find . -type f -name '*.go' | grep -vE '/(tools|test|vendor)/')
GO_CMDS := $(notdir $(wildcard ./cmd/*))

GO_LDFLAGS := '-extldflags "-static"
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.version=$(VERSION)
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.commit=$(COMMIT)
GO_LDFLAGS += -X $(GO_MODULE)/pkg/version.date=$(BUILD_DATE)
GO_LDFLAGS += -w -s # Drop debugging symbols.
GO_LDFLAGS += '

ifeq ($(GOOS),windows)
	PROGRAMS := $(addsuffix .exe,$(GO_CMDS))
else
	PROGRAMS := $(GO_CMDS)
endif

.PHONY: default
default: all

.PHONY: all
all: lint test build

.PHONY: clean
clean:
	-$(RM) -r ./out
	@$(GO) clean -x

.PHONY: lint
lint: $(GO_SOURCES)
	$(GOLINT) run ./...

.PHONY: format
format: $(GO_SOURCES)
	@$(GOFMT) -s -w $<

.PHONY: test
test: $(GO_SOURCES)
	@$(GO) test $(GO_PACKAGES)

.PHONY: build
build: $(addprefix out/,$(PROGRAMS))

out/%: $(GO_SOURCES)
	$(GO) build \
		-o $@ \
		-ldflags $(GO_LDFLAGS) \
		$(GO_MODULE)/cmd/$(notdir $(basename $@))
