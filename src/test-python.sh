#!/usr/bin/env bash

# shellcheck disable=1091

set -e -u -o pipefail
set -x

python3 -m venv venv
source venv/bin/activate
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
