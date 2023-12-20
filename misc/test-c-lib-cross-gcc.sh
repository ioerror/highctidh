#!/bin/bash
#
# Test cross compile of c library using gcc cross compilers
#
set -eu;
set -x;

export HOST_ARCH=`uname -m`;
CHECKMARK="\xE2\x9C\x94";

# these are passed on to `make`:
export AR CC CC_MARCH CFLAGS LD PLATFORM PLATFORM_SIZE prefix;

make_and_clean() {
    rm -fv *.o *.so;
    echo "${PLATFORM} ${CC_MARCH:-} (${PLATFORM_SIZE}):";
    make;
    mkdir -p cross/$PLATFORM/$PLATFORM_SIZE/;
    mv -v *.so cross/$PLATFORM/$PLATFORM_SIZE/;
    echo -e "$CHECKMARK";
}

echo "Checking to see if CI needs clean up";
rm -fv *.o *.so;

echo "Cross compile for GNU/Linux with gcc on $HOST_ARCH...";

# See: https://gcc.gnu.org/bugzilla/show_bug.cgi?id=60846
# XXX TODO: fix this
# # gcc-13-arm-linux-gnueabi
# PLATFORM=arm PLATFORM_SIZE=32 \
# AR=/usr/bin/$PLATFORM-linux-gnueabi-ar \
# LD=/usr/bin/$PLATFORM-linux-gnueabi-ld \
# CC_MARCH=arm \
# CC="$PLATFORM-linux-gnueabi-gcc" \
# make_and_clean;

# gcc-13-powerpc64-linux-gnu
PLATFORM=ppc64 PLATFORM_SIZE=64 \
AR=/usr/bin/powerpc64-linux-gnu-ar \
LD=/usr/bin/powerpc64-linux-gnu-ld \
CC="powerpc64-linux-gnu-gcc -pipe" \
make_and_clean;

# gcc-13-powerpc64le-linux-gnu
PLATFORM=ppc64le PLATFORM_SIZE=64 \
AR=/usr/bin/powerpc64le-linux-gnu-ar \
LD=/usr/bin/powerpc64le-linux-gnu-ld \
CC="powerpc64le-linux-gnu-gcc -pipe" \
make_and_clean;

# gcc-13-riscv64-linux-gnu
PLATFORM=riscv64 PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnu-ld \
CC="$PLATFORM-linux-gnu-gcc -pipe" \
make_and_clean;

# gcc-13-s390x-linux-gnu
PLATFORM=s390x PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnu-ld \
CC="$PLATFORM-linux-gnu-gcc -pipe" \
make_and_clean;

# gcc-13-sparc64-linux-gnu
PLATFORM=sparc64 PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnu-ld \
CC="$PLATFORM-linux-gnu-gcc -pipe" \
make_and_clean;

# gcc-mips64-linux-gnuabi64
PLATFORM=mips64 PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnu-ld \
CC="$PLATFORM-linux-gnuabi64-gcc -pipe" \
make_and_clean;

# See https://gcc.gnu.org/bugzilla/show_bug.cgi?id=60846
# XXX TODO: fix this
# error: ‘__int128’ is not supported on this target
# gcc-i686-linux-gnu
# PLATFORM=i686 PLATFORM_SIZE=32 \
# LD=/usr/$PLATFORM-linux-gnu/bin/ld.gold \
# CC_MARCH=i686 \
# CC="$PLATFORM-linux-gnu-gcc -pipe" \
# make_and_clean;

# gcc-13-x86-64-linux-gnu
PLATFORM=x86_64 PLATFORM_SIZE=64 \
CC="$PLATFORM-linux-gnu-gcc -pipe" \
make_and_clean;

# gcc-13-aarch64-linux-gnu
PLATFORM=arm64 PLATFORM_SIZE=64 \
AR=/usr/bin/aarch64-linux-gnu-ar \
LD=/usr/bin/aarch64-linux-gnu-ld \
prefix=/usr/aarch64-linux-gnu/ \
CC="aarch64-linux-gnu-gcc -pipe" \
make_and_clean;

echo "Cross compile with gcc successful:";
sha256sum cross/*/*/*.so;
