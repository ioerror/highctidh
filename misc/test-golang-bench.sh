#!/usr/bin/env bash
#
# Benchmark each golang module by field size
#

set -e -u -o pipefail

export GOOS=linux
export CGO_ENABLED=1

echo "Benchmarking on $GOOS with CGO_ENABLED=$CGO_ENABLED..."

for BITS in 511 512 1024 2048
do
  (
    cd src/ctidh$BITS
    go test -v -bench=.
  )
done

echo "Benchmark for $GOOS with CGO_ENABLED=$CGO_ENABLED successful."
