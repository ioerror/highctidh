#!/bin/sh
set -e;

#dnf update #> /dev/null 2>&1;;
echo "Installing required packages...";
dnf install -y gcc make clang python3 python3-devel python3-flit python3-wheel python3-setuptools python3-build #> /dev/null 2>&1;
echo "Required packages installed";
