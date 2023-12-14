#!/bin/sh
set -e;

apk update > /dev/null 2>&1;
echo "Installing required packages...";
apk add --no-cache alpine-sdk python3 python3-dev py3-build py3-flit \
    py3-setuptools py3-wheel clang > /dev/null 2>&1;
echo "Required packages installed";
