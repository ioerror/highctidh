#!/bin/bash

set -e;
set -x;

python3 -m venv venv
source venv/bin/activate
pip3 install -U pip
pip3 install build pytest pytest-xdist setuptools
mkdir -p build/src
mkdir -p dist/tmp
python3 -m build
pip install --force-reinstall dist/highctidh-*.whl

python3 ./../misc/highctidh-simple-benchmark.py
python3 -m pytest -v -n auto --doctest-modules -k 512

./test512
./testrandom
