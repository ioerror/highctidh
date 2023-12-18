#!/bin/bash -x
set -e

GID=`id -g`

docker run --mount type=bind,source="$(pwd)/docker_build_output/",target=/docker_build_output/ \
   -u $UID:$GID \
   --rm -it debian-libhighctidh-fiat-crypto:latest /highctidh/generate-fields.sh
