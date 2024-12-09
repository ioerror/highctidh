#!/bin/bash
#
# Test cross compile of golang module using clang
# This is meant to be called as part of `.woodpecker/golang.yaml`
#
set -e

test -n "$1"
ARCH="$1"

TARGET="$ARCH-pc-linux-gnu"
GNUTRIPLE="$ARCH-linux-gnu"
export GOARCH="$ARCH"

case "$ARCH" in
amd64*)
	GNUTRIPLE="x86_64-linux-gnu"
	TARGET="x86_64-pc-linux-gnu"
	test "${ARCH/amd64.v}" != "$ARCH" && test "${ARCH/amd64.v}" -ge 3 && CONFIGURE_ARGS="--enable-asm"

	export GOARCH=amd64
	export GOAMD64="${ARCH/amd64.}"
	test "$GOAMD64" = "$ARCH" && unset GOAMD64
	;;

arm | arm32v*)
	GNUTRIPLE="arm-linux-gnueabi"
	TARGET="arm-pc-linux-gnu"
	CFLAGS="-mfloat-abi=soft"

	export GOARCH=arm
	export GOARM="${ARCH/arm32v}"
	test "$GOARM" = "$ARCH" && export GOARM=7
	;;

arm64*)
	GNUTRIPLE="aarch64-linux-gnu"

	export GOARM64="${ARCH/arm64}"
	test "$GOARM64" = "$ARCH" && unset GOARM64
	V="${ARCH/arm64v}"
	if test "$V" != "$ARCH"; then
		V="${V/.*}"
		V="${V/,*}"
		test "$V" -ge 9 && CONFIGURE_ARGS="--enable-asm"
	fi
	;;

i386 | 386 | i686)
	GNUTRIPLE="i686-linux-gnu"
	TARGET="i386-pc-linux-gnu"
	export GOARCH=386
	;;

mips)
	CONFIGURE_ARGS="--disable-shared"
	;;

mipsle | mipsel)
	GNUTRIPLE="mipsel-linux-gnu"
	TARGET="mipsel-pc-linux-gnu"
	CONFIGURE_ARGS="--disable-shared"
	export GOARCH=mipsle
	;;

mips64le | mips64el)
	GNUTRIPLE="mips64el-linux-gnuabi64"
	TARGET="mips64el-pc-linux-gnu"
	export GOARCH=mips64le
	;;

ppc64)
	GNUTRIPLE="powerpc64-linux-gnu"
	BUILD_LIB_WITH_GCC=1
	;;

ppc64le)
	GNUTRIPLE="powerpc64le-linux-gnu"
	;;

riscv64 | s390x)
	BUILD_LIB_WITH_GCC=1
	;;
esac

if test -n "$BUILD_LIB_WITH_GCC"; then
	CC="${GNUTRIPLE}-gcc"
	LD="${GNUTRIPLE}-ld"
else
	CC="clang --target=$TARGET"
	LD="ld.lld"
	LDFLAGS="-fuse-ld=lld"
fi

./configure --host="$GNUTRIPLE" --prefix="$HOME/inst" --disable-silent-rules $CONFIGURE_ARGS CC="$CC" CFLAGS="$CFLAGS" LD="$LD" LDFLAGS="$LDFLAGS"
make install
export PKG_CONFIG_PATH="$HOME/inst/lib/pkgconfig"

export GOOS=linux
export CGO_ENABLED=1
export CC="clang --target=$TARGET $CFLAGS"

CHECKMARK="\xE2\x9C\x94"

for BITS in 511 512 1024 2048; do
	cd src/ctidh$BITS
	if test -n "$GOAMD64"; then
		echo -n "$GOARCH/$GOAMD64 $BITS bits:"
	elif test -n "$GOARM"; then
		echo -n "$GOARCH/$GOARM $BITS bits:"
	else
		echo -n "$GOARCH $BITS bits:"
	fi
	go build
	echo -e "$CHECKMARK"
	cd -
done

if [ "$GOARCH" == "amd64" ]; then
	echo "Running tests on `uname -m`"
	for BITS in 511 512 1024 2048; do
		cd src/ctidh$BITS
		echo "$GOARCH $BITS bits:"
		go test -v
		echo -e "$GOARCH $BITS bits:$CHECKMARK"
		cd -
	done
fi
