#!/bin/sh -x

export LD_LIBRARY_PATH=.
export PYTHONPATH=$PYTHONPATH:`pwd`
python3 tests/test_highctidh.py
./test511
./test512
./test1024
./test2048
./testrandom
python3 highctidh-simple-benchmark.py
