.DEFAULT_GOAL: default

VERSION ?= $(shell git describe --abbrev=0 --tags 2>/dev/null || echo 0.0.0)
COMMIT ?= $(shell git rev-parse --short HEAD 2>/dev/null || echo HEAD)
BUILD_DATE ?= $(shell date -u +"%Y-%m-%dT%H:%M:%SZ")

GOOS ?= $(shell go env GOOS)
GOARCH ?= $(shell go env GOARCH)
GOPATH ?= $(shell pwd)/_workspace
GOBASE := $(firstword $(subst :, ,$(GOPATH)))
GOBIN := $(GOBASE)/bin
GO ?= go
GOFMT ?= gofmt
GOLINT := $(GOBIN)/golangci-lint
TAR ?= tar

export CGO_ENABLED = 0
export GO111MODULE = on
export GOARCH
export GOOS
export GOBIN

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
	-$(RM) *.gz
	-$(RM) -r ./out
	@$(GO) clean -x

.PHONY: lint
lint: $(GOLINT) $(GO_SOURCES)
	$(GOLINT) run ./...

.PHONY: format
format: $(GO_SOURCES)
	@$(GOFMT) -s -w $^

.PHONY: test
test: $(GO_SOURCES)
	@$(GO) test $(GO_PACKAGES)

.PHONY: build
build: $(addprefix out/,$(PROGRAMS))

.PHONY: dist
dist: $(notdir $(GO_MODULE))-$(VERSION)-$(GOARCH).tar.gz

$(GOBIN)/%:
	# go install -v -tags tools ./...
	- grep '_ "' tools/tools.go | \
		awk '{ print $$2 }' | \
		xargs -n1 $(GO) install -v

out/%: $(GO_SOURCES)
	$(GO) build \
		-o $@ \
		-ldflags $(GO_LDFLAGS) \
		$(GO_MODULE)/cmd/$(notdir $(basename $@))

%.tar.gz: build
	$(TAR) -cvzf $@ \
		-C out \
		$(notdir $(wildcard ./out/*))
