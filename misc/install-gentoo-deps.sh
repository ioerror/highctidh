#!/bin/sh
set -e;

echo "Installing required packages...";
emerge -y sys-devel/gcc #> /dev/null 2>&1;
emerge -y app-admin/sudo #> /dev/null 2>&1;
emerge -y dev-lang/python:3.11 #> /dev/null 2>&1;
echo "Required packages installed";
