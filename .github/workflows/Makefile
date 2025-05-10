SHELL := bash
ROOT := $(shell pwd)

YS-VERSION := 0.1.96
GIT-PREFIX := ../../.git/.local
GIT-BIN := $(GIT-PREFIX)/bin
YS := $(GIT-BIN)/ys-$(YS-VERSION)

YAML-FILES := $(wildcard .*.yaml)

export YSPATH := $(ROOT)/lib


default:
	@echo $(YAML-FILES)

update: $(YS) $(YAML-FILES)

.%.yaml: %.ys
	$(YS) -Y $< > $@

$(YS):
	@echo "Installing ys as '$(YS)'"
	curl -s https://yamlscript.org/install | \
	  QUIET=1 BIN=1 VERSION=$(YS-VERSION) PREFIX=$(GIT-PREFIX) bash
