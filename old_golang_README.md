
What is this?
=============

CGO Go bindings to the CTIDH reference implementation.
CTIDH is a post quantum cryptographic primitive called a NIKE,
a noninteractive key exchange.

Learn more about CTIDH: https://ctidh.isogeny.org/


How to Build
============

Step 1
------

Get highctidh:

```
git clone https://codeberg.org/io/highctidh.git
```

Step 2
------

Build the CTIDH C shared library files, portable:

```
cd highctidh
make libhighctidh_511.so libhighctidh_512.so libhighctidh_1024.so libhighctidh_2048.so
sudo make install

cd ..
```

Step 3
------

Build your Go application.


The the header file in place and these environment variables sets you
should now be able to build your Go application which imports and
makes use of the CTIDH Golang bindings.

Here's how we generate the bindings for each key size:

```
export CTIDH_BITS=511
go run gen/main.go gen/bench.go gen/tests.go gen/binding.go -name=Ctidh511 -type=ctidh511 -bits=511
export CTIDH_BITS=512
go run gen/main.go gen/bench.go gen/tests.go gen/binding.go -name=Ctidh512 -type=ctidh512 -bits=512
export CTIDH_BITS=1024
go run gen/main.go gen/bench.go gen/tests.go gen/binding.go -name=Ctidh1024 -type=ctidh1024 -bits=1024
export CTIDH_BITS=2048
go run gen/main.go gen/bench.go gen/tests.go gen/binding.go -name=Ctidh2048 -type=ctidh2048 -bits=2048

```

Let $P point to the ctidh_cgo directory:

```
export P=/home/human/code/ctidh_cgo
```

Run tests for a particular key size like this:

```
export CTIDH_BITS=511
export CGO_CFLAGS="-w -g -I${P} -I${P}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${P}/highctidh -Wl,-rpath,${P}/highctidh -lhighctidh_${CTIDH_BITS}"
go test -v -tags Ctidh511

export CTIDH_BITS=512
export CGO_CFLAGS="-w -g -I${P} -I${P}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${P}/highctidh -Wl,-rpath,${P}/highctidh -lhighctidh_${CTIDH_BITS}"
go test -v -tags Ctidh512

export CTIDH_BITS=1024
export CGO_CFLAGS="-w -g -I${P} -I${P}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${P}/highctidh -Wl,-rpath,${P}/highctidh -lhighctidh_${CTIDH_BITS}"
go test -v -tags Ctidh1024

export CTIDH_BITS=2048
export CGO_CFLAGS="-w -g -I${P} -I${P}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${P}/highctidh -Wl,-rpath,${P}/highctidh -lhighctidh_${CTIDH_BITS}"
go test -v -tags Ctidh2048
```

CTIDH Tests and Benchmarks
===========================

In order to run the unit tests you'll have to set the CGO CFLAGS and
LDFLAGS to indicate the absolute path to the library and header
files. Here's an example using the LD_LIBRARY_PATH environment
variable:

```
export CTIDH_BITS=512
cp binding${CTIDH_BITS}.h binding.h
export PWD=`pwd`
export CGO_CFLAGS="-w -g -I${PWD}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${PWD}/highctidh -l:libhighctidh_${CTIDH_BITS}.so"
export LD_LIBRARY_PATH="${PWD}/highctidh"
go test -v
```

It's also possible to compile your cgo binary using a set rpath which
instructs it to load libraries from a relative path instead of setting
LD_LIBRARY_PATH:

```
export CTIDH_BITS=512
cp binding${CTIDH_BITS}.h binding.h
export PWD=`pwd`
export CGO_CFLAGS="-w -g -I${PWD}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${PWD}/highctidh -Wl,-rpath,./highctidh -lhighctidh_${CTIDH_BITS}"
go test -v
```


benchmarks
----------

Benchmark the DeriveSecret function for each public key size:

```
VALID_BIT_SIZES=('511' '512' '1024' '2048')
for bits in "${VALID_BIT_SIZES[@]}"
do
export CTIDH_BITS=$bits
cp binding${CTIDH_BITS}.h binding.h
export PWD=`pwd`
export CGO_CFLAGS="-w -g -I${PWD}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${PWD}/highctidh -Wl,-rpath,./highctidh -lhighctidh_${CTIDH_BITS}"
go test -bench=DeriveSecret
done

```


test vectors
------------

Test vectors are a work in progress.

```
VALID_BIT_SIZES=('511' '512' '1024' '2048')
for bits in "${VALID_BIT_SIZES[@]}"
do
export CTIDH_BITS=$bits
cp binding${CTIDH_BITS}.h binding.h
export PWD=`pwd`
export CGO_CFLAGS="-w -g -I${PWD}/highctidh -DBITS=${CTIDH_BITS}"
export CGO_LDFLAGS="-L${PWD}/highctidh -Wl,-rpath,./highctidh -lhighctidh_${CTIDH_BITS}"
go test -v -tags=bits${CTIDH_BITS} -run=${CTIDH_BITS}
done
```


License
=======

This is public domain.
