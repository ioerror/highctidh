#!/bin/sh
set -e;

#pacman -Syu #> /dev/null 2>&1;
echo "Installing required packages...";
pacman -Syu gcc clang make python3 python-devtools python-build python-flit \
    python-setuptools python-wheel #> /dev/null 2>&1;
echo "Required packages installed";
