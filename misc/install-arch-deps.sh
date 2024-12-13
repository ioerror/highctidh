#!/bin/sh
set -e;

echo "Installing required packages...";
pacman -Syu --noconfirm gcc clang make sudo python3 python-devtools python-build python-flit \
    python-setuptools python-wheel time autoconf automake libtool #> /dev/null 2>&1;
echo "Required packages installed";
