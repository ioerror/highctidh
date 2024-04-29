#!/bin/sh
set -e;
set -x;

echo "Installing required packages...";
swupd -y bundle-add -y sudo llvm make python3-basic # > /dev/null 2>&1;
echo "Required packages installed";
