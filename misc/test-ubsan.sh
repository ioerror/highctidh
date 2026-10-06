#!/bin/bash
set -euo pipefail

root=$(cd "$(dirname "$0")/.." && pwd)
san="-fsanitize=undefined -fno-sanitize-recover=undefined"
export CC="${CC:-clang}"
export UBSAN_OPTIONS=print_stacktrace=1:halt_on_error=1

work=$(mktemp -d)
trap 'rm -rf "$work"' EXIT
cp -R "$root/src" "$work/src"
cd "$work/src"
"${MAKE:-make}" clean > /dev/null
CFLAGS="$san ${EXTRA_CFLAGS:-}" "${MAKE:-make}" -j"$(getconf _NPROCESSORS_ONLN)" \
	testrandom test512 test1024 test2048
for t in test512 test1024 test2048 testrandom; do
	echo "ubsan C $t"
	"./$t" > /dev/null
done
