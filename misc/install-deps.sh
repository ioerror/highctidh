#!/bin/sh

set -e
test -n "$1"

CONTAINER_FULL="$1"
CONTAINER="${CONTAINER_FULL%%:*}"
CONTAINER="${CONTAINER%%/*}"
shift

with_san=""
with_go=""
with_py=""
only=""
while test -n "$1"; do
	case "$1" in
	--sanitizer) with_san=1 ;;
	--golang) with_go=1 ;;
	--python) with_py=1 ;;
	--only)
		if ! test "$2"; then
			echo "--only must be followed by package names" >&2
			exit 1
		fi
		only=1
		;;
	*) break ;;
	esac
	shift
done


RPM_PKGS="gcc make clang time autoconf automake libtool"
RPM_PKGS_PY="python3-devel python3-setuptools"

DEB_PKGS="ca-certificates build-essential make gcc clang \
	time autoconf automake libtool pkg-config"
DEB_PKGS_PY="python3-all-dev python3-build python3-setuptools \
	python3-venv python3-wheel python3-pip python3-stdeb dh-python \
	flit fakeroot python3-pytest python3-pytest-xdist"


PREPARE=""

case "$CONTAINER" in
almalinux|amazonlinux|fedora|oraclelinux|rockylinux)
	RUN="dnf install -y"
	PKGS="$RPM_PKGS"
	test "$with_py" && PKGS="$PKGS $RPM_PKGS_PY"
	;;

alpine)
	PREPARE="apk update"
	RUN="apk add --no-cache" ;;

archlinux)
	RUN="pacman -Syu --noconfirm" ;;

clearlinux)
	swupd info
	RUN="swupd bundle-add -y" ;;

debian|devuan|ubuntu|x-ubuntu*)
	export DEBIAN_FRONTEND=noninteractive

	PREPARE="apt update"
	RUN="apt install -y --no-install-recommends"

	if which sudo > /dev/null; then
		PREPARE="sudo $PREPARE"
		RUN="sudo $RUN"
	fi

	PKGS="$DEB_PKGS"
	test "$with_san" && PKGS="$PKGS libclang-rt-dev"
	test "$with_go" && PKGS="$PKGS golang"
	test "$with_py" && PKGS="$PKGS $DEB_PKGS_PY"
	;;

opensuse)
	RUN="zypper install -y" ;;

*)
	echo "Unknown container '${CONTAINER}'" >&2
	exit 1
esac


if test "$PREPARE"; then
	echo "Preparing with: $PREPARE"
	$PREPARE > /dev/null
fi


if test "$only"; then
	echo "Installing only packages: $*"
	$RUN "$@"
	exit
fi


case "$CONTAINER" in
almalinux)
	# XXX : python3-flit python3-build
	test "$with_py" && PKGS="$PKGS python3-wheel"
	;;

alpine)
	PKGS="alpine-sdk bash clang llvm compiler-rt \
		autoconf automake libtool pkgconfig"
	test "$with_go" && PKGS="$PKGS go"
	test "$with_py" && PKGS="$PKGS python3-dev \
		py3-build py3-flit py3-setuptools \
		py3-wheel py3-pip py3-pytest py3-pytest-xdist"
	;;

amazonlinux)
	# XXX 2023.6.20250203.1 missing python3-wheel
	test "$with_py" && PKGS="$PKGS python3-flit python3-build"
	;;

archlinux)
	PKGS="gcc make clang time autoconf automake libtool"
	test "$with_py" && PKGS="$PKGS python-build \
		python-flit python-setuptools python-wheel"
	;;

clearlinux)
	PKGS="dev-utils make c-basic llvm sysadmin-basic"
	test "$with_py" && PKGS="$PKGS python3-basic"
	;;

fedora)
	test "$with_py" && PKGS="$PKGS \
		python3-flit python3-build python3-wheel"
	;;

opensuse)
	# XXX missing python3-flit python3-build
	# XXX some releases are also missing python3-wheel
	PKGS="$RPM_PKGS gawk diffutils"
	test "$with_py" && PKGS="$PKGS $RPM_PKGS_PY"
	;;

oraclelinux|rockylinux)
	# XXX missing python3-flit python3-build
	# XXX some releases are also missing python3-wheel
	;;

x-ubuntu-22.04-webassembly)
	PKGS="$DEB_PKGS emscripten"
	;;

x-ubuntu-24.04-cross)
	PACKAGES=""
	case "$1" in
	386|i386)
		PACKAGES="gcc-i686-linux-gnu libc6-dev-i386-cross" ;;

	amd64*)
		PACKAGES=" " ;;

	arm|arm32v*)
		PACKAGES="gcc-arm-linux-gnueabi libc6-dev-armel-cross libc6-dev-armhf-cross" ;;

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
	if test "$PACKAGES"; then
		shift
		PKGS="$PKGS $PACKAGES"
	else
		PKGS="$PKGS \
		gcc gcc-arm-linux-gnueabi \
		gcc-arm-linux-gnueabihf gcc-aarch64-linux-gnu gcc-i686-linux-gnu \
		gcc-mips64-linux-gnuabi64 gcc-powerpc64le-linux-gnu \
		gcc-powerpc64-linux-gnu gcc-riscv64-linux-gnu gcc-s390x-linux-gnu \
		gcc-sparc64-linux-gnu binutils-aarch64-linux-gnu binutils-i686-linux-gnu \
		binutils-mips64el-linux-gnuabi64 binutils-mips64-linux-gnuabi64 \
		binutils-mips-linux-gnu binutils-mipsel-linux-gnu \
		binutils-powerpc64le-linux-gnu binutils-powerpc64-linux-gnu \
		binutils-riscv64-linux-gnu binutils-sparc64-linux-gnu \
		binutils-s390x-linux-gnu libc6-dev-arm64-cross \
		libc6-dev-i386-cross libc6-dev-armel-cross libc6-dev-armhf-cross \
		libc6-dev-mips-cross libc6-dev-mipsel-cross libc6-dev-mipsn32-mips64el-cross \
		libc6-dev-mips64-cross libc6-dev-mips64el-cross libc6-dev-ppc64-cross \
		libc6-dev-ppc64el-cross libc6-dev-riscv64-cross libc6-dev-s390x-cross \
		libc6-dev libc6-dev-sparc64-cross libc6-dev-sparc-sparc64-cross"
	fi
	;;
esac


echo "Installing required packages..."
test -n "$1" && echo "Installing extra packages: $*"

$RUN $PKGS "$@" > /dev/null

echo "Packages installed"
