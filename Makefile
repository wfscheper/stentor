include .bingo/Variables.mk

TARGET   = bin/stentor

# commands
TESTER   = tools/gotestsum

# variables
BUILDDATE   := $(shell date +%Y-%m-%d)
COMMIT      := $(shell git rev-parse HEAD)
RELEASEFLAGS = $(if $(filter false,$(DRYRUN)),,--snapshot) --rm-dist
REPORTDIR    = reports
STENTORFLAGS = $(if $(filter false,$(DRYRUN)),-release)
TESTFLAGS    = -cover -covermode=atomic

# output controls
override Q = $(if $(filter 1,$(V)),,@)
override M = ▶

.DEFAULT_GOAL := all

.PHONY: all
all: lint build test

.PHONY: build
build: $(TARGET)

.PHONY: changelog
changelog: $(TARGET) | $(GOTAGGER) ; $(info $(M) generating changelog)
	$Q $(TARGET) $(STENTORFLAGS) $(shell $(GOTAGGER)) $(shell git tag --list --sort=-version:refname | head -n1)

.PHONY: clean
clean:
	$Q $(RM) -r bin/ dist/ reports/

.PHONY: dist
dist: | $(GORELEASER) ; $(info $(M) building dist…)
	$(GORELEASER) release $(RELEASEFLAGS)

.PHONY: fmt format ; $(info $(M) formatting…)
fmt format: LINTERFLAGS += --fix
fmt format: lint | $(GOLANGCI_LINT)

.PHONY: lint
lint: | $(GOLANGCI_LINT) ; $(info $(M) linting…)
	$Q $(GOLANGCI_LINT) run $(LINTERFLAGS)

.PHONY: test tests
test tests: | $(GOTESTSUM) ; $(info $(M) running tests…)
	$Q $(GOTESTSUM) $(TESTERFLAGS) -- $(TESTFLAGS) ./...

.PHONY: test-report
test-report: TESTERFLAGS += --junitfile reports/junit.xml
test-report: TESTFLAGS += -coverprofile=reports/cover.out
test-report: $(REPORTDIR)
test-report: test

.PHONY: test-watch
test-watch: TESTERFLAGS += --watch
test-watch: test

.PHONY: show
show:
	@echo $(VALUE)=$($(VALUE))

.PHONY: version
version: | $(GOTAGGER)
	@$(GOTAGGER)

.PHONY: FORCE
$(TARGET): FORCE | $(GOTAGGER); $(info $(M) building $(TARGET)…)
	$Q go build -o $@ -mod=readonly -ldflags "-X main.buildDate=$(BUILDDATE) -X main.commit=$(COMMIT) -X main.version=$$($(GOTAGGER))" ./cmd/stentor/

$(REPORTDIR):
	@mkdir -p $@
