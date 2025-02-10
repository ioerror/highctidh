#!/bin/sh

set -e
test -n "$1"

CONTAINER="${1%%:*}"
CONTAINER="${CONTAINER%%/*}"
shift


PREPARE=""
RUN=""

case "$CONTAINER" in
almalinux|fedora|oraclelinux|rockylinux)
	RUN="dnf install -y" ;;

alpine)
	PREPARE="apk update"
	RUN="apk add --no-cache" ;;

amazonlinux)
	RUN="yum install -y" ;;

archlinux)
	RUN="pacman -Syu --noconfirm" ;;

clearlinux)
	PREPARE="swupd info"
	RUN="swupd bundle-add -y" ;;

debian|devuan|ubuntu*|x-ubuntu*)
	export DEBIAN_FRONTEND=noninteractive
	PREPARE="apt update"
	RUN="apt install -y --no-install-recommends" ;;

opensuse)
	RUN="zypper install -y" ;;

*)
	echo "Unknown container '${CONTAINER}'" >&2
	exit 1
esac


RPM_PKGS="gcc make clang time autoconf automake libtool"
RPM_PKGS_PY="python3 python3-devel python3-setuptools"

case "$CONTAINER" in
almalinux)
	# XXX missing packages: python3-flit python3-build
	PKGS="$RPM_PKGS \
		python3.11 python3.11-devel python3.11-wheel python3-setuptools" ;;

alpine)
	PKGS="alpine-sdk bash clang llvm compiler-rt \
		autoconf automake libtool pkgconfig \
		go \
		python3 python3-dev py3-build py3-flit py3-setuptools \
		py3-wheel py3-pip py3-pytest py3-pytest-xdist" ;;

amazonlinux)
	# XXX missing python3-flit python3-build
	# XXX some releases are also missing python3-wheel
	PKGS="$RPM_PKGS $RPM_PKGS_PY" ;;

archlinux)
	PKGS="gcc make clang time autoconf automake libtool \
		python3 python-build python-flit python-setuptools python-wheel" ;;

clearlinux)
	PKGS="dev-utils make c-basic llvm sysadmin-basic python3-basic" ;;

debian|devuan|ubuntu)
	PKGS="coreutils build-essential make gcc clang time \
		autoconf automake libtool \
		python3 python3-build python3-setuptools python3-stdeb \
		python3-venv python3-wheel python3-pip \
		dh-python python3-all-dev flit fakeroot \
		python3-pytest python3-pytest-xdist" ;;

fedora)
	PKGS="$RPM_PKGS $RPM_PKGS_PY \
		python3-flit python3-build python3-wheel" ;;

opensuse)
	# XXX missing python3-flit python3-build
	# XXX some releases are also missing python3-wheel
	PKGS="$RPM_PKGS gawk diffutils $RPM_PKGS_PY" ;;

oraclelinux|rockylinux)
	# XXX missing python3-flit python3-build
	# XXX some releases are also missing python3-wheel
	PKGS="$RPM_PKGS $RPM_PKGS_PY" ;;

x-ubuntu-24.10-cross)
	PKGS="make gcc clang autoconf automake libtool \
		python3 \
		python3-build python3-setuptools build-essential python3-venv \
		python3-wheel python3-pip flit gcc  gcc-arm-linux-gnueabi \
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
		libc6-dev-ppc64el-cross  libc6-dev-riscv64-cross libc6-dev-s390x-cross \
		libc6-dev libc6-dev-sparc64-cross libc6-dev-sparc-sparc64-cross" ;;

x-ubuntu-webassembly)
	PKGS="make emscripten autoconf automake libtool" ;;
esac


echo "Installing required packages..."
test -n "$2" && echo "Installing extra packages: $*"

(
	$PREPARE
	$RUN $PKGS "$@"
) > /dev/null

echo "Packages installed"
