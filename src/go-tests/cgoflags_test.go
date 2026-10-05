package tests

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
)

var sizes = []string{"511", "512", "1024", "2048"}

func cgoCFLAGS(t *testing.T, goos, goarch, bits string) []string {
	t.Helper()
	gobin, err := exec.LookPath("go")
	if err != nil {
		t.Skip("go not in PATH")
	}
	cmd := exec.Command(gobin, "list", "-f", "{{join .CgoCFLAGS \" \"}}", ".")
	cmd.Dir = filepath.Join("..", "ctidh"+bits)
	cmd.Env = append(os.Environ(), "GOOS="+goos, "GOARCH="+goarch, "CGO_ENABLED=1", "GOFLAGS=")
	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("go list %s/%s ctidh%s: %v\n%s", goos, goarch, bits, err, out)
	}
	return strings.Fields(string(out))
}

func TestCgoNoCPUTuning(t *testing.T) {
	for _, bits := range sizes {
		src, err := os.ReadFile(filepath.Join("..", "ctidh"+bits, "common.go"))
		if err != nil {
			t.Fatal(err)
		}
		for _, f := range []string{"-march=", "-mtune=", "-mcpu="} {
			if strings.Contains(string(src), f) {
				t.Errorf("ctidh%s/common.go contains %s", bits, f)
			}
		}
		for _, p := range [][2]string{{"linux", "loong64"}, {"linux", "amd64"}, {"darwin", "amd64"}} {
			for _, f := range cgoCFLAGS(t, p[0], p[1], bits) {
				if strings.HasPrefix(f, "-march") || strings.HasPrefix(f, "-mtune") || strings.HasPrefix(f, "-mcpu") {
					t.Errorf("%s/%s ctidh%s: %s", p[0], p[1], bits, f)
				}
			}
		}
	}
}
