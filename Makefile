ifneq (GNU,$(firstword $(shell $(MAKE) --version)))
  $(error Error: 'make' must be 'GNU make')
endif

SHELL := bash

export MAKE ?= make
export PYTEST ?= pytest-3

library: _prep
	$(MAKE) -C src
	cp src/*.so dist/

release: clean _prep
	$(MAKE) -C src -f Makefile.packages sdist
	./misc/docker-multi-arch-package-build.sh
	@echo "Watch the build process: tail -f docker_build_output/*/*/build.log"
	@echo "When the build is finished, run: $(MAKE) release-upload"

release-upload:
	export WORKDIR=`pwd` && cd src/ && \
		: $(MAKE) -f Makefile.packages prepare-artifacts-for-upload \
		pypi-upload
	@echo "Please upload docker_build_output/upload/build-artifacts/*:"
	ls -1 docker_build_output/upload/build-artifacts/

packages deb: _prep
	$(MAKE) -C src -f Makefile.packages $@

deb-and-wheel-in-podman:
	podman run --rm -it \
	  -v `pwd`:/highctidh \
	  --workdir /highctidh \
	  debian:bookworm \
	    bash -c 'apt update && \
	      ./misc/install-debian-deps.sh && \
	      $(MAKE) wheel && \
	      CC=clang $(MAKE) deb'

wasm: _prep
	CC=emcc $(MAKE) -C src highctidh.wasm

wheel: _prep
	$(MAKE) -C src -f Makefile.packages wheel

pytest:
	cd src; $(PYTEST) -v -n auto --doctest-modules

install test examples:
	$(MAKE) -C src $@

update-golang-modules test-python:
	$(MAKE) -C src -f Makefile.packages $@

test-quick:
	$(MAKE) -C src testrandom test512
	`pwd`/src/$@.sh

test-go:
	cd src; go test -v ./...
	cd src/ctidh511; go test -v ./...
	cd src/ctidh512; go test -v ./...
	cd src/ctidh1024; go test -v ./...
	cd src/ctidh2048; go test -v ./...

examples-run:
	cd src; time ./example-ctidh511
	cd src; time ./example-ctidh512
	cd src; time ./example-ctidh1024
	cd src; time ./example-ctidh2048

alpine-multi-arch-deps:
	@echo 'Attempting to install dependencies'
	uname -a
	./misc/install-alpine-deps.sh
	@echo 'Attempting to build, install, and test'

ALPINE_DEPS := \
  alpine-multi-arch-deps \
  library \
  install \
  examples \
  examples-run \
  test-python \

alpine-multi-arch: $(ALPINE_DEPS)
	@echo 'Hopefully the above was successful'
	ls -alh dist/*.so && sha256sum dist/*.so

_prep:
	-mkdir src/dist
	-ln -s src/build build
	-ln -s src/dist dist

clean:
	-$(RM) -r build
	-$(RM) -r dist
	-$(RM) -r deb_dist
	-$(RM) -r docker_build_output
	-$(RM) -r .pytest_cache
	-$(MAKE) -C src clean
	-$(MAKE) -C src -f Makefile.packages clean
