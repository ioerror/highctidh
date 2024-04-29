#!/bin/sh
set -e;

echo "Installing required packages...";
mkdir -p /var/db/repos/gentoo;
emerge sys-devel/gcc #> /dev/null 2>&1;
emerge app-admin/sudo #> /dev/null 2>&1;
emerge dev-lang/python:3.11 #> /dev/null 2>&1;
echo "Required packages installed";
