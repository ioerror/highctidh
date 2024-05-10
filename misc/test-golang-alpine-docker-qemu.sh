#!/bin/sh
set -e;
set -x;

echo "Updating alpine package sources";
apk update > /dev/null 2>&1;
echo "Installing required packages...";
apk add --no-cache alpine-sdk clang gcc go musl-dev libc6-compat > /dev/null 2>&1;
echo "Required packages installed";
echo "Configuring git";
git config --global --add safe.directory .;
echo "Running Golang tests";
export GOCACHE=`mktemp -d`
go clean -v;
go test -v ./...;
echo "Golang tests completed";
