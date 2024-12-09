#!/bin/sh
set -e;

echo "Installing required packages...";
# XXX missing python3-flit python3-build
# XXX some releases are also missing python3-wheel
yum install -y gcc make clang python3 python3-devel python3-setuptools sudo tar gzip git time autoconf automake libtool > /dev/null 2>&1
echo "Required packages installed";
