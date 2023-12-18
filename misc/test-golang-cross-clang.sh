#!/bin/bash
#
# Test cross compile of golang module using gcc cross compilers
#
set -e;

export GOOS=linux;
export CGO_ENABLED=1;
export HOST_ARCH=`uname -m`;
CHECKMARK="\xE2\x9C\x94";

echo "Running tests on $HOST_ARCH";
for BITS in 511 512 1024 2048;
do
    cd ctidh$BITS;

    export GOARCH=amd64;
    echo "$GOARCH $BITS bits:";
    CC=clang go test -v;
    echo -n "$GOARCH $BITS bits:";
    echo -e "$CHECKMARK";

    cd ../;
done

echo "Cross compile for $GOOS with CGO_ENABLED=$CGO_ENABLED...";
for BITS in 511 512 1024 2048;
do
    cd ctidh$BITS;

    export GOARCH=amd64;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=amd64 PLATFORM_SIZE=64 \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export GOARCH=arm64;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=arm64 PLATFORM_SIZE=64 \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export GOARCH=ppc64le;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=ppc64le PLATFORM_SIZE=64 \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export GOARCH=ppc64;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=ppc64le PLATFORM_SIZE=64 \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export GOARCH=riscv64;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=riscv64 PLATFORM_SIZE=64 \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export GOARCH=s390x;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=s390x PLATFORM_SIZE=64 \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export GOARCH=mips64;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=mips64 PLATFORM_SIZE=64 \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export GOARCH=mips64le;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=mips64el PLATFORM_SIZE=64 \
    LD=/usr/bin/$PLATFORM-linux-gnuabi$PLATFORM_SIZE-ld \
    CC="clang --target=$PLATFORM-pc-linux-gnu" \
    go build;
    echo -e "$CHECKMARK";

    export CGO_CFLAGS_ALLOW="-fforce-enable-int128";
    export GOARCH=mipsle;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=mipsel PLATFORM_SIZE=32 \
    CC="clang --target=$PLATFORM-pc-linux-gnu -fforce-enable-int128" \
    go build;
    echo -e "$CHECKMARK";

    export CGO_CFLAGS_ALLOW="-fforce-enable-int128";
    export GOARCH=mips;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=mips PLATFORM_SIZE=32 \
    CC="clang --target=$PLATFORM-pc-linux-gnu -fforce-enable-int128" \
    go build;
    echo -e "$CHECKMARK";

    export CGO_CFLAGS_ALLOW="-fforce-enable-int128";
    export GOARCH=386;
    echo -n "$GOARCH $BITS bits:";
    PLATFORM=i386 PLATFORM_SIZE=32 \
    CC="clang --target=$PLATFORM-pc-linux-gnu -fforce-enable-int128" \
    go build;
    echo -e "$CHECKMARK";

    for SUBARCH in 5 6 7
    do
       export CGO_CFLAGS_ALLOW="-fforce-enable-int128";
       export GOARCH=arm;
       export GOARM=$SUBARCH;
       echo -n "$GOARCH/$GOARM $BITS bits:";
       CC="clang --target=arm-pc-linux-gnu -fforce-enable-int128 -mfloat-abi=soft" \
           go build;
       echo -e "$CHECKMARK";
    done

    cd ../;
done;
echo "Cross compile for $GOOS with CGO_ENABLED=$CGO_ENABLED successful.";
