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
while test -n "$1"; do
	case "$1" in
	--sanitizer) with_san=1 ;;
	--golang) with_go=1 ;;
	--python) with_py=1 ;;
	*) break ;;
	esac
	shift
done


RPM_PKGS="gcc make clang time autoconf automake libtool"
RPM_PKGS_PY="python3 python3-devel python3-setuptools"

DEB_PKGS="ca-certificates build-essential make gcc clang \
	time autoconf automake libtool pkg-config"
DEB_PKGS_PY="python3 python3-all-dev python3-build python3-setuptools \
	python3-venv python3-wheel python3-pip python3-stdeb dh-python \
	flit fakeroot python3-pytest python3-pytest-xdist"


PREPARE=""

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
	swupd info
	RUN="swupd bundle-add -y" ;;

debian|devuan|ubuntu|x-ubuntu*)
	export DEBIAN_FRONTEND=noninteractive

	PREPARE="apt update"
	RUN="apt install -y --no-install-recommends"

	if which sudo > /dev/null; then
		test -n "$PREPARE" && PREPARE="sudo $PREPARE"
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


case "$CONTAINER" in
almalinux)
	# XXX missing packages: python3-flit python3-build
	PKGS="$RPM_PKGS"
	test "$with_py" && PKGS="$PKGS python3.11 \
		python3.11-devel python3.11-wheel python3-setuptools"
	;;

alpine)
	PKGS="alpine-sdk bash clang llvm compiler-rt \
		autoconf automake libtool pkgconfig"
	test "$with_go" && PKGS="$PKGS go"
	test "$with_py" && PKGS="$PKGS python3 \
		python3-dev py3-build py3-flit py3-setuptools \
		py3-wheel py3-pip py3-pytest py3-pytest-xdist"
	;;

amazonlinux|oraclelinux|rockylinux)
	# XXX missing python3-flit python3-build
	# XXX some releases are also missing python3-wheel
	PKGS="$RPM_PKGS"
	test "$with_py" && PKGS="$PKGS $RPM_PKGS_PY"
	;;

archlinux)
	PKGS="gcc make clang time autoconf automake libtool"
	test "$with_py" && PKGS="$PKGS python3 python-build \
		python-flit python-setuptools python-wheel"
	;;

clearlinux)
	PKGS="dev-utils make c-basic llvm sysadmin-basic"
	test "$with_py" && PKGS="$PKGS python3-basic"
	;;

debian|devuan|ubuntu)
	;;

fedora)
	PKGS="$RPM_PKGS"
	test "$with_py" && PKGS="$PKGS $RPM_PKGS_PY \
		python3-flit python3-build python3-wheel"
	;;

opensuse)
	# XXX missing python3-flit python3-build
	# XXX some releases are also missing python3-wheel
	PKGS="$RPM_PKGS gawk diffutils"
	test "$with_py" && PKGS="$PKGS $RPM_PKGS_PY"
	;;

x-ubuntu-24.04-cross)
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
	;;

x-ubuntu-22.04-webassembly)
	PKGS="$DEB_PKGS emscripten" ;;
esac


echo "Installing required packages..."
test -n "$1" && echo "Installing extra packages: $*"

(
	$PREPARE
	$RUN $PKGS "$@"
) > /dev/null

echo "Packages installed"
