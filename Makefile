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
ZIP ?= zip
INSTALL ?= install
INSTALL_PROGRAM = $(INSTALL)
INSTALL_DATA = $(INSTALL) -m 644
prefix ?= /usr/local
datarootdir = $(prefix)/share
datadir = $(datarootdir)
exec_prefix = $(prefix)
bindir = $(exec_prefix)/bin
libexecdir = $(exec_prefix)/libexec
infodir = $(datarootdir)/info

export CGO_ENABLED = 0
export GO111MODULE = on
export GOARCH
export GOOS
export GOBIN

GO_MODULE := $(shell $(GO) list -m)
GO_PACKAGES := $(shell $(GO) list ./... | grep -vE '/(tools|test|vendor)')
GO_SOURCES := $(shell find . -type f -name '*.go' | grep -vE '/(tools|test|vendor)/')
GO_CMDS := $(notdir $(wildcard ./cmd/*))

PROJECT_NAME ?= $(notdir $(GO_MODULE))
BUILD_DIR ?= out

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
	-$(RM) *.gz *.xz *.tar *.zip
	-$(RM) -r $(BUILD_DIR)
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
build: $(addprefix $(BUILD_DIR)/,$(PROGRAMS))

.PHONY: binary-dist
binary-dist: $(PROJECT_NAME)-$(VERSION)-$(GOARCH).tar.gz

.PHONY: dist
dist: tar

.PHONY: tar
tar: $(PROJECT_NAME)-$(VERSION).tar

.PHONY: install
install: build
	$(INSTALL) -d $(DESTDIR)$(bindir)/
	$(INSTALL_PROGRAM) $(wildcard $(BUILD_DIR)/*) $(DESTDIR)$(bindir)/

$(GOBIN)/%:
	# go install -v -tags tools ./...
	- grep '_ "' tools/tools.go | \
		awk '{ print $$2 }' | \
		xargs -n1 $(GO) install -v

$(BUILD_DIR)/%: $(GO_SOURCES)
	$(GO) build \
		-o $@ \
		-ldflags $(GO_LDFLAGS) \
		$(GO_MODULE)/cmd/$(notdir $(basename $@))

%.tar: $(GO_SOURCES)
	$(TAR) -cf $@ \
		--exclude-vcs \
		--exclude-vcs-ignores \
		--exclude .github \
		$(shell git ls-files --exclude-standard)

%.zip: build
	$(ZIP) -jr $@ $(BUILD_DIR)

%.tar.gz: build
	$(TAR) -czf $@ \
		-C $(BUILD_DIR) \
		$(notdir $(wildcard $(BUILD_DIR)/*))
