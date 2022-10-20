#!/bin/bash
set -e

ARCHES="amd64 arm32v5 arm32v7 arm64v8 i386 mips64le ppc64le riscv64 s390x"

docker pull multiarch/qemu-user-static
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

for arch in $ARCHES
do
    docker pull $arch/debian:sid-slim
done

for arch in $ARCHES
do
    mkdir -p docker_build_output/$arch
    echo "Building for $arch..."
    docker build . --build-arg ARCH=$arch -t debian-$arch-libhighctidh  --progress=plain 2>&1 >> docker_build_output/$arch/build.log
    if [ $? -eq 0 ]; then
        echo "Build base image for $arch"
    fi
done

docker image ls|grep libhighctidh
