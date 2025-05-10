SHELL := bash
ROOT := $(shell pwd)

YS-VERSION := 0.1.96
GIT-PREFIX := ../../.git/.local
GIT-BIN := $(GIT-PREFIX)/bin

YAML-FILES := $(wildcard .*.yaml)

export YSPATH := $(ROOT)/lib

ifneq (,$(force))
  _ := $(shell touch Makefile)
endif


default:

env:
	@echo PATH=$(GIT-BIN):$(PATH)

# Using 'make update devel=1` is much faster than the default Docker container.
ifdef devel
YS := $(GIT-BIN)/ys-$(YS-VERSION)

update: $(YS) $(YAML-FILES)

$(YS):
	@echo "Installing ys as '$(YS)'"
	curl -s https://yamlscript.org/install | \
	  QUIET=1 BIN=1 VERSION=$(YS-VERSION) PREFIX=$(GIT-PREFIX) bash
else
YS := docker run \
        -v $(ROOT):/work \
	-w /work \
	--env YSPATH=lib \
	ingy/ys:$(YS-VERSION) ys -Y

update: $(YAML-FILES)
endif

.%.yaml: %.ys Makefile lib/* bash/*
	$(YS) -Y $< > $@
