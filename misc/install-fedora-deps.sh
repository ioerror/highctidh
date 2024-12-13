#!/bin/sh
set -e;

echo "Installing required packages...";
dnf install -y gcc make clang python3 python3-devel python3-flit python3-wheel python3-setuptools python3-build sudo time autoconf automake libtool > /dev/null 2>&1;
echo "Required packages installed";
