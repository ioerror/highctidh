#!/bin/bash
set -e;

if test -z "$1"; then
	echo "ARCH appears to be unset"
	exit 1
fi
ARCH="$1"

# This BASE_PACKAGES list assumes a Docker image with golang installed
BASE_PACKAGES="ca-certificates clang lld make pkg-config";

case "$ARCH" in
386|i386)
	PACKAGES="gcc-i686-linux-gnu libc6-dev-i386-cross" ;;

arm|arm32v*)
	PACKAGES="gcc-arm-linux-gnueabi gcc-arm-linux-gnueabihf libc6-dev-armel-cross libc6-dev-armhf-cross" ;;

arm64)
	PACKAGES="gcc-aarch64-linux-gnu libc6-dev-arm64-cross" ;;

ppc64)
	PACKAGES="gcc-powerpc64-linux-gnu libc6-dev-powerpc-ppc64-cross" ;;

mipsle|mipsel)
	PACKAGES="gcc-mipsel-linux-gnu libc6-dev-mipsel-cross" ;;

mips64)
	PACKAGES="gcc-mips64-linux-gnuabi64 libc6-dev-mips64-cross" ;;

mips64le|mips64el)
	PACKAGES="gcc-mips64el-linux-gnuabi64 libc6-dev-mipsn32-mips64el-cross" ;;

ppc64le)
	PACKAGES="gcc-powerpc64le-linux-gnu libc6-dev-ppc64el-cross" ;;

mips|riscv64|s390x)
	PACKAGES="gcc-$ARCH-linux-gnu libc6-dev-$ARCH-cross" ;;
esac

echo "Installing required packages for $ARCH: $BASE_PACKAGES $PACKAGES"
apt update > /dev/null 2>&1
apt install -y --no-install-recommends $BASE_PACKAGES $PACKAGES > /dev/null 2>&1
echo "Required packages installed"
