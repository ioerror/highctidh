#!/bin/bash
set -e

export DIST="sid"
export SOURCE_DATE_EPOCH=`git log -1 --pretty=%ct`
export DEBIAN_FRONTEND=noninteractive;
export DEB_BUILD_OPTIONS=nocheck;
export ARCHES_UNSUPPORTED="POWER8/ppc64 loongarch64/Loongson sparc64";
export ARCHES_GCC="amd64 arm64v8 ppc64le riscv64 s390x";
export ARCHES_CLANG="i386 mips64le arm32v5 arm32v7";
echo "Starting building of libhighctidh packages: `date -R`";
echo "Currently skipping builds for $ARCHES_UNSUPPORTED";

if [ -d docker_build_output ];
then
    echo "Moving old docker_build_output...";
    mv docker_build_output docker_build_output.old-`date +%s`;
fi

for arch in $ARCHES_CLANG $ARCHES_GCC
do
  echo "Building artifacts on $arch with clang";
  mkdir -p docker_build_output/$arch/{build,dist,deb_dist};
done

for arch in $ARCHES_CLANG
do
  echo "Building artifacts on $arch with clang";
  CC=clang PODMAN_ARCH=$arch make -f Makefile.packages deb-and-wheel-in-podman-arch-sid
done

for arch in $ARCHES_GCC
do
  echo "Building artifacts on $arch with gcc";
  CC=gcc PODMAN_ARCH=$arch make -f Makefile.packages deb-and-wheel-in-podman-arch-sid
done

echo "Watch the build process: tail -f docker_build_output/*/*/build.log";

find docker_build_output/ | egrep '.tar.gz$|.whl$|.deb$' | xargs -n1 sha256sum;
echo "Finished building libhighctidh packages: `date -R`";
