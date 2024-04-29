#!/bin/sh
set -e;

echo "Installing required packages...";
swupd -y bundle-add dev-utils sudo llvm make python3-basic # > /dev/null 2>&1;
echo "Required packages installed";
