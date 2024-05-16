#!/usr/bin/env bash

set -e -u -o pipefail

echo "Installing required packages..."

apt update &> /dev/null

apt install -y --no-install-recommends \
  clang make gcc git build-essential \
  &> /dev/null

echo "Required packages installed"
