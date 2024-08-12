# This is a minimal Makefile that is called by bmake as gmake will autodetect
# and use GNUmakefile
#
export CC ?= clang
export MAKE ?= make

library: _prep
	$(MAKE) -C src -f BSDmakefile
	cp src/*.so dist/

install test examples:
	$(MAKE) -C src -f BSDmakefile $@

test-quick:
	$(MAKE) -C src -f BSDmakefile testrandom test512
	`pwd`/src/$@.sh

examples-run:
	cd src; time ./example-ctidh511
	cd src; time ./example-ctidh512
	cd src; time ./example-ctidh1024
	cd src; time ./example-ctidh2048

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
