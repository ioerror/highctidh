#!/bin/bash
set -e

export DIST="sid"
export SOURCE_DATE_EPOCH=`git log -1 --pretty=%ct`
export DEBIAN_FRONTEND=noninteractive;
export DEB_BUILD_OPTIONS=nocheck;
export ARCHES_GCC="amd64 arm64v8 ppc64le riscv64 s390x";
export ARCHES_CLANG="arm32v5 arm32v7 mips64le i386";
echo "Starting building of libhighctidh packages: `date -R`";
echo "Currently skipping builds for $ARCHES_UNSUPPORTED";

for arch in $ARCHES_CLANG
do
  echo "Building artifacts on $arch with clang";
  mkdir -p docker_build_output/$arch/{build,dist,deb_dist};
  docker run \
    --mount type=bind,source="$(pwd)/docker_build_output/$arch/build/",target=/highctidh/build/ \
    --mount type=bind,source="$(pwd)/docker_build_output/$arch/dist/",target=/highctidh/dist/ \
    --mount type=bind,source="$(pwd)/docker_build_output/$arch/deb_dist/",target=/highctidh/deb_dist/ \
    -e "ARCH=$arch" \
    -e "CC=clang" \
    -e "DEBIAN_FRONTEND=noninteractive" \
    -e "SOURCE_DATE_EPOCH=$SOURCE_DATE_EPOCH" \
    -e "WORKDIR=/highctidh/" \
    --rm -d -it \
    debian-$arch-$DIST-libhighctidh:latest \
    /bin/bash -c 'cd /highctidh && make -f Makefile.packages packages 2>&1 >> /highctidh/dist/build.log';
done

for arch in $ARCHES_GCC
do
  echo "Building artifacts on $arch with gcc";
  mkdir -p docker_build_output/$arch/{build,dist,deb_dist};
  docker run \
    --mount type=bind,source="$(pwd)/docker_build_output/$arch/build/",target=/highctidh/build/ \
    --mount type=bind,source="$(pwd)/docker_build_output/$arch/dist/",target=/highctidh/dist/ \
    --mount type=bind,source="$(pwd)/docker_build_output/$arch/deb_dist/",target=/highctidh/deb_dist/ \
    -e "ARCH=$arch" \
    -e CC=gcc \
    -e "DEBIAN_FRONTEND=noninteractive" \
    -e "SOURCE_DATE_EPOCH=$SOURCE_DATE_EPOCH" \
    -e "WORKDIR=/highctidh/" \
    --rm -d -it \
    debian-$arch-$DIST-libhighctidh:latest \
    /bin/bash -c 'cd /highctidh && make -f Makefile.packages packages 2>&1 >> /highctidh/dist/build.log';
done

tail -f docker_build_output/*/*/build.log

find docker_build_output/ | egrep '.tar.gz$|.whl$|.deb$' | xargs -n1 sha256sum;
echo "Finished building libhighctidh packages: `date -R`";
