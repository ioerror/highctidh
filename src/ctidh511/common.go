package ctidh511

/*
 #define BITS 511
 #include "binding511.h"
*/
import "C"
import (
	"fmt"
	"io"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

// This function wraps go_fillrandom, so we can emulate the calls from the
// C library and test the results
func test_go_fillrandom(context unsafe.Pointer, outptr []byte) {
	highctidh_511_go_fillrandom(context, unsafe.Pointer(&outptr[0]), C.size_t(len(outptr)))
}

// This is called from the C library, DO NOT CHANGE THE FUNCTION INTERFACE
//
//export highctidh_511_go_fillrandom
func highctidh_511_go_fillrandom(context unsafe.Pointer, outptr unsafe.Pointer, outsz C.size_t) {
	rng := gopointer.Restore(context).(io.Reader)
	buf := make([]byte, outsz)
	count, err := rng.Read(buf)
	if err != nil {
		panic(err)
	}
	if count != int(outsz) {
		panic("rng fail")
	}
	for i := 0; i < int(outsz); i++ {
		p := unsafe.Pointer(uintptr(outptr) + uintptr(i))
		*(*uint8)(p) = uint8(buf[i])
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
