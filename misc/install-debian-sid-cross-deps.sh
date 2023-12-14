#!/bin/bash
set -e

apt update
apt install -y --no-install-recommends make gcc clang git python3 \
    python3-build python3-setuptools build-essential python3-venv \
    python3-wheel python3-pip flit gcc-13 gcc-aarch64-linux-gnu \
    binutils-aarch64-linux-gnu gcc-13-powerpc64-linux-gnu \
    gcc-13-riscv64-linux-gnu gcc-13-s390x-linux-gnu \
    gcc-mips64-linux-gnuabi64 binutils-mips64-linux-gnuabi64 \
    libc6-dev-arm64-cross libc6-dev-armel-cross libc6-dev-armhf-cross \
    libc6-dev-mips64-cross libc6-dev-mips64el-cross libc6-dev-ppc64-cross \
    libc6-dev-ppc64el-cross libc6-dev-riscv64-cross libc6-dev-s390x-cross \
    libc6-dev-sparc64-cross

