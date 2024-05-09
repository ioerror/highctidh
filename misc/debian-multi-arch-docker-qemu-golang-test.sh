#!/bin/sh
set -e;
set -x;

echo "Updating Debian package sources";
apt update # > /dev/null 2>&1;
echo "Installing required packages...";
apt install -y --no-install-recommends make clang git \
  golang build-essential fakeroot coreutils ca-certificates #> /dev/null 2>&1;
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
