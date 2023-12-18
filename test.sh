#!/bin/bash

set -e;
set -x;

python3 -m venv venv
source venv/bin/activate
python3 -m build
pip install --force-reinstall dist/highctidh-*.whl
python3 -m unittest -v
python3 ./misc/highctidh-simple-benchmark.py
./test511
./test512
./test1024
./test2048
./testrandom
