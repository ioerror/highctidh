#!/bin/bash
#
# Cross compile C libraries and optionally golang module
#
set -e

test "$1"
ARCH="$1"

GOARCH="$ARCH"
GNUTRIPLE="$ARCH-linux-gnu"
CLANG_TARGET="$ARCH-pc-linux-gnu"

case "$ARCH" in
amd64*)
	GOARCH="amd64"
	GNUTRIPLE="x86_64-linux-gnu"
	CLANG_TARGET="x86_64-pc-linux-gnu"

	V="${ARCH#amd64}"
	if test "$V" != "$ARCH"; then
		export GOAMD64="$V"
		test "${V#v}" -ge 3 && CONFIGURE_ARGS="--enable-asm"
	fi
	;;

arm|arm32v5)
	GOARCH="arm"
	export GOARM="5"
	GNUTRIPLE="arm-linux-gnueabi"
	CLANG_TARGET="arm-pc-linux-gnu"
	CFLAGS="-mfloat-abi=soft"
	;;

arm32v6|arm32v7)
	GOARCH="arm"
	export GOARM="${ARCH#arm32v}"
	GNUTRIPLE="arm-linux-gnueabi"
	CLANG_TARGET="arm-pc-linux-gnu"
	CFLAGS="-mfloat-abi=hard"
	;;

arm64*)
	GOARCH="arm64"
	GNUTRIPLE="aarch64-linux-gnu"
	CLANG_TARGET="$GNUTRIPLE"

	V="${ARCH#arm64}"
	if test "$V" != "$ARCH"; then
		export GOARM64="$V"
		V="${V#v}"
		V="${V%%.*}"
		V="${V%%,*}"
		test "$V" -ge 9 && CONFIGURE_ARGS="--enable-asm"
	fi
	;;

i386|386|i686)
	GOARCH=386
	GNUTRIPLE="i686-linux-gnu"
	CLANG_TARGET="i386-pc-linux-gnu"
	;;

mips)
	CONFIGURE_ARGS="--disable-shared"
	;;

mipsle|mipsel)
	GOARCH=mipsle
	GNUTRIPLE="mipsel-linux-gnu"
	CLANG_TARGET="mipsel-pc-linux-gnu"
	CONFIGURE_ARGS="--disable-shared"
	;;

mips64le|mips64el)
	GOARCH=mips64le
	GNUTRIPLE="mips64el-linux-gnuabi64"
	CLANG_TARGET="mips64el-pc-linux-gnu"
	;;

ppc64)
	GNUTRIPLE="powerpc64-linux-gnu"
	BUILD_LIB_WITH_GCC=1
	;;

ppc64le)
	GNUTRIPLE="powerpc64le-linux-gnu"
	;;

riscv64|s390x)
	BUILD_LIB_WITH_GCC=1
	;;
esac

if test "$BUILD_LIB_WITH_GCC"; then
	CC="${GNUTRIPLE}-gcc"
	LD="${GNUTRIPLE}-ld"
else
	CC="clang --target=$CLANG_TARGET"
	LD="ld.lld"
	LDFLAGS="-fuse-ld=lld"
fi

./configure --host="$GNUTRIPLE" --prefix="$HOME/inst" --disable-silent-rules $CONFIGURE_ARGS CC="$CC" CFLAGS="$CFLAGS" LD="$LD" LDFLAGS="$LDFLAGS"
make

make install
export PKG_CONFIG_PATH="$HOME/inst/lib/pkgconfig"

export GOARCH
export GOOS=linux
export CGO_ENABLED=1
export CC="clang --target=$CLANG_TARGET $CFLAGS"

CHECKMARK="\xE2\x9C\x94"

for BITS in 511 512 1024 2048; do
	if test "$GOAMD64"; then
		echo -n "$GOARCH/$GOAMD64 $BITS bits:"
	elif test "$GOARM"; then
		echo -n "$GOARCH/$GOARM $BITS bits:"
	else
		echo -n "$GOARCH $BITS bits:"
	fi

	( cd src/ctidh$BITS && go build )

	echo -e "$CHECKMARK"
done

test "$GOARCH" != "amd64" && exit

echo "Running tests"
for BITS in 511 512 1024 2048; do
	echo "$GOARCH $BITS bits:"

	( cd src/ctidh$BITS && go test -v )

	echo -e "$GOARCH $BITS bits:$CHECKMARK"
done
