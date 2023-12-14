#!/bin/bash
set -e

apt update > /dev/null 2>&1;
echo "Installing required packages...";
apt install -y --no-install-recommends make gcc clang git python3 \
    python3-build python3-setuptools build-essential python3-venv python3-wheel \
    python3-pip flit gcc  gcc-arm-linux-gnueabi gcc-arm-linux-gnueabihf \
    gcc-aarch64-linux-gnu gcc-i686-linux-gnu gcc-mips64-linux-gnuabi64 \
    gcc-powerpc64le-linux-gnu gcc-powerpc64-linux-gnu gcc-riscv64-linux-gnu \
    gcc-s390x-linux-gnu gcc-sparc64-linux-gnu binutils-aarch64-linux-gnu \
    binutils-i686-linux-gnu gcc-mips64el-linux-gnuabi64 \
    binutils-mips64el-linux-gnuabi64 binutils-mips64-linux-gnuabi64 \
    binutils-powerpc64le-linux-gnu binutils-powerpc64-linux-gnu \
    binutils-riscv64-linux-gnu binutils-sparc64-linux-gnu \
    binutils-s390x-linux-gnu libc6-dev-arm64-cross libc6-dev-i386-cross \
    libc6-dev-armel-cross libc6-dev-armhf-cross libc6-dev-mips-cross \
    libc6-dev-mipsel-cross libc6-dev-mips64-cross libc6-dev-mips64el-cross \
    libc6-dev-ppc64-cross libc6-dev-ppc64el-cross libc6-dev-riscv64-cross \
    libc6-dev-s390x-cross libc6-dev libc6-dev-sparc64-cross > /dev/null 2>&1;
echo "Required packages installed";
