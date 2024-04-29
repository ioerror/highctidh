#!/bin/sh
set -e;

echo "Installing required packages...";
sudo swupd bundle-add llvm # > /dev/null 2>&1;
sudo swupd bundle-add make # > /dev/null 2>&1;
sudo swupd bundle-add python3-basic # > /dev/null 2>&1;
echo "Required packages installed";
