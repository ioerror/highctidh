#!/usr/bin/env bash

set -e -u -o pipefail
set -x

echo "Installing required packages..."

apt update &> /dev/null

apt install -y --no-install-recommends \
  make git emscripten \
  &> /dev/null

echo "Required packages installed"
