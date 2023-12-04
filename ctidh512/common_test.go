package ctidh512

import (
	"bytes"
	"testing"

	gopointer "github.com/mattn/go-pointer"
	"github.com/stretchr/testify/require"
)

func TestFillRandom(t *testing.T) {
	message := []byte("AAAA")
	rng := bytes.NewReader(message)
	p := gopointer.Save(rng)
	outsz := len(message)
	out := test_c_buf(outsz)
	test_go_fillrandom(p, out, outsz)
	t.Logf("out: `%s`", test_GoString(out, outsz))
	require.Equal(t, message, []byte(test_GoString(out, outsz)))

	message = []byte("how now brown cow")
	rng = bytes.NewReader(message)
	p = gopointer.Save(rng)
	outsz = len(message)
	out = test_c_buf(outsz)
	test_go_fillrandom(p, out, outsz)
	t.Logf("out: `%s`", test_GoString(out, outsz))
	require.Equal(t, message, []byte(test_GoString(out, outsz)))
}
