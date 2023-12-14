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
    echo -n "$GOARCH $BITS bits:";
    go test -v;
    echo -e "$CHECKMARK";

    cd ../;
done

echo "Cross compile for $GOOS with CGO_ENABLED=$CGO_ENABLED...";
for BITS in 511 512 1024 2048;
do
    cd ctidh$BITS;

    export GOARCH=amd64;
    echo -n "$GOARCH $BITS bits:";
    CC=x86_64-linux-gnu-gcc go build;
    echo -e "$CHECKMARK";

    export GOARCH=arm64;
    echo -n "$GOARCH $BITS bits:";
    CC=aarch64-linux-gnu-gcc go build;
    echo -e "$CHECKMARK";

    export GOARCH=ppc64le;
    echo -n "$GOARCH $BITS bits:";
    CC=powerpc64le-linux-gnu-gcc go build;
    echo -e "$CHECKMARK";

    export GOARCH=riscv64;
    echo -n "$GOARCH $BITS bits:";
    CC=riscv64-linux-gnu-gcc GOARCH=riscv64 go build;
    echo -e "$CHECKMARK";

    export GOARCH=s390x;
    echo -n "$GOARCH $BITS bits:";
    CC=s390x-linux-gnu-gcc go build;
    echo -e "$CHECKMARK";

    export GOARCH=mips64;
    echo -n "$GOARCH $BITS bits:";
    CC=mips64-linux-gnuabi64-gcc go build;
    echo -e "$CHECKMARK";

  # XXX TODO: fix this: error: '__int128' is not supported on this target
  # export CGO_CFLAGS_ALLOW="-fforce-enable-int128";
  # export GOARCH=386;
  # echo -n "$GOARCH $BITS bits:";
  # CC="i686-linux-gnu-gcc -fno-stack-check -fno-stack-protector" \
  # CGO_CFLAGS="-fno-stack-protector" \
  # go build;
  # echo -e "$CHECKMARK";
  
  # for SUBARCH in 5 6 7
  # do
  #     export GOARCH=arm;
  #     export GOARM=$SUBARCH;
  #     echo -n "$GOARCH/$GOARM $BITS bits:";
  #     CC="arm-linux-gnueabi-gcc -mfloat-abi=soft" \
  #     CGO_CFLAGS="-fno-stack-protector" \
  #         go build -x;
  #     echo -e "$CHECKMARK";
  # done
  # export GOARM="";

    cd ../;
done;
echo "Cross compile for $GOOS with CGO_ENABLED=$CGO_ENABLED successful.";
