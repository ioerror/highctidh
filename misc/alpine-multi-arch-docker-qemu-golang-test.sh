#!/bin/sh
set -e;
set -x;

echo "Updating alpine package sources";
apk update > /dev/null 2>&1;
echo "Installing required packages...";
apk add --no-cache alpine-sdk clang go > /dev/null 2>&1;
echo "Required packages installed";
echo "Configuring git";
git config --global --add safe.directory .;
echo "Running Golang tests";
#export GOEXPERIMENT=cgocheck2;
#export GODEBUG=cgocheck=1;
#export CGO_LDFLAGS="-Wl,--no-as-needed -Wl,-allow-multiple-definition";
go test -v ./...;
echo "Golang tests completed";
