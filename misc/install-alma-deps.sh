#!/bin/sh
set -e;

echo "Installing required packages...";
dnf install -y gcc make clang python3.11 python3.11-devel python3.11-flit python3.11-wheel python3-setuptools python3.11-build sudo #> /dev/null 2>&1;
echo "Required packages installed";
