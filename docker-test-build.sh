#!/bin/bash -x
set -e

ARCHES="amd64 arm32v5 arm32v7 arm64v8 i386 mips64le ppc64le riscv64 s390x"

for arch in $ARCHES
do
    mkdir -p docker_build_output/$arch/
    $(docker run --mount type=bind,source="$(pwd)/docker_build_output/$arch/",target=/docker_build_output/ -e ARCH=$arch -e DESTDIR=/docker_build_output/ --rm -it debian-$arch-libhighctidh:latest make install) &
    echo "Building artifacts on $arch"
done
