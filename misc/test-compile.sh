#!/usr/bin/env bash

set -e -u -o pipefail
set -x

(
  cd src/
  make clean
  export CC=gcc
  make libhighctidh.a
  make libhighctidh.so
)

python3 -m build
make -f Makefile.packages deb
echo "gcc builds are okay"

(
  cd src/
  make clean
  export CC=clang
  make libhighctidh.a
  make libhighctidh.so
)

python3 -m build
make -f Makefile.packages deb
echo "clang builds are okay"

echo "cleaning up..."
(
  cd src/
  make clean
)
echo "test-compile successful"
