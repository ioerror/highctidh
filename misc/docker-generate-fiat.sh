#!/usr/bin/env bash

set -e -u -o pipefail

docker=$(type -p podman || type -p docker)

$docker run --mount type=bind,source="$(pwd)/docker_build_output/",target=/docker_build_output/ \
  --rm -it debian-libhighctidh-fiat-crypto \
  /highctidh/generate-fields.sh
