#!/bin/sh
set -e;

echo "Installing required packages...";
mkdir /var/lock/
opkg update #> /dev/null 2>&1;
opkg install gcc clang make sudo python3 #> /dev/null 2>&1;
echo "Required packages installed";
