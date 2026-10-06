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
tests=${UBSAN_C_TESTS-test511 test512 test1024 test2048 testrandom}
if [ -n "$tests" ]; then
	"${MAKE:-make}" clean > /dev/null
	CFLAGS="$san ${EXTRA_CFLAGS:-}" "${MAKE:-make}" \
		-j"$(getconf _NPROCESSORS_ONLN)" $tests
fi
for t in $tests; do
	echo "ubsan C $t"
	"./$t" > /dev/null
done

cd "$root"
for b in ${UBSAN_GO_SIZES-511 512 1024 2048}; do
	echo "ubsan Go ctidh$b"
	CGO_CFLAGS="-O2 -g $san" CGO_LDFLAGS="$san" go test -count=1 "./src/ctidh$b/"
done
