#!/bin/bash
set -e

echo "Installing required packages...";
apt update  > /dev/null 2>&1;
apt install -y --no-install-recommends clang fakeroot > /dev/null 2>&1;
echo "Required packages installed";
