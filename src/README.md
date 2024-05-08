![highctidh build status](https://ci.codeberg.org/api/badges/vula/highctidh/status.svg "highctidh build status")

# Warning and notice

This is an unofficial fork of high-ctidh. This is highly experimental software
and it has not yet been reviewed for security considerations. We encourage
users of this software to assume that there are vulnerabilities until a review
confirms otherwise.

# Quick start

The `highctidh` project offers three compatible interfaces that all use the
same underlying C or assembler implementations:
- C Library
- Python module
- Golang module

## C library

Build and install the C library for each field size with required headers:
```
make
sudo make install
sudo ldconfig
```

After installing the libraries it is possible to use them in normal C programs.
To build and run the example C programs that use the libraries:
```
make examples
make examples-run
```

## Python bindings

Python bindings include the required shared libraries and do not require
`libhighctidh_*.so` to be installed on the system.

Quickly install the latest release from `pypi` in a virtual environment:
```
python3 -m venv venv
source venv/bin/activate
pip install highctidh
```

Build the highctidh Python bindings as a Python wheel and then install it:
```
make wheel
pip install --force-reinstall dist/highctidh-*.whl
```

Build and install a Debian GNU/Linux package containing the Python module:
```
make deb
sudo dpkg -i python3-highctidh_*.deb
```

## Golang bindings

Add the golang bindings to a golang project:
```
go get -u codeberg.org/vula/highctidh/
```

Use the golang bindings in a `.go` file by importing the field size that is
desired or by importing any or all of the four field sizes:
```
import (
    ctidh511  "codeberg.org/vula/highctidh/src/ctidh511"
    ctidh512  "codeberg.org/vula/highctidh/src/ctidh512"
    ctidh1024 "codeberg.org/vula/highctidh/src/ctidh1024"
    ctidh2048 "codeberg.org/vula/highctidh/src/ctidh2048"
) 
```

# highctidh improvements

This fork enhances high-ctidh with additional Makefile targets including
building high-ctidh as four shared libraries, one for each key size of 511,
512, 1024, and 2048. Python bindings are additionally added, as well as
optional Debian packaging of both the shared library object files and the
Python module. The Python bindings were made in concert with the author of the
Golang bindings which are now included. Both bindings were built around the
same shared objects for cross verification purposes. Currently this library is
fast on the `amd64`/`x86_64` CPU architecture and functional but much slower
with other CPU architectures. The portable backend was generated using the
`fiat-crypto` project which uses a "Correct-by-Construction" approach; see
`doc/PRIMES.md` for more information.  Tested architectures for the C library
include: `amd64`/`x86_64`  (with and without avx2), `arm32v5`, `arm32v6`,
`arm32v7`, `arm64v8`/`aarch64`/`arm64`, `i386`, `loongarch64/Loongson`, `mips`,
`mipsel`, `mips64`, `mips64el`, `POWER8/ppc64`, `POWER9/ppc64le`, `riscv64`,
`s390x`, `sun4v`, and `sparc64`.

## Golang bindings

The Golang bindings compile and should be functional on `amd64`/`x86_64`,
`arm32v5`, `arm32v6`, `arm32v7`, `arm64v8`/`aarch64`/`arm64`, `i386`, `ppc64`,
`ppc64le`, `riscv64`, `s390x`, `mips`, `mipsle`, `mips64`, `mips64le`.  The
`misc/test-golang-cross.sh` script runs tests on the host build architecture
and then attempts to cross-compile for each listed architecture.  The
`.woodpecker/golang.yml` attempts to perform a cross-compile for all listed
golang versions and enumerated architectures.  Native builds for the Golang
bindings should be functional on `loong64` and `sparc64` but this is currently
untested. Go version 1.21 is used to build and test in
`.woodpecker/golang.yml`.

## Python bindings

The Python bindings build and should be functional on `amd64`, `arm32/armv7l`,
`arm32v5`, `arm32v6`, `arm32v7`, `arm64`, `i386`, `ppc64le`, `riscv64`,
`s390x`, and `mips64el`. Python 3.9, 3.10, 3.11, and 3.12 are used to build and
test in `.woodpecker/qemu-python-clang.yml`.

Debian packages and Python wheels that contain everything needed to use
`highctidh` build with the `make -f Makefile.packages packages` Makefile target
for `amd64`, `arm32/armv7l`, `arm32/armv5`, `arm64`, `i386`, `mips64el`,
`ppc64el`, `riscv64`, and `s390x`.

## Performance

Rough performance numbers are available in `docs/BENCHMARKS.md`. We recommend
using `gcc` 10 or later as the compiler except on 32-bit platforms where we
recommend `clang` 14 or later.

## Testing

We attempt to comprehensively test changes to this software in a continuous
integration environment. Testing includes `clang`, `gcc`, native builds on
`linux/amd64` and `linux/arm64`, as well as `qemu` builds for almost all other
supported architectures. Please consult the relevant configuration files in
`.woodpecker/` for more information.

To test without installing run the `test` target:
```
   make test
```

The C library and bindings have been tested on the following operating systems:
- AlmaLinux 9.3 (GNU libc)
- Alpine v3.17 - v3.19.1 (musl libc)
- Arch latest (GNU libc)
- CheriBSD 14.0-CURRENT (FreeBSD libc)
- Clear Linux 41560 (GNU libc)
- Debian stable, testing, unstable (GNU libc)
- Devuan latest (GNU libc)
- DragonFlyBSD 6.4.0 (FreeBSD libc)
- Fedora 38, 39, 40, 41 (GNU libc)
- FreeBSD 14.0 (FreeBSD libc)
- HardenedBSD (FreeBSD libc)
- MacOS 11, 12, 13, 14 (BSD libc)
- NetBSD 10.0 (NetBSD libc)
- Omnios r151046 (illumos libc)
- OpenBSD 7.5 (OpenBSD libc)
- Oracle Linux 9 (GNU libc)
- Rockylinux 9, 9.3 (GNU libc)
- Solaris 11.4 (Solaris libc)
- Ubuntu 22.03, 23.10, 24.04 (GNU libc)
- Windows Server 2019, 2022 (MSVCRT, CYGWIN, UCRT)

## Notes on building

Building on Solaris, CheriBSD, FreeBSD, NetBSD, and OpenBSD building is
supported using the `gmake` command. GNU/Linux and MacOS are supported with the
`make` command.

MacOS 11, 12, 13, and 14 support is functional for building the C library.
MacOS 14 support is functional for the Golang bindings with Golang 1.19, 1.20,
1.21.x, and 1.22.0.
MacOS 14 supports the Python module with Python 3.9, 3.10, 3.11, and 3.12.

Windows support is extremely experimental. The Golang modules are almost
certianly not be functional on Windows due to the state of cgo on Windows
platforms. Building the main C library on Windows Server 2019 and Windows
Server 2022 should be possible with `clang` as is demonstrated in the
continuous integration configuration `windows-fiat-c-library-test.yml`.  It has
only been tested with the Windows Server 2022 image preloaded with `clang`,
`bash`, `make`, and other related tools
[https://github.com/actions/runner-images/blob/main/images/windows/Windows2022-Readme.md](available
as a part of the CI configuration).

Building the C library and performing minimal testing manually requires using
`bash` as provided by `git` on Windows, GNU `make`, and `clang` using the
following commands:
```
export HIGHCTIDH_PORTABLE=1
export WINDOWS=1
export CC=${{ matrix.CC }} MAKE=make
mkdir -p src/build/src
mkdir -p src/dist/tmp
cd src/ && make && make testrandom test512 && ./testrandom && ./test512
```

The Python module on Windows is functional when installed with `pip` under
MSYS2 using `gcc` in the `MSYS`, `UCRT64`, `MINGW64` environments or using
`clang` in the `MINGW64` environment.  These different environments are tested
in `windows-msys-64bit-gcc-cygwin-ucrt-msvcrt-python-pip-test.yml` and
`windows-msys-64bit-clang-msvcrt-python-pip-test.yml` respectively.

### Additional notes on building the C library
 
To build and install we recommend:
```
   sudo apt install gcc clang make
   make
   sudo make install
```

To build and install the shared library files using the
"Correct-by-Construction" fiat-crypto portable C backend:
```
    make libhighctidh.so HIGHCTIDH_PORTABLE=1
    sudo make install
```
The fiat-crypto portable C backend works on all platforms.

To build and install the shared library files using the original artisanal
`x86_64` assembler backend:
```
    make libhighctidh.so HIGHCTIDH_PORTABLE=0
    sudo make install
```

The original artisanal assembler backend works only on the `x86_64` platform.
It has been modified slightly for compatibility with LLVM-`as`/`clang`.
Hand written assembler contributions for other platforms are welcome.

By default `HIGHCTIDH_PORTABLE=1` is enabled for all platforms unless
the library is installed via the Python package, in which case optimized
implementations will be used where possible.


## Example C library usage

An example C program that can use any of the
libhighctidh_{511,512,1024,2048}.so libraries is available in
`src/example-ctidh.c`. Use the `make examples` target to build `example-ctidh511`,
`/example-ctidh512`, `example-ctidh1024`, and `example-ctidh2048` programs.

## Example Python module usage

A basic Python benchmarking program `misc/highctidh-simple-benchmark.py` shows
general performance numbers. Python tests may be run with `pytest` and should be
functional without `pytest` assuming the library is installed. If the library
path includes the build directory as is done in `test.sh`, `pytest` or python
should be able to run the tests without installation. 

### Additional information

More information about the Python bindings including installation instructions
are available in the `docs/README.python.md` file.

The Golang bindings behave as any normal Golang module/package.

The `fiat-crypto` code synthesis is described in `misc/fiat-docker/README.md`
and the documentation provides instructions for reproducible C code generation.

The original released README is `docs/README.original.md`.
The original website was https://ctidh.isogeny.org/software.html

# Acknowledgements

The original authors of this software released high-ctidh in the public domain.
All contributions made in this fork are also in the public domain.

Please consult `docs/AUTHORS.md` for the original authorship information.

This forked project is funded through the [NGI Assure
Fund](https://nlnet.nl/assure), a fund established by [NLnet](https://nlnet.nl)
with financial support from the European Commission's [Next Generation
Internet](https://ngi.eu) program. Learn more on the [NLnet project
page](https://nlnet.nl/project/Vula#ack).
