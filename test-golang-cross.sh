#!/bin/bash
#
# Test cross compile of golang module using gcc cross compilers
#
set -e;

export GOOS=linux;
export CGO_ENABLED=1;
CHECKMARK="\xE2\x9C\x94"

echo "Cross compile for $GOOS with CGO_ENABLED=$CGO_ENABLED...";
for BITS in 511 512 1024 2048;
do
    cd ctidh$BITS;

    export GOARCH=amd64;
    echo -n "$GOARCH $BITS bits:";
    CC=x86_64-linux-gnu-gcc go build
    echo -e "$CHECKMARK";

    export GOARCH=arm64;
    echo -n "$GOARCH $BITS bits:";
    CC=aarch64-linux-gnu-gcc go build
    echo -e "$CHECKMARK";

    export GOARCH=ppc64le;
    echo -n "$GOARCH $BITS bits:";
    CC=powerpc64le-linux-gnu-gcc go build
    echo -e "$CHECKMARK";

    export GOARCH=riscv64;
    echo -n "$GOARCH $BITS bits:";
    CC=riscv64-linux-gnu-gcc GOARCH=riscv64 go build
    echo -e "$CHECKMARK";

    export GOARCH=s390x;
    echo -n "$GOARCH $BITS bits:";
    CC=s390x-linux-gnu-gcc go build
    echo -e "$CHECKMARK";

    export GOARCH=mips64;
    echo -n "$GOARCH $BITS bits:";
    CC=mips64-linux-gnuabi64-gcc go build
    echo -e "$CHECKMARK";

    cd ../;
done
echo "Cross compile for $GOOS with CGO_ENABLED=$CGO_ENABLED successful.";
