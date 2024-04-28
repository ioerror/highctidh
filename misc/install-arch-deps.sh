#!/bin/sh
set -e;

pacman -Syu #> /dev/null 2>&1;
echo "Installing required packages...";
pacman -yS gcc clang make python3 python3-dev py3-build py3-flit \
    py3-setuptools py3-wheel #> /dev/null 2>&1;
echo "Required packages installed";
