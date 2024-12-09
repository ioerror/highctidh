#!/bin/sh
set -e;

echo "Installing required packages...";
# XXX missing packages: python3-flit python3-build
dnf install -y gcc make clang python3.11 python3.11-devel python3.11-wheel python3-setuptools sudo time autoconf automake libtool > /dev/null 2>&1
echo "Required packages installed";
