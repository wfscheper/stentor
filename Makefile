TARGET   = bin/stentor
SOURCES := $(shell find . -name '*.go' -a -not -name '*_test.go')

# commands
GOTAGGER = bin/gotagger
LINTER   = bin/golangci-lint
RELEASER = bin/goreleaser
TESTER   = bin/gotestsum

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

.PHONY: all
all: lint build test

.PHONY: build
build: $(TARGET)

.PHONY: changelog
changelog: $(TARGET) | $(GOTAGGER) ; $(info $(M) generating changelog)
	$Q $(TARGET) $(STENTORFLAGS) $(shell bin/gotagger) $(shell git tag --list --sort=-version:refname | head -n1)

.PHONY: clean
clean:
	$Q $(RM) -r bin/ dist/ reports/

.PHONY: dist
dist: | $(RELEASER) ; $(info $(M) building dist…)
	$(RELEASER) release $(RELEASEFLAGS)

.PHONY: fmt format ; $(info $(M) formatting…)
fmt format: LINTERFLAGS += --fix
fmt format: lint | $(LINTER)

.PHONY: lint
lint: | $(LINTER) ; $(info $(M) linting…)
	$Q $(LINTER) run $(LINTERFLAGS)

.PHONY: test tests
test tests: | $(TESTER) ; $(info $(M) running tests…)
	$Q $(TESTER) $(TESTERFLAGS) -- $(TESTFLAGS) ./...

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

$(GOTAGGER): tools/go.mod tools/go.sum ; $(info $(M) building $(GOTAGGER)…)
	cd tools && GOBIN=$(CURDIR)/bin go install github.com/sassoftware/gotagger/cmd/gotagger

$(LINTER): tools/go.mod tools/go.sum ; $(info $(M) building $(LINTER)…)
	cd tools && GOBIN=$(CURDIR)/bin go install github.com/golangci/golangci-lint/cmd/golangci-lint

$(RELEASER): tools/go.mod tools/go.sum ; $(info $(M) building $(RELEASER)…)
	cd tools && GOBIN=$(CURDIR)/bin go install github.com/goreleaser/goreleaser

$(REPORTDIR):
	@mkdir -p $@

$(TARGET): $(SOURCES) go.mod go.sum | $(GOTAGGER); $(info $(M) building $(TARGET)…)
	go build -o $@ -mod=readonly -ldflags "-X main.buildDate=$(BUILDDATE) -X main.commit=$(COMMIT) -X main.version=$$(bin/gotagger)" ./cmd/stentor/

$(TESTER): tools/go.mod tools/go.sum ; $(info $(M) building $(TESTER)…)
	cd tools && GOBIN=$(CURDIR)/bin go install gotest.tools/gotestsum
