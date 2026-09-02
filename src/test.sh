#!/bin/sh

set -e;
set -x;

./test511
./test512
./test1024
./test2048
./testrandom

python3 -m venv venv
. venv/bin/activate
pip3 install build pytest pytest-xdist
mkdir -p build/src
mkdir -p dist/tmp
python3 -m build
pip install --force-reinstall dist/highctidh-*.whl

python3 ./../misc/highctidh-simple-benchmark.py
python3 -m pytest -v -n auto --doctest-modules -k 511
python3 -m pytest -v -n auto --doctest-modules -k 512
python3 -m pytest -v -n auto --doctest-modules -k 1024
python3 -m pytest -v -n auto --doctest-modules -k 2048
