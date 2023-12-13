#!/bin/bash
#
# Test cross compile of c library using gcc and clang cross compilers
#
set -e;

export HOST_ARCH=`uname -m`;
CHECKMARK="\xE2\x9C\x94";

echo "Cross compile for GNU/Linux on $HOST_ARCH...";
for BITS in 511 512 1024 2048;
do
    # clang
    export ARCH=amd64;
    echo -n "$ARCH $BITS bits:";
    CC=clang \
    CFLAGS="-fPIC --target=$ARCH-pc-linux-gnu" \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

    # gcc-13 
    export ARCH=x86_64;
    echo -n "$ARCH $BITS bits:";
    CC=x86_64-linux-gnu-gcc \
    CFLAGS=-fPIC \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

    # gcc-13-aarch64-linux-gnu
    # binutils-aarch64-linux-gnu
    export ARCH=arm64;
    echo -n "$ARCH $BITS bits:";
    CFLAGS="-fPIC" \
    CC=aarch64-linux-gnu-gcc \
    AR=/usr/bin/aarch64-linux-gnu-ar \
    LD=/usr/bin/aarch64-linux-gnu-ld.gold \
    prefix=/usr/aarch64-linux-gnu/ \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

    # gcc-13-powerpc64-linux-gnu
    export ARCH=ppc64le;
    echo -n "$ARCH $BITS bits:";
    #target=ppc64le \
    CC=powerpc64le-linux-gnu-gcc \
    LD=/usr/bin/powerpc64le-linux-gnu-ld.gold \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

    # gcc-13-riscv64-linux-gnu
    export ARCH=riscv64;
    echo -n "$ARCH $BITS bits:";
    CC=$ARCH-linux-gnu-gcc \
    LD=/usr/bin/riscv64-linux-gnu-ld.gold \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

    # gcc-13-s390x-linux-gnu
    export ARCH=s390x;
    echo -n "$ARCH $BITS bits:";
    CC=$ARCH-linux-gnu-gcc \
    LD=/usr/bin/s390x-linux-gnu-ld.gold \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

    # gcc-mips64-linux-gnuabi64
    # binutils-mips64-linux-gnuabi64
    export ARCH=mips64;
    echo -n "$ARCH $BITS bits:";
    CC=$ARCH-linux-gnuabi64-gcc \
    LD=/usr/bin/mips64-linux-gnu-ld.gold \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

    # clang
    export ARCH=i386;
    echo -n "$ARCH $BITS bits:";
    CC=clang \
    CFLAGS="--target=$ARCH-pc-linux-gnu -fforce-enable-int128" \
    make;
    sha256sum *.so;
    rm *.o *.so;
    echo -e "$CHECKMARK";

done;
echo "Cross compile for successful.";
