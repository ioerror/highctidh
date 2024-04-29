#!/bin/sh
set -e;

echo "Installing required packages...";
# XXX missing python3-flit python3-build
dnf install -y gcc make clang python3 python3-devel python3-wheel python3-setuptools sudo #> /dev/null 2>&1;
echo "Required packages installed";
