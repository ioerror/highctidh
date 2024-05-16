#!/bin/bash
set -e
set -x

cd src/
make clean
export CC=gcc
make libhighctidh.a
make libhighctidh.so
cd ../
python3 -m build
make -f Makefile.packages deb
echo "gcc builds are okay"

cd src/
make clean
export CC=clang
make libhighctidh.a
make libhighctidh.so
cd ../
python3 -m build
make -f Makefile.packages deb
echo "clang builds are okay"

echo "cleaning up..."
cd src/
make clean
cd ../
echo "test-compile successful"
