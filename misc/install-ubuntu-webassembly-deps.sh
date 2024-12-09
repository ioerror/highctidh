#!/bin/bash
set -e
set -x

echo "Installing required packages...";
apt update > /dev/null 2>&1;
apt install -y --no-install-recommends make git emscripten autoconf automake libtool > /dev/null 2>&1
echo "Required packages installed";
