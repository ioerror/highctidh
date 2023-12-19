#!/bin/bash
set -e

dpkg --add-architecture i386;
apt update > /dev/null #2>&1;
echo "Installing required packages...";
apt install -y --no-install-recommends make gcc clang git build-essential \
    gcc-mipsel-linux-gnu gcc-mips64-linux-gnuabi64 gcc-mips64el-linux-gnuabi64 \
    gccgo-i686-linux-gnu gcc-powerpc64le-linux-gnu gcc-riscv64-linux-gnu \
    gcc-s390x-linux-gnu gcc-aarch64-linux-gnu gcc-arm-linux-gnueabi golang \
    ca-certificates libc6-arm64-cross libc6-armel-cross libc6-armhf-cross \
    libc6-dev-arm64-cross libc6-dev-armel-cross libc6-dev-armhf-cross \
    libc6-dev-mips32-mips64-cross libc6-dev-mips32-mips64el-cross \
    libc6-dev-mips64-cross libc6-dev-mips64el-cross libc6-dev-mips-cross \
    libc6-dev-mipsn32-mips64-cross libc6-dev-mipsn32-mips64el-cross \
    libc6-dev-mipsel-cross libc6-dev-powerpc-ppc64-cross libc6-dev-ppc64-cross \
    libc6-dev-ppc64el-cross libc6-dev-riscv64-cross libc6-dev-s390x-cross \
    libc6-dev-sparc64-cross libc6-dev-sparc-sparc64-cross libc6-dev-i386 \
    libc6-i386 libc6-dev linux-libc-dev linux-libc-dev:i386 # > /dev/null 2>&1;
echo "Required packages installed";
