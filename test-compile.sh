#!/bin/bash
set -e
set -x

make clean
export CC=gcc
make libhighctidh.a
make libhighctidh.so
python3 -m build
make -f Makefile.packages deb
echo "gcc builds are okay"

make clean
export CC=clang
make libhighctidh.a
make libhighctidh.so
python3 -m build
make -f Makefile.packages deb
echo "clang builds are okay"

echo "cleaning up..."
make clean
echo "test-compile successful"
