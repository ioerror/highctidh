#!/bin/sh

apk update;
apk add --no-cache alpine-sdk python3 python3-dev py3-build py3-flit py3-setuptools py3-wheel clang;
