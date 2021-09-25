MAGE = tools/bin/mage

.PHONY: all
all: | $(MAGE)
	$(MAGE)

.PHONY: help
help: | $(MAGE)
	$(MAGE) -l

$(MAGE): mage/go.mod mage/go.sum mage/mage.go mage/magefile.go
	cd mage/ && go run mage.go -compile ../$@

mage/go.mod mage/go.sum mage/mage.go mage/magefile.go:
	@:

%: | $(MAGE)
	@$(MAGE) $*
