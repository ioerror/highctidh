#!/bin/sh
set -e;

echo "Installing required packages...";
pacman -Syu --noconfirm gcc clang make sudo python3 python-pip python-build python-flit \
    python-setuptools python-wheel #> /dev/null 2>&1;
echo "Required packages installed";
