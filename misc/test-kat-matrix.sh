#!/bin/bash

set -u;

OUT=$(realpath -m "${1:?usage: $0 OUTDIR [ARCH...]}");
shift;
ARCHES=${*:-amd64 386 arm64 armhf s390x ppc64 ppc64le mips riscv64};
ROOT=$(cd "$(dirname "$0")/.." && pwd);
SIZES="511 512 1024 2048";
CSRC="crypto_classify.c crypto_declassify.c csidh.c elligator.c fp2fiat.c
mont.c poly.c randombytes.c random.c skgen.c steps.c steps_untuned.c
validate.c int32_sort.c";

mkdir -p "$OUT/log" "$OUT/bin" "$OUT/gen" "$OUT/result";

arch_env() {
    GOARM="";
    CGO_ALLOW="";
    CGO_LD="";
    PYPRE="";
    QEMU="";
    TRIPLE="";
    SIZE=64;
    EXTRA="";
    case $1 in
        amd64) GOARCH=amd64; CC=gcc;;
        386) GOARCH=386; TRIPLE=i686-linux-gnu; QEMU=qemu-i386; SIZE=32;
             PYPRE="setarch i686";;
        arm64) GOARCH=arm64; TRIPLE=aarch64-linux-gnu; QEMU=qemu-aarch64;;
        armhf) GOARCH=arm; GOARM=7; TRIPLE=arm-linux-gnueabihf;
               QEMU=qemu-arm; SIZE=32;;
        s390x) GOARCH=s390x; TRIPLE=s390x-linux-gnu; QEMU=qemu-s390x;;
        ppc64) GOARCH=ppc64; TRIPLE=powerpc64-linux-gnu; QEMU=qemu-ppc64;;
        ppc64le) GOARCH=ppc64le; TRIPLE=powerpc64le-linux-gnu;
                 QEMU=qemu-ppc64le;;
        mips) GOARCH=mips; TRIPLE=mips-linux-gnu; QEMU=qemu-mips; SIZE=32;
              CGO_LD="-no-pie";;
        riscv64) GOARCH=riscv64; TRIPLE=riscv64-linux-gnu;
                 QEMU=qemu-riscv64;;
        *) return 1;;
    esac
    if [ -n "$TRIPLE" ];
    then
        if [ "$SIZE" = 32 ];
        then
            CC="clang --target=$TRIPLE";
            EXTRA="-fforce-enable-int128";
            CGO_ALLOW="-fforce-enable-int128";
        else
            CC="$TRIPLE-gcc";
        fi
    fi
}

run() {
    if [ -n "$QEMU" ];
    then
        "$QEMU" -L "/usr/$TRIPLE" "$@";
    else
        "$@";
    fi
}

mark() {
    echo "$3" > "$OUT/result/$1-$2";
}

go_cell() {
    local a=$1 bin=$OUT/bin/kat-go-$1.test;
    if ! (cd "$ROOT" && GOOS=linux GOARCH=$GOARCH GOARM=$GOARM \
            CGO_ENABLED=1 CGO_CFLAGS_ALLOW="$CGO_ALLOW" \
            CGO_LDFLAGS="$CGO_LD" \
            CC="$CC $EXTRA" go test -c -o "$bin" ./src/kat) \
            > "$OUT/log/go-build-$a.txt" 2>&1;
    then
        mark go-check "$a" "no-build";
        mark go-gen "$a" "no-build";
        return;
    fi
    if (cd "$ROOT/src/kat" && run "$bin" -test.v -test.timeout 0) \
            > "$OUT/log/go-check-$a.txt" 2>&1;
    then
        mark go-check "$a" pass;
    else
        mark go-check "$a" FAIL;
    fi
    mkdir -p "$OUT/gen/go-$a";
    if (cd "$ROOT/src/kat" && HIGHCTIDH_KAT_OUT="$OUT/gen/go-$a" \
            run "$bin" -test.timeout 0) > "$OUT/log/go-gen-$a.txt" 2>&1;
    then
        mark go-gen "$a" done;
    else
        mark go-gen "$a" FAIL;
    fi
}

c_cell() {
    local a=$1 n bin ok=pass gen=done;
    mkdir -p "$OUT/gen/c-$a";
    for n in $SIZES;
    do
        bin=$OUT/bin/kat-c$n-$a;
        if ! (cd "$ROOT/src" && $CC $EXTRA -O2 -fwrapv -DGETRANDOM \
                -DHIGHCTIDH_PORTABLE=1 -DPLATFORM=$GOARCH \
                -DPLATFORM_SIZE=$SIZE -DBITS=$n \
                "-DNAMESPACEBITS(x)=highctidh_${n}_##x" \
                "-DNAMESPACEGENERIC(x)=highctidh_##x" -I. -o "$bin" kat.c \
                $CSRC fiat_p$n.c fp_inv$n.c fp_sqrt$n.c primes$n.c) \
                > "$OUT/log/c-build-$n-$a.txt" 2>&1;
        then
            mark c-check "$a" "no-build";
            mark c-gen "$a" "no-build";
            return;
        fi
        set -- $(sed -n 's/^seed_. = //p' "$ROOT/src/kat/ctidh$n.kat");
        if ! run "$bin" "$@" > "$OUT/gen/c-$a/ctidh$n.kat" \
                2> "$OUT/log/c-gen-$n-$a.txt";
        then
            gen=FAIL;
        fi
        if ! cmp -s "$OUT/gen/c-$a/ctidh$n.kat" "$ROOT/src/kat/ctidh$n.kat";
        then
            ok=FAIL;
        fi
    done
    mark c-check "$a" "$ok";
    mark c-gen "$a" "$gen";
}

py_cell() {
    local a=$1 image;
    if [ -z "${PY_IMAGE:-}" ];
    then
        mark py-check "$a" "no-image";
        mark py-gen "$a" "no-image";
        return;
    fi
    image=$(printf "$PY_IMAGE" "$a");
    if ! podman image exists "$image";
    then
        mark py-check "$a" "no-image";
        mark py-gen "$a" "no-image";
        return;
    fi
    mkdir -p "$OUT/gen/py-$a";
    podman run --rm --name "${CONTAINER_PREFIX:-kat-matrix-}$a-$$" \
        -v "$ROOT:/src:ro,Z" -v "$OUT/gen/py-$a:/out:Z" \
        -e CC=$([ "$SIZE" = 32 ] && echo clang || echo gcc) \
        -e CFLAGS="$EXTRA" -e HIGHCTIDH_PORTABLE=1 \
        "$image" $PYPRE sh -c '
            cp -r /src /w && cd /w/src &&
            pip install -q --break-system-packages --no-build-isolation . \
                > /tmp/build.txt 2>&1 || { cat /tmp/build.txt; exit 3; }
            python3 -m pytest -v -p no:cacheprovider -o addopts= \
                tests/test_kat.py;
            echo "check-exit $?";
            HIGHCTIDH_KAT_OUT=/out python3 tests/test_kat.py;
            echo "gen-exit $?"' > "$OUT/log/py-$a.txt" 2>&1;
    if grep -q '^check-exit 0$' "$OUT/log/py-$a.txt";
    then
        mark py-check "$a" pass;
    elif grep -q '^check-exit' "$OUT/log/py-$a.txt";
    then
        mark py-check "$a" FAIL;
    else
        mark py-check "$a" "no-build";
    fi
    if grep -q '^gen-exit 0$' "$OUT/log/py-$a.txt";
    then
        mark py-gen "$a" done;
    else
        mark py-gen "$a" FAIL;
    fi
}

cell() {
    local a=$1;
    if ! arch_env "$a";
    then
        for c in go-check go-gen c-check c-gen py-check py-gen;
        do
            mark "$c" "$a" "unknown";
        done
        return;
    fi
    if [ -n "$QEMU" ] && ! command -v "$QEMU" > /dev/null;
    then
        for c in go-check go-gen c-check c-gen;
        do
            mark "$c" "$a" "no-qemu";
        done
    elif ! command -v ${CC%% *} > /dev/null;
    then
        for c in go-check go-gen c-check c-gen;
        do
            mark "$c" "$a" "no-cc";
        done
    else
        go_cell "$a" &
        c_cell "$a" &
    fi
    py_cell "$a" &
    wait;
}

for a in $ARCHES;
do
    cell "$a" &
done
wait;

mkdir -p "$OUT/gen/committed";
cp "$ROOT"/src/kat/*.kat "$OUT/gen/committed/";

declare -A GROUP;
LABELS="";
group_of() {
    local d=$1 h;
    G="-";
    [ -f "$d/ctidh511.kat" ] || return 0;
    h=$(cd "$d" && cat ctidh511.kat ctidh512.kat ctidh1024.kat \
        ctidh2048.kat 2> /dev/null | sha256sum | cut -c1-12);
    if [ -z "${GROUP[$h]:-}" ];
    then
        LABELS="${LABELS}x";
        GROUP[$h]=$(echo ABCDEFGHIJKLMNOPQRSTUVWXYZ | cut -c${#LABELS});
    fi
    G=${GROUP[$h]};
}

group_of "$OUT/gen/committed";
canon=$G;

get() {
    cat "$OUT/result/$1-$2" 2> /dev/null || echo "-";
}

printf "\nKAT matrix (generated sets: %s = committed canonical set)\n\n" \
    "$canon";
printf "%-8s %-9s %-9s %-9s %-6s %-6s %-6s\n" \
    arch go-check c-check py-check go-gen c-gen py-gen;
FAILED=0;
for a in $ARCHES;
do
    row=();
    for c in go-check c-check py-check;
    do
        r=$(get "$c" "$a");
        [ "$r" = FAIL ] && FAILED=1;
        row+=("$r");
    done
    for l in go c py;
    do
        r=$(get "$l-gen" "$a");
        if [ "$r" = done ];
        then
            group_of "$OUT/gen/$l-$a";
            [ "$G" != "$canon" ] && FAILED=1;
            row+=("$G");
        else
            row+=("$r");
        fi
    done
    printf "%-8s %-9s %-9s %-9s %-6s %-6s %-6s\n" "$a" "${row[@]}";
done

printf "\nGenerated on (row) == generated on (column):\n\n%-11s" "";
SETS="committed";
for a in $ARCHES;
do
    for l in go c py;
    do
        [ "$(get "$l-gen" "$a")" = done ] && SETS="$SETS $l-$a";
    done
done
declare -A SG;
for s in $SETS;
do
    group_of "$OUT/gen/$s";
    SG[$s]=$G;
    printf "%s" "$G";
done
printf "\n";
for s in $SETS;
do
    printf "%-11s" "$s";
    for t in $SETS;
    do
        if [ "${SG[$s]}" = "${SG[$t]}" ];
        then
            printf "=";
        else
            printf "x";
        fi
    done
    printf "\n";
done

exit $FAILED;
