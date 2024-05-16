#!/usr/bin/env bash

set -e -u -o pipefail

echo "Installing required packages..."
apt update  > /dev/null 2>&1
apt install -y --no-install-recommends clang make gcc git build-essential > /dev/null 2>&1
echo "Required packages installed"
