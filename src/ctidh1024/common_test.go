package ctidh1024

import (
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"os/exec"
	"runtime"
	"strings"
	"testing"
	"testing/iotest"

	gopointer "github.com/mattn/go-pointer"
	"github.com/stretchr/testify/require"
)

func TestFillRandom(t *testing.T) {
	message := []byte("AAAA")
	rng := bytes.NewReader(message)
	p := gopointer.Save(&rngContext{rng: rng})
	outsz1 := len(message)
	outbuf1 := make([]byte, outsz1)
	test_go_fillrandom(p, outbuf1)
	t.Logf("out: `%s`", outbuf1)
	require.Equal(t, message, outbuf1)

	message = []byte("how now brown cow holy cow")
	rng = bytes.NewReader(message)
	p = gopointer.Save(&rngContext{rng: rng})
	outsz2 := len(message) - len(" holy cow")
	outbuf2 := make([]byte, outsz2)
	test_go_fillrandom(p, outbuf2)
	t.Logf("out: `%s`", outbuf2)
	require.Equal(t, message[:outsz2], outbuf2)

	outsz3 := len(message) - outsz2
	outbuf3 := make([]byte, outsz3)
	test_go_fillrandom(p, outbuf3)
	t.Logf("out: `%s`", outbuf3)
	require.Equal(t, message[outsz2:], outbuf3)
}

func TestFillRandomShortRead(t *testing.T) {
	message := []byte("how now brown cow holy cow")
	p := gopointer.Save(&rngContext{rng: iotest.OneByteReader(bytes.NewReader(message))})
	defer gopointer.Unref(p)
	outbuf := make([]byte, len(message))
	require.NotPanics(t, func() { test_go_fillrandom(p, outbuf) })
	require.Equal(t, message, outbuf)
}

func TestFillRandomCheckptr(t *testing.T) {
	if os.Getenv("HIGHCTIDH_CHECKPTR_CHILD") != "" {
		t.Skip("running as the checkptr child")
	}
	if testing.Short() {
		t.Skip("short mode")
	}
	if strings.Contains(os.Getenv("CGO_CFLAGS"), "-fsanitize") {
		t.Skip("sanitizer CGO_CFLAGS")
	}
	gobin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go not in PATH")
	}
	host, err := exec.Command(gobin, "env", "GOHOSTOS", "GOHOSTARCH").Output()
	if err != nil || strings.Join(strings.Fields(string(host)), "/") != runtime.GOOS+"/"+runtime.GOARCH {
		t.Skip("not a native build")
	}
	for _, args := range [][]string{
		{"vet", "."},
		{"test", "-count=1", "-gcflags=all=-d=checkptr", "-run", "^TestFillRandom", "."},
	} {
		cmd := exec.Command(gobin, args...)
		cmd.Env = append(os.Environ(), "HIGHCTIDH_CHECKPTR_CHILD=1")
		out, err := cmd.CombinedOutput()
		if err != nil {
			t.Fatalf("go %s: %v\n%s", strings.Join(args, " "), err, out)
		}
	}
}

var errTestRNG = errors.New("test rng failure")

type testRNG struct {
	ctr   uint64
	buf   []byte
	n     int
	limit int
	panic interface{}
}

func (r *testRNG) Read(p []byte) (int, error) {
	if r.limit >= 0 && r.n+len(p) > r.limit {
		if r.panic != nil {
			panic(r.panic)
		}
		k := copy(p, r.next(r.limit-r.n))
		r.n += k
		return k, errTestRNG
	}
	k := copy(p, r.next(len(p)))
	r.n += k
	return k, nil
}

func (r *testRNG) next(n int) []byte {
	for len(r.buf) < n {
		var c [8]byte
		binary.LittleEndian.PutUint64(c[:], r.ctr)
		r.ctr++
		h := sha256.Sum256(c[:])
		r.buf = append(r.buf, h[:]...)
	}
	out := r.buf[:n]
	r.buf = r.buf[n:]
	return out
}

type countingReader struct {
	r     io.Reader
	calls int
}

func (c *countingReader) Read(p []byte) (int, error) {
	c.calls++
	return c.r.Read(p)
}

func TestFillRandomErrorFillsPattern(t *testing.T) {
	rng := &countingReader{r: &testRNG{limit: 3}}
	ctx := &rngContext{rng: rng}
	p := gopointer.Save(ctx)
	defer gopointer.Unref(p)
	out := make([]byte, 40)
	require.NotPanics(t, func() { test_go_fillrandom(p, out) })
	require.ErrorIs(t, ctx.err, errTestRNG)
	for i := range out {
		require.Equal(t, byte(i/4), out[i])
	}
	calls := rng.calls
	out = make([]byte, 8)
	test_go_fillrandom(p, out)
	require.Equal(t, calls, rng.calls)
	require.Equal(t, []byte{0, 0, 0, 0, 1, 1, 1, 1}, out)
}

func TestGeneratePrivateKeyCheckedReaderError(t *testing.T) {
	full := &testRNG{limit: -1}
	want, err := GeneratePrivateKeyChecked(full)
	require.NoError(t, err)
	require.Equal(t, GeneratePrivateKey(&testRNG{limit: -1}).Bytes(), want.Bytes())
	require.NoError(t, new(PrivateKey).FromBytes(want.Bytes()))
	total := full.n
	require.Greater(t, total, 0)

	for n := 0; n < total; n++ {
		require.NotPanics(t, func() {
			k, err := GeneratePrivateKeyChecked(&testRNG{limit: n})
			require.ErrorIs(t, err, errTestRNG, "limit %d", n)
			require.Nil(t, k)
		})
	}
	k, err := GeneratePrivateKeyChecked(&testRNG{limit: total})
	require.NoError(t, err)
	require.Equal(t, want.Bytes(), k.Bytes())
}

func TestGenerateKeyPairWithRNGChecked(t *testing.T) {
	priv, pub, err := GenerateKeyPairWithRNGChecked(&testRNG{limit: -1})
	require.NoError(t, err)
	priv2, pub2 := GenerateKeyPairWithRNG(&testRNG{limit: -1})
	require.Equal(t, priv2.Bytes(), priv.Bytes())
	require.Equal(t, pub2.Bytes(), pub.Bytes())
	require.Equal(t, DerivePublicKey(priv).Bytes(), pub.Bytes())

	require.NotPanics(t, func() {
		priv, pub, err := GenerateKeyPairWithRNGChecked(&testRNG{limit: 100})
		require.ErrorIs(t, err, errTestRNG)
		require.Nil(t, priv)
		require.Nil(t, pub)
	})
}

func TestGenerateRNGPanics(t *testing.T) {
	value := errors.New("reader panic value")
	require.NotPanics(t, func() {
		k, err := GeneratePrivateKeyChecked(&testRNG{limit: 10, panic: value})
		require.Error(t, err)
		require.Nil(t, k)
		priv, pub, err := GenerateKeyPairWithRNGChecked(&testRNG{limit: 10, panic: value})
		require.Error(t, err)
		require.Nil(t, priv)
		require.Nil(t, pub)
		k, err = GeneratePrivateKeyChecked(nil)
		require.Error(t, err)
		require.Nil(t, k)
	})
	require.PanicsWithValue(t, value, func() { GeneratePrivateKey(&testRNG{limit: 10, panic: value}) })
	require.PanicsWithValue(t, value, func() { GenerateKeyPairWithRNG(&testRNG{limit: 10, panic: value}) })
	require.PanicsWithValue(t, errTestRNG, func() { GeneratePrivateKey(&testRNG{limit: 10}) })
	require.PanicsWithValue(t, errTestRNG, func() { GenerateKeyPairWithRNG(&testRNG{limit: 10}) })
}
