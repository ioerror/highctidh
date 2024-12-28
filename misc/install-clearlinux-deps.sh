#!/bin/sh
set -e;
set -x;

echo "Installing required packages...";
swupd bundle-add -y dev-utils sudo llvm make python3-basic c-basic sysadmin-basic > /dev/null 2>&1
echo "Required packages installed";
