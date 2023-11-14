#!/bin/bash
# This requires `apt install -y docker.io qemu-user-static
# qemu-system-{arm,mips,ppc,s390x,sparc,x86} docker-buildx`
docker pull multiarch/qemu-user-static

# This following command must be run as root or with sudo or doas, etc.
# Trusting trust hyperoperation
docker run --rm --privileged multiarch/qemu-user-static --reset -p yes

echo "Please run docker-setup.sh to finish the build setup process..."
