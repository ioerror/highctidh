package ctidh512

/*
 #cgo CFLAGS: -DBITS=512 -DGETRANDOM -DCGONUTS
 #cgo aarch64 CFLAGS: -DPLATFORM=aarch64 -DPLATFORM_SIZE=64 -march=native -mtune=native -DHIGHCTIDH_PORTABLE
 #cgo armv7l CFLAGS: -DPLATFORM=armv7l -DPLATFORM_SIZE=32 -fforce-enable-int128 -D__ARM32__ -DHIGHCTIDH_PORTABLE
 #cgo loongarch64 CFLAGS: -DPLATFORM=loongarch64 -DPLATFORM_SIZE=64 -march=native -mtune=native -DHIGHCTIDH_PORTABLE
 #cgo mips64 CFLAGS: -DPLATFORM=mips64 -DPLATFORM_SIZE=64 -fforce-enable-int128 -DHIGHCTIDH_PORTABLE
 #cgo ppc64le CFLAGS: -DPLATFORM=ppc64le -DPLATFORM_SIZE=64 -mtune=native -DHIGHCTIDH_PORTABLE
 #cgo riscv64 CFLAGS: -DPLATFORM=riscv64 -DPLATFORM_SIZE=64 -DHIGHCTIDH_PORTABLE
 #cgo s390x CFLAGS: -DPLATFORM=s390x -DPLATFORM_SIZE=64 -march=native -mtune=native -DHIGHCTIDH_PORTABLE
 #cgo sparc64 CFLAGS: -DPLATFORM=sparc64 -DPLATFORM_SIZE=64 -march=native -mtune=native -DHIGHCTIDH_PORTABLE
 #cgo amd64 386 CFLAGS: -DPLATFORM=x86_64 -DPLATFORM_SIZE=64 -march=native -mtune=native
 #cgo 386 CFLAGS: -DPLATFORM=i386 -DPLATFORM_SIZE=32 -fforce-enable-int128 -D__i386__ -DHIGHCTIDH_PORTABLE
 #include <stdlib.h>
 #include <stdint.h>
 #include "binding512.h"
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
