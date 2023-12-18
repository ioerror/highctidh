#!/bin/bash

ARCHES="amd64 arm32v5 arm32v7 arm64v8 i386 mips64le ppc64le riscv64 s390x"
DIST="sid"

for arch in $ARCHES
do
    docker pull $arch/debian:$DIST-slim
done

for arch in $ARCHES
do
    mkdir -p docker_build_output/$arch
    echo "Building for $arch/$DIST..."
    docker build . --build-arg ARCH=$arch --build-arg DEBIAN_FRONTEND=noninteractive -t debian-$arch-$DIST-libhighctidh 2>&1 | tee docker_build_output/$arch-$DIST-docker-image-build.log
    if [ $? -eq 0 ]; then
        echo "Built base image for $arch/$DIST"
    fi
done

docker image ls|grep libhighctidh
