#!/bin/sh
#

mkdir -p {dist/source,build/tmp};
getconf LONG_BIT
uname -s
uname -a
uname -m
git config --global --add safe.directory /home/runner/work/highctidh/highctidh;
