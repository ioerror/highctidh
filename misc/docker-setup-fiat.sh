#!/usr/bin/env bash

set -e -u -o pipefail

docker=$(type -p podman || type -p docker)

mkdir -p docker_build_output/fiat-crypto-generated-code/

$docker build misc/fiat-docker/ \
  -t debian-libhighctidh-fiat-crypto \
  --progress=plain 2>&1 |
    tee -a docker_build_output/fiat-crypto-generated-code/build.log

$docker image ls |
  grep debian-libhighctidh-fiat-crypto
