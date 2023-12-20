#!/bin/bash
#
# Test cross compile of c library using clang cross compilers
#
set -eu;
set -x;

export HOST_ARCH=`uname -m`;
export CODENAME=`lsb_release -c|grep Codename`;
echo "Building on: $HOST_ARCH";
CHECKMARK="\xE2\x9C\x94";

# these are passed on to `make`:
export AR CC CC_MARCH CFLAGS LD PLATFORM PLATFORM_SIZE prefix;

make_and_clean() {
    rm -fv *.o *.so;
    echo "${PLATFORM} ${CC_MARCH:-} (${PLATFORM_SIZE}):";
    make;
    mkdir -p cross/$PLATFORM/$PLATFORM_SIZE/;
    mv -v *.so cross/$PLATFORM/$PLATFORM_SIZE/;
    echo -n "${PLATFORM} ${CC_MARCH:-} (${PLATFORM_SIZE}):";
    echo -e "$CHECKMARK";
}

echo "Checking to see if CI needs clean up";
rm -fv *.o *.so;

echo "Cross compile for GNU/Linux with clang on $HOST_ARCH...";

CC_MARCH=x86-64 \
PLATFORM=x86-64 PLATFORM_SIZE=64 \
CC="clang --target=x86_64 -fPIC -I /usr/include/x86_64-linux-gnu/ " \
make_and_clean;

PLATFORM=arm PLATFORM_SIZE=32 \
CC_MARCH=arm \
LD=/usr/bin/$PLATFORM-linux-gnueabi-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fforce-enable-int128 -fuse-ld=$LD" \
make_and_clean;

PLATFORM=ppc64 PLATFORM_SIZE=64 \
LD=/usr/powerpc64-linux-gnu/bin/ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD" \
make_and_clean;

PLATFORM=ppc64le PLATFORM_SIZE=64 \
LD=/usr/powerpc64le-linux-gnu/bin/ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD" \
make_and_clean;

PLATFORM=riscv64 PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnu-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD" \
make_and_clean;

PLATFORM=s390x PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnu-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD" \
make_and_clean;

if [ $CODENAME == "Codename:	mantic" ];
then
PLATFORM=sparcv9 PLATFORM_SIZE=64 \
AR=/usr/bin/sparc64-linux-gnu-ar \
LD=/usr/bin/sparc64-linux-gnu-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD -I /usr/sparc64-linux-gnu/" \
make_and_clean;
else
    echo "Skipping sparc 64-bit";
fi

PLATFORM=mips64 PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnuabi$PLATFORM_SIZE-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD" \
make_and_clean;

if [ $CODENAME == "Codename:	mantic" ];
then
PLATFORM=mips64el PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnuabi$PLATFORM_SIZE-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD" \
make_and_clean;
else
    echo "Skipping mips64el 64-bit";
fi

PLATFORM=i386 PLATFORM_SIZE=32 \
CC_MARCH=i686 \
LD=/usr/bin/x86_64-linux-gnu-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fforce-enable-int128 -fuse-ld=$LD" \
make_and_clean;

PLATFORM=aarch64 PLATFORM_SIZE=64 \
LD=/usr/bin/$PLATFORM-linux-gnu-ld \
CC="clang --target=$PLATFORM-pc-linux-gnu -fuse-ld=$LD" \
make_and_clean;

echo "Cross compile with clang successful:";
sha256sum cross/*/*/*.so;
