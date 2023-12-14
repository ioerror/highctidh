#!/bin/bash
#
# Test cross compile of c library using gcc and clang cross compilers
#
# XXX TODO:
# 
#  This should cross compile _each_ target with clang and also with gcc
#
set -eu;

export HOST_ARCH=`uname -m`;
CHECKMARK="\xE2\x9C\x94";

# these are passed on to `make`:
export AR CC CC_MARCH CFLAGS LD PLATFORM PLATFORM_SIZE prefix

make_and_clean() {
    echo -n "${PLATFORM} ${CC_MARCH:-} (${PLATFORM_SIZE} bit) ctidh-$BITS:";
    make;
    rm *.o;
    mkdir -p cross/$PLATFORM/$PLATFORM_SIZE/;
    mv *.so cross/$PLATFORM/$PLATFORM_SIZE/;
    echo -e "$CHECKMARK";
}

echo "Cross compile for GNU/Linux on $HOST_ARCH...";
for BITS in 511 512 1024 2048;
do
    # clang
    PLATFORM=amd64 PLATFORM_SIZE=64 \
    CC=clang \
    CFLAGS="--target=$PLATFORM-pc-linux-gnu" \
    make_and_clean

    # gcc-13 TODO missing --target for this
    PLATFORM=x86_64 PLATFORM_SIZE=64 \
    CC=x86_64-linux-gnu-gcc \
    make_and_clean

    # gcc-13-aarch64-linux-gnu
    # binutils-aarch64-linux-gnu
    PLATFORM=arm64 PLATFORM_SIZE=64 \
    CC=aarch64-linux-gnu-gcc \
    AR=/usr/bin/aarch64-linux-gnu-ar \
    LD=/usr/bin/aarch64-linux-gnu-ld.gold \
    prefix=/usr/aarch64-linux-gnu/ \
    make_and_clean

    # gcc-13-powerpc64-linux-gnu
    #target=ppc64le \
    PLATFORM=ppc64le PLATFORM_SIZE=64 \
    CC=powerpc64le-linux-gnu-gcc \
    LD=/usr/bin/powerpc64le-linux-gnu-ld.gold \
    make_and_clean

    # gcc-13-riscv64-linux-gnu
    PLATFORM=riscv64 PLATFORM_SIZE=64 \
    CC=$PLATFORM-linux-gnu-gcc \
    LD=/usr/bin/riscv64-linux-gnu-ld.gold \
    make_and_clean

    # gcc-13-s390x-linux-gnu
    PLATFORM=s390x PLATFORM_SIZE=64 \
    CC=$PLATFORM-linux-gnu-gcc \
    LD=/usr/bin/s390x-linux-gnu-ld.gold \
    make_and_clean

    # gcc-mips64-linux-gnuabi64
    # binutils-mips64-linux-gnuabi64
    PLATFORM=mips64 PLATFORM_SIZE=64 \
    CC=$PLATFORM-linux-gnuabi64-gcc \
    LD=/usr/bin/mips64-linux-gnu-ld.gold \
    make_and_clean

    # clang
    PLATFORM=i386 PLATFORM_SIZE=32 \
    CC=clang \
    CC_MARCH=i686 CFLAGS="--target=$PLATFORM-pc-linux-gnu" \
    make_and_clean

done;
    sha256sum cross/*/*.so;
echo "Cross compile for successful.";
