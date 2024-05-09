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
go clean -v;
go test -v ./...;
echo "Golang tests completed";
