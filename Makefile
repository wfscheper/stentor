TARGET   = bin/stentor
SOURCES := $(shell find . -name '*.go' -a -not -name '*_test.go')

# commands
GOTAGGER = tools/gotagger
LINTER   = tools/golangci-lint
RELEASER = tools/goreleaser
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

.PHONY: FORCE
$(TARGET): FORCE | $(GOTAGGER); $(info $(M) building $(TARGET)…)
	$Q go build -o $@ -mod=readonly -ldflags "-X main.buildDate=$(BUILDDATE) -X main.commit=$(COMMIT) -X main.version=$$($(GOTAGGER))" ./cmd/stentor/

$(REPORTDIR):
	@mkdir -p $@

define build_tool
$(1): tools/go.mod tools/go.sum ; $$(info $$(M) building $(1)…)
	$Q cd tools && go build -mod=readonly $(2)
endef

define update_tool
.PHONY: update-$(notdir $(1))
update-$(notdir $(1)): ; $$(info $$(M) updating $(notdir $(1))…)
	$Q cd tools && go get $(2)
endef

$(eval $(call build_tool,$(GOTAGGER),github.com/sassoftware/gotagger/cmd/gotagger))
$(eval $(call build_tool,$(LINTER),github.com/golangci/golangci-lint/cmd/golangci-lint))
$(eval $(call build_tool,$(RELEASER),github.com/goreleaser/goreleaser))
$(eval $(call build_tool,$(TESTER),gotest.tools/gotestsum))
$(eval $(call update_tool,$(GOTAGGER),github.com/sassoftware/gotagger/cmd/gotagger))
$(eval $(call update_tool,$(LINTER),github.com/golangci/golangci-lint/cmd/golangci-lint))
$(eval $(call update_tool,$(RELEASER),github.com/goreleaser/goreleaser))
$(eval $(call update_tool,$(TESTER),gotest.tools/gotestsum))
