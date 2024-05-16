#!/usr/bin/env bash

set -e -u -o pipefail

mkdir -p {dist/source,build/tmp}
git config --global --add safe.directory /highctidh
