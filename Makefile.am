ACLOCAL_AMFLAGS = -I m4
CFLAGS = $(AM_CFLAGS) -Isrc


include_HEADERS = \
	src/highctidh.h

lib_LTLIBRARIES = \
	libhighctidh_511.la \
	libhighctidh_512.la \
	libhighctidh_1024.la \
	libhighctidh_2048.la \
	libhighctidh.la

pkgconfigdir = $(libdir)/pkgconfig
pkgconfig_DATA = \
	libhighctidh_511.pc \
	libhighctidh_512.pc \
	libhighctidh_1024.pc \
	libhighctidh_2048.pc \
	libhighctidh.pc

EXTRA_DIST = README.md

CLEANFILES = \
	libbase.la \
	libcommon511.la \
	libcommon512.la \
	libcommon1024.la \
	libcommon2048.la \
	libhighctidh_511.pc \
	libhighctidh_512.pc \
	libhighctidh_1024.pc \
	libhighctidh_2048.pc \
	libhighctidh.pc \
	__pycache__ \
	src/__pycache__

check_PROGRAMS = \
	src/example-ctidh511 \
	src/example-ctidh512 \
	src/example-ctidh1024 \
	src/example-ctidh2048 \
	src/test511 \
	src/test512 \
	src/test1024 \
	src/test2048 \
	src/testrandom


EXTRA_LTLIBRARIES = libbase.la libcommon511.la libcommon512.la libcommon1024.la libcommon2048.la

libbase_la_SOURCES = \
	src/crypto_classify.c \
	src/crypto_declassify.c \
	src/randombytes.c \
	src/int32_sort.c \
	src/random.c \
	src/steps.c \
	src/steps_untuned.c

libcommon511_la_CFLAGS = -DBITS=511
if ENABLE_ASM
libcommon511_la_SOURCES = src/uintbig511.S src/fp511.S
else
libcommon511_la_SOURCES = src/fiat_p511.c src/fp2fiat.c
endif
libcommon511_la_SOURCES += src/fp_inv511.c src/fp_sqrt511.c src/primes511.c \
	src/csidh.c src/elligator.c src/mont.c src/poly.c src/skgen.c \
	src/validate.c

libcommon512_la_CFLAGS = -DBITS=512
if ENABLE_ASM
libcommon512_la_SOURCES = src/uintbig512.S src/fp512.S
else
libcommon512_la_SOURCES = src/fiat_p512.c src/fp2fiat.c
endif
libcommon512_la_SOURCES += src/fp_inv512.c src/fp_sqrt512.c src/primes512.c \
	src/csidh.c src/elligator.c src/mont.c src/poly.c src/skgen.c \
	src/validate.c

libcommon1024_la_CFLAGS = -DBITS=1024
if ENABLE_ASM
libcommon1024_la_SOURCES = src/uintbig1024.S src/fp1024.S
else
libcommon1024_la_SOURCES = src/fiat_p1024.c src/fp2fiat.c
endif
libcommon1024_la_SOURCES += src/fp_inv1024.c src/fp_sqrt1024.c src/primes1024.c \
	src/csidh.c src/elligator.c src/mont.c src/poly.c src/skgen.c \
	src/validate.c

libcommon2048_la_CFLAGS = -DBITS=2048
if ENABLE_ASM
libcommon2048_la_SOURCES = src/uintbig2048.S src/fp2048.S
else
libcommon2048_la_SOURCES = src/fiat_p2048.c src/fp2fiat.c
endif
libcommon2048_la_SOURCES += src/fp_inv2048.c src/fp_sqrt2048.c src/primes2048.c \
	src/csidh.c src/elligator.c src/mont.c src/poly.c src/skgen.c \
	src/validate.c


test511: src/test511$(EXEEXT)
src_test511_CFLAGS = -DBITS=511
src_test511_LDADD = libbase.la libcommon511.la
src_test511_LDFLAGS = $(LTPROGLDFLAGS)
src_test511_SOURCES = src/test.c

test512: src/test512$(EXEEXT)
src_test512_CFLAGS = -DBITS=512
src_test512_LDADD = libbase.la libcommon512.la
src_test512_LDFLAGS = $(LTPROGLDFLAGS)
src_test512_SOURCES = src/test.c

test1024: src/test1024$(EXEEXT)
src_test1024_CFLAGS = -DBITS=1024
src_test1024_LDADD = libbase.la libcommon1024.la
src_test1024_LDFLAGS = $(LTPROGLDFLAGS)
src_test1024_SOURCES = src/test.c

test2048: src/test2048$(EXEEXT)
src_test2048_CFLAGS = -DBITS=2048
src_test2048_LDADD = libbase.la libcommon2048.la
src_test2048_LDFLAGS = $(LTPROGLDFLAGS)
src_test2048_SOURCES = src/test.c

src_testrandom_LDADD = libbase.la


libhighctidh_511_la_CFLAGS = -DBITS=511
libhighctidh_511_la_SOURCES =
libhighctidh_511_la_LDFLAGS = $(LTLDFLAGS) -export-symbols-regex '^highctidh_511_(base|public_key_(from|to)_bytes|csidh(|_private(|_withrng))|validate|action)$$'
libhighctidh_511_la_LIBADD = libbase.la libcommon511.la

libhighctidh_512_la_CFLAGS = -DBITS=512
libhighctidh_512_la_SOURCES =
libhighctidh_512_la_LDFLAGS = $(LTLDFLAGS) -export-symbols-regex '^highctidh_512_(base|public_key_(from|to)_bytes|csidh(|_private(|_withrng))|validate|action)$$'
libhighctidh_512_la_LIBADD = libbase.la libcommon512.la

libhighctidh_1024_la_CFLAGS = -DBITS=1024
libhighctidh_1024_la_SOURCES =
libhighctidh_1024_la_LDFLAGS = $(LTLDFLAGS) -export-symbols-regex '^highctidh_1024_(base|public_key_(from|to)_bytes|csidh(|_private(|_withrng))|validate|action)$$'
libhighctidh_1024_la_LIBADD = libbase.la libcommon1024.la

libhighctidh_2048_la_CFLAGS = -DBITS=2048
libhighctidh_2048_la_SOURCES =
libhighctidh_2048_la_LDFLAGS = $(LTLDFLAGS) -export-symbols-regex '^highctidh_2048_(base|public_key_(from|to)_bytes|csidh(|_private(|_withrng))|validate|action)$$'
libhighctidh_2048_la_LIBADD = libbase.la libcommon2048.la

libhighctidh_la_SOURCES =
libhighctidh_la_LDFLAGS = $(LTLDFLAGS) -export-symbols-regex '^highctidh_(511|512|1024|2048)_(base|public_key_(from|to)_bytes|csidh(|_private(|_withrng))|validate|action)$$'
libhighctidh_la_LIBADD = libbase.la \
	libcommon511.la \
	libcommon512.la \
	libcommon1024.la \
	libcommon2048.la

example511: src/example-ctidh511$(EXEEXT)
src_example_ctidh511_CFLAGS = -DBITS=511
src_example_ctidh511_LDADD = libhighctidh_511.la
src_example_ctidh511_LDFLAGS = $(LTPROGLDFLAGS)
src_example_ctidh511_SOURCES = src/example-ctidh.c

example512: src/example-ctidh512$(EXEEXT)
src_example_ctidh512_CFLAGS = -DBITS=512
src_example_ctidh512_LDADD = libhighctidh_512.la
src_example_ctidh512_LDFLAGS = $(LTPROGLDFLAGS)
src_example_ctidh512_SOURCES = src/example-ctidh.c

example1024: src/example-ctidh1024$(EXEEXT)
src_example_ctidh1024_CFLAGS = -DBITS=1024
src_example_ctidh1024_LDADD = libhighctidh_1024.la
src_example_ctidh1024_LDFLAGS = $(LTPROGLDFLAGS)
src_example_ctidh1024_SOURCES = src/example-ctidh.c

example2048: src/example-ctidh2048$(EXEEXT)
src_example_ctidh2048_CFLAGS = -DBITS=2048
src_example_ctidh2048_LDADD = libhighctidh_2048.la
src_example_ctidh2048_LDFLAGS = $(LTPROGLDFLAGS)
src_example_ctidh2048_SOURCES = src/example-ctidh.c


examples: example511 example512 example1024 example2048
	ls -l src/example-ctidh511$(EXEEXT) src/example-ctidh512$(EXEEXT) src/example-ctidh1024$(EXEEXT) src/example-ctidh2048$(EXEEXT)

examples-run: examples
	time src/example-ctidh511
	time src/example-ctidh512
	time src/example-ctidh1024
	time src/example-ctidh2048

test: src/test511 src/test512 src/test1024 src/test2048 src/testrandom
	cd src ; ./test.sh

test-quick: src/test512 src/testrandom
	src/test512
	src/testrandom

test-go:
	cd src; go test -v ./...
	cd src/ctidh511; go test -v ./...
	cd src/ctidh512; go test -v ./...
	cd src/ctidh1024; go test -v ./...
	cd src/ctidh2048; go test -v ./...


PYTEST ?= pytest-3
pytest:
	$(PYTEST) --rootdir=src -v -n auto --doctest-modules


# This produces wasm as a shared library module
wasm: libhighctidh.la
	emcc -v --no-entry -O1 .libs/libhighctidh.a -sSIDE_MODULE=1 -o src/dist/highctidh.wasm
	ls -al src/dist/highctidh.wasm













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

deb-and-wheel-in-podman:
	podman run --rm -it \
	  -v `pwd`:/highctidh \
	  --workdir /highctidh \
	  -e "CC=$(CC)" \
	  -e "DEBIAN_FRONTEND=noninteractive" \
	  -e "DEB_BUILD_OPTIONS=nocheck" \
	  debian:bookworm \
	    bash -c 'apt update && \
	      ./misc/install-debian-deps.sh && \
	      $(MAKE) wheel && \
	      $(MAKE) deb'


packages deb wheel: _prep
	$(MAKE) -C src -f Makefile.packages $@

update-golang-modules test-python:
	$(MAKE) -C src -f Makefile.packages $@



alpine-multi-arch-deps:
	@echo 'Attempting to install dependencies'
	uname -a
	./misc/install-alpine-deps.sh
	@echo 'Attempting to build, install, and test'

ALPINE_DEPS := \
  alpine-multi-arch-deps \
  install \
  examples-run \
  test-python

alpine-multi-arch: $(ALPINE_DEPS)
	@echo 'Hopefully the above was successful'
	ls -alh src/dist/*.so && sha256sum src/dist/*.so

_prep:
	-mkdir src/dist
	-ln -s src/build build
	-ln -s src/dist dist

xxclean:
	-$(RM) -r build
	-$(RM) -r dist
	-$(RM) -r deb_dist
	-$(RM) -r docker_build_output
	-$(RM) -r .pytest_cache
	-$(MAKE) -C src clean
	-$(MAKE) -C src -f Makefile.packages clean
