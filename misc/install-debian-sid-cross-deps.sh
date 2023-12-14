#!/bin/bash
set -e

apt update
apt install -y --no-install-recommends make gcc clang git python3 \
    python3-build python3-setuptools build-essential python3-venv \
    python3-wheel python3-pip flit gcc-13 gcc-aarch64-linux-gnu \
    binutils-aarch64-linux-gnu gcc-13-powerpc64-linux-gnu \
    gcc-13-riscv64-linux-gnu gcc-13-s390x-linux-gnu \
    gcc-mips64-linux-gnuabi64 binutils-mips64-linux-gnuabi64
