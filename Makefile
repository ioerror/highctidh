export MAKE ?= make
export PYTEST ?=pytest
library: _prep
	cd src; $(MAKE);
	cp src/*.so dist/;

release: clean _prep
	cd src; $(MAKE) -f Makefile.packages sdist;
	./misc/docker-multi-arch-package-build.sh;
	echo "Watch the build process: tail -f docker_build_output/*/*/build.log";
	echo "When the build is finished, run: $(MAKE) release-upload";

release-upload:
	export WORKDIR=`pwd` && cd src/ && \
		echo $(MAKE) -f Makefile.packages prepare-artifacts-for-upload \
		pypi-upload;
	echo "Please upload docker_build_output/upload/build-artifacts/*:"
	ls -1 docker_build_output/upload/build-artifacts/;

packages: _prep
	cd src/; $(MAKE) -f Makefile.packages packages;

deb: _prep
	cd src/; $(MAKE) -f Makefile.packages deb;

deb-and-wheel-in-podman:
	podman run -v `pwd`:/highctidh --workdir /highctidh --rm -it debian:bookworm bash -c 'apt update && ./misc/install-debian-deps.sh && $(MAKE) wheel && CC=clang $(MAKE) deb'

wheel: _prep
	cd src/; $(MAKE) -f Makefile.packages wheel;

pytest:
	cd src; $(PYTEST) -v -n auto --doctest-modules;

update-golang-modules:
	cd src/; $(MAKE) -f Makefile.packages update-golang-modules;

install:
	cd src; $(MAKE) install;

test:
	cd src; $(MAKE) test;

test-quick:
	cd src; $(MAKE) testrandom test512;
	cd src; ./test-quick.sh

test-go:
	cd src; go test -v ./...;
	cd src/ctidh511; go test -v ./...;
	cd src/ctidh512; go test -v ./...;
	cd src/ctidh1024; go test -v ./...;
	cd src/ctidh2048; go test -v ./...;

examples:
	cd src; $(MAKE) examples;

examples-run:
	cd src; time ./example-ctidh511;
	cd src; time ./example-ctidh512;
	cd src; time ./example-ctidh1024;
	cd src; time ./example-ctidh2048;

_prep:
	-mkdir src/dist;
	-ln -s src/build build;
	-ln -s src/dist dist;

clean:
	-rm -rf build;
	-rm -rf dist;
	-rm -rf deb_dist;
	-rm -rf docker_build_output;
	-rm -rf .pytest_cache;
	-cd src; $(MAKE) clean;
	-cd src; $(MAKE) -f Makefile.packages clean;
