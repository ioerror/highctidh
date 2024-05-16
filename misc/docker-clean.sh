#!/usr/bin/env bash

set -e -u -o pipefail

SELECTOR=libhighctidh
IFS=" " read -r -a CONTAINERS <<< "$(
  docker image ls |
    grep $SELECTOR |
    cut -d\ -f1
)"

for c in "${CONTAINERS[@]}"
do
  echo "Removing $c..."
  docker rmi "$c"
  echo "Done removing $c"
done
