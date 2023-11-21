This is an unofficial fork of high-ctidh. This is highly experimental software
and it has not yet been reviewed for security considerations.

This fork enhances high-ctidh with additional Makefile targets including
building high-ctidh as four shared libraries, one for each key size of 511,
512, 1024, and 2048. Python bindings are additionally added, as well as
optional Debian packaging of both the shared library object files and the
Python module.  The Python bindings were made in concert with the author of the
highctidh [ctidh_cgo](https://git.xx.network/elixxir/ctidh_cgo/) Golang
bindings. Both bindings were built around the same shared objects for cross
verification purposes. Currently this library is fast on the `x86_64` CPU
architecture and functional but much slower with other CPU architectures. The
portable backend was generated using the `fiat-crypto` project which uses a
"Correct-by-Construction" approach; see `PRIMES.md` for more information.
Tested architectures include: `amd64`, `arm32/armv7l`, `arm64/aarch64`, `i386`,
`loongarch64/Loongson`, `mips64/mips64el`, `POWER8/ppc64`, `POWER9/ppc64le`,
`riscv64`, `s390x`, `sparc64`, and `x86_64`. To see rough performance numbers,
look at `BENCHMARKS.md`. We recommend using gcc 10 or later as the compiler
except on `arm32/armv7l`, `i386`, and `mips64/mips64el` where we recommend
clang 14.

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
Hand written assembler contributions for other platforms are welcome.

By default `HIGHCTIDH_PORTABLE=1` is enabled for all platforms.

To test without installing run the `test` target:
```
   make test
```
An example C program that can use any of the
libhighctidh_{511,512,1024,2048}.so libraries is available in
`example-ctidh.c`. Use the `make examples` target to build `example-ctidh511`,
`example-ctidh512`, `example-ctidh1024`, and `example-ctidh2048` programs.

A basic Python benchmarking program `highctidh-simple-benchmark.py` shows
general performance numbers. Python tests may be run with pytest and should be
functional without pytest assuming the library is installed. If the library
path includes the build directory as is done in `test.sh`, pytest or python
should be able to run the tests without installation. 

More information about the Python bindings including installation instructions
are available in the `README.python.md` file.

The original authors of this software released high-ctidh in the public domain.
All contributions made in this fork are also in the public domain.

The original released README is `README.original.md`.
The original website was https://ctidh.isogeny.org/software.html
