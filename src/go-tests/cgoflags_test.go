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

func TestCgoPlatformLines(t *testing.T) {
	want := map[[2]string][]string{
		{"darwin", "amd64"}:  {"-D__Darwin__", "-DGETRANDOM", "-D__x86_64__", "-DHIGHCTIDH_PORTABLE=1"},
		{"darwin", "arm64"}:  {"-D__ARM64__", "-D__Darwin__", "-DGETRANDOM", "-DHIGHCTIDH_PORTABLE=1"},
		{"windows", "amd64"}: {"-D__Windows__", "-DCGONUTS", "-DPLATFORM_SIZE=64", "-DHIGHCTIDH_PORTABLE=1"},
		{"windows", "arm64"}: {"-D__Windows__", "-DPLATFORM_SIZE=64", "-DHIGHCTIDH_PORTABLE=1"},
		{"solaris", "amd64"}: {"-m64", "-Wno-attributes", "-D__sun", "-D__i86pc__", "-DHIGHCTIDH_PORTABLE=1"},
		{"illumos", "amd64"}: {"-m64", "-Wno-attributes", "-D__sun", "-D__i86pc__", "-DHIGHCTIDH_PORTABLE=1"},
	}
	for _, bits := range sizes {
		src, err := os.ReadFile(filepath.Join("..", "ctidh"+bits, "common.go"))
		if err != nil {
			t.Fatal(err)
		}
		for _, line := range strings.Split(string(src), "\n") {
			f := strings.Fields(line)
			if len(f) > 1 && f[0] == "#cgo" && strings.Contains(f[1], "/") {
				t.Errorf("ctidh%s: constraint with a slash is ignored: %s", bits, line)
			}
		}
		for p, defs := range want {
			flags := cgoCFLAGS(t, p[0], p[1], bits)
			seen := map[string]string{}
			for _, f := range flags {
				if k, v, ok := strings.Cut(f, "="); ok && strings.HasPrefix(k, "-D") {
					if old, dup := seen[k]; dup && old != v {
						t.Errorf("%s/%s ctidh%s: %s=%s and %s=%s", p[0], p[1], bits, k, old, k, v)
					}
					seen[k] = v
				}
			}
			for _, d := range defs {
				found := false
				for _, f := range flags {
					found = found || f == d
				}
				if !found {
					t.Errorf("%s/%s ctidh%s: missing %s in %v", p[0], p[1], bits, d, flags)
				}
			}
		}
	}
}
