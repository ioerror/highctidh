#!/bin/bash
set -e
set -x

echo "Installing required packages...";
apt update # > /dev/null 2>&1;
apt install -y --no-install-recommends make gcc clang git \
        python3 python3-build python3-setuptools python3-stdeb \
        build-essential python3-venv python3-wheel python3-pip \
        dh-python python3-all-dev flit fakeroot coreutils \
        python3-pytest python3-pytest-xdist #> /dev/null 2>&1;
echo "Required packages installed";
