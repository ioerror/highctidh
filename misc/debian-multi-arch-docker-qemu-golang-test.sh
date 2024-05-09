#!/bin/sh
set -e;
set -x;

echo "Updating Debian package sources";
apt update # > /dev/null 2>&1;
echo "Installing required packages...";
apt install -y --no-install-recommends make gcc clang git \
        python3 python3-build python3-setuptools python3-stdeb \
        build-essential python3-venv python3-wheel python3-pip \
        dh-python python3-all-dev flit fakeroot coreutils \
        python3-pytest python3-pytest-xdist #> /dev/null 2>&1;
echo "Required packages installed";
echo "Configuring git";
git config --global --add safe.directory .;
echo "Running Golang tests";
#export GOEXPERIMENT=cgocheck2;
#export GODEBUG=cgocheck=1;
#export CGO_LDFLAGS="-Wl,--no-as-needed -Wl,-allow-multiple-definition";
export CC=clang
go test -v ./...;
echo "Golang tests completed";
