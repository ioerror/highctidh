package ctidh512

import (
	"bytes"
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
	p := gopointer.Save(rng)
	outsz1 := len(message)
	outbuf1 := make([]byte, outsz1)
	test_go_fillrandom(p, outbuf1)
	t.Logf("out: `%s`", outbuf1)
	require.Equal(t, message, outbuf1)

	message = []byte("how now brown cow holy cow")
	rng = bytes.NewReader(message)
	p = gopointer.Save(rng)
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
	p := gopointer.Save(iotest.OneByteReader(bytes.NewReader(message)))
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
