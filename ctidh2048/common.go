package ctidh2048

/*

 #cgo CFLAGS: -DBITS=2048 -DGETRANDOM -DCGONUTS -Wformat -Werror=format-security -D_FORTIFY_SOURCE=2 -fstack-protector-all -fpie -fPIC -Wextra -O3 -Os
 #cgo LDFLAGS: -Wl,-z,noexecstack -Wl,-z,relro

 // The following should work as native builds and as cross compiled builds.
 // Example cross compile build lines are provided as examples.

 // CC=aarch64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=arm64 go build -v
 #cgo arm64 CFLAGS: -DPLATFORM=aarch64 -DPLATFORM_SIZE=64 -DHIGHCTIDH_PORTABLE

 // CC=gcc CGO_ENABLED=1 GOOS=linux GOARCH=amd64 go build -v
 #cgo amd64 CFLAGS: -DPLATFORM=x86_64 -DPLATFORM_SIZE=64 -march=native -mtune=native

 // CC=powerpc64le-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=ppc64le go build -v
 #cgo ppc64le CFLAGS: -DPLATFORM=ppc64le -DPLATFORM_SIZE=64 -DHIGHCTIDH_PORTABLE

 // CC=riscv64-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=riscv64 go build -v
 #cgo riscv64 CFLAGS: -DPLATFORM=riscv64 -DPLATFORM_SIZE=64 -DHIGHCTIDH_PORTABLE

 // CC=s390x-linux-gnu-gcc CGO_ENABLED=1 GOOS=linux GOARCH=s390x go build -v
 #cgo s390x CFLAGS: -DPLATFORM=s390x -DPLATFORM_SIZE=64 -DHIGHCTIDH_PORTABLE

 // CC=mips64-linux-gnuabi64-gcc CGO_ENABLED=1 GOOS=linux GOARCH=mips64  go build
 // With clang, -fforce-enable-int128 must be added to the CFLAGS
 #cgo mips64 CFLAGS: -DPLATFORM=mips64 -DPLATFORM_SIZE=64 -DHIGHCTIDH_PORTABLE

 // The following should work as native builds with clang:

 #cgo arm CFLAGS: -DPLATFORM=armv7l -DPLATFORM_SIZE=32 -fforce-enable-int128 -D__ARM32__ -DHIGHCTIDH_PORTABLE
 #cgo loong64 CFLAGS: -DPLATFORM=loongarch64 -DPLATFORM_SIZE=64 -march=native -mtune=native -DHIGHCTIDH_PORTABLE
 #cgo sparc64 CFLAGS: -DPLATFORM=sparc64 -DPLATFORM_SIZE=64 -march=native -mtune=native -DHIGHCTIDH_PORTABLE

 // The following likely will not work as a cgo extension at this time:
 //  CC=clang CGO_ENABLED=1 GOOS=linux GOARCH=386  go build
 //  Results in: "common.go:57:18: cannot use _Ctype_ulong(size) (value of type _Ctype_ulong) as _Ctype_uint value in argument to (_Cfunc__CMalloc)"
 #cgo 386 CFLAGS: -DPLATFORM=i386 -DPLATFORM_SIZE=32 -fforce-enable-int128 -D__i386__ -DHIGHCTIDH_PORTABLE

 #include <stdlib.h>
 #include <stdint.h>
 #include "binding2048.h"
*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

func test_go_fillrandom(context unsafe.Pointer, outptr unsafe.Pointer, outsz int) {
	go_fillrandom(context, outptr, C.size_t(outsz))
}

func test_c_buf(size int) unsafe.Pointer {
	return C.malloc(C.ulong(size))
}

func test_free(p unsafe.Pointer) {
	C.free(p)
}

func test_GoString(x unsafe.Pointer, size int) string {
	ret := C.GoBytes(x, C.int(size))
	return string(ret)
}

//export go_fillrandom
func go_fillrandom(context unsafe.Pointer, outptr unsafe.Pointer, outsz C.size_t) {
	rng := gopointer.Restore(context).(io.Reader)
	buf := make([]byte, outsz)
	count, err := rng.Read(buf)
	if err != nil {
		panic(err)
	}
	if count != int(outsz) {
		panic("rng fail")
	}
	p := uintptr(outptr)
	for i := 0; i < int(outsz); {
		(*(*uint8)(unsafe.Pointer(p))) = uint8(buf[i])
		p += 1
		i += 1
	}
}

// Name returns the string naming of the current
// CTIDH that this binding is being used with;
// Valid values are:
//
// CTIDH-511, CTIDH-512, CTIDH-1024 and, CTIDH-2048.
func Name() string {
	return fmt.Sprintf("CTIDH-%d", C.BITS)
}

func validateBitSize(bits int) {
	switch bits {
	case 511:
	case 512:
	case 1024:
	case 2048:
	default:
		panic("CTIDH/cgo: BITS must be 511 or 512 or 1024 or 2048")
	}
}
