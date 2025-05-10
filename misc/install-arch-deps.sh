#!/bin/sh
set -ex;

echo "Installing required packages...";
pacman -Syu --noconfirm gcc clang make sudo python3 python-build python-flit \
    python-setuptools python-wheel #> /dev/null 2>&1;
echo "Required packages installed";
