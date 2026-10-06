package kat_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"codeberg.org/vula/highctidh/src/ctidh1024"
	"codeberg.org/vula/highctidh/src/ctidh2048"
	"codeberg.org/vula/highctidh/src/ctidh511"
	"codeberg.org/vula/highctidh/src/ctidh512"
)

type splitmix64 struct {
	state uint64
	buf   []byte
}

func (r *splitmix64) Read(p []byte) (int, error) {
	for len(r.buf) < len(p) {
		r.state += 0x9e3779b97f4a7c15
		z := r.state
		z = (z ^ (z >> 30)) * 0xbf58476d1ce4e5b9
		z = (z ^ (z >> 27)) * 0x94d049bb133111eb
		z ^= z >> 31
		var b [8]byte
		binary.LittleEndian.PutUint64(b[:], z)
		r.buf = append(r.buf, b[:]...)
	}
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

type scheme struct {
	gen    func(io.Reader) []byte
	pub    func([]byte) ([]byte, error)
	action func(sk, pk []byte) ([]byte, error)
}

func loadSK511(b []byte) (*ctidh511.PrivateKey, error) {
	k := new(ctidh511.PrivateKey)
	return k, k.FromBytes(b)
}

func loadPK511(b []byte) (*ctidh511.PublicKey, error) {
	k := new(ctidh511.PublicKey)
	return k, k.FromBytes(b)
}

func loadSK512(b []byte) (*ctidh512.PrivateKey, error) {
	k := new(ctidh512.PrivateKey)
	return k, k.FromBytes(b)
}

func loadPK512(b []byte) (*ctidh512.PublicKey, error) {
	k := new(ctidh512.PublicKey)
	return k, k.FromBytes(b)
}

func loadSK1024(b []byte) (*ctidh1024.PrivateKey, error) {
	k := new(ctidh1024.PrivateKey)
	return k, k.FromBytes(b)
}

func loadPK1024(b []byte) (*ctidh1024.PublicKey, error) {
	k := new(ctidh1024.PublicKey)
	return k, k.FromBytes(b)
}

func loadSK2048(b []byte) (*ctidh2048.PrivateKey, error) {
	k := new(ctidh2048.PrivateKey)
	return k, k.FromBytes(b)
}

func loadPK2048(b []byte) (*ctidh2048.PublicKey, error) {
	k := new(ctidh2048.PublicKey)
	return k, k.FromBytes(b)
}

var schemes = map[int]scheme{
	511: {
		func(r io.Reader) []byte { return ctidh511.GeneratePrivateKey(r).Bytes() },
		func(b []byte) ([]byte, error) {
			k, err := loadSK511(b)
			if err != nil {
				return nil, err
			}
			return ctidh511.DerivePublicKey(k).Bytes(), nil
		},
		func(s, p []byte) ([]byte, error) {
			k, err := loadSK511(s)
			if err != nil {
				return nil, err
			}
			pk, err := loadPK511(p)
			if err != nil {
				return nil, err
			}
			return ctidh511.DeriveSecret(k, pk), nil
		},
	},
	512: {
		func(r io.Reader) []byte { return ctidh512.GeneratePrivateKey(r).Bytes() },
		func(b []byte) ([]byte, error) {
			k, err := loadSK512(b)
			if err != nil {
				return nil, err
			}
			return ctidh512.DerivePublicKey(k).Bytes(), nil
		},
		func(s, p []byte) ([]byte, error) {
			k, err := loadSK512(s)
			if err != nil {
				return nil, err
			}
			pk, err := loadPK512(p)
			if err != nil {
				return nil, err
			}
			return ctidh512.DeriveSecret(k, pk), nil
		},
	},
	1024: {
		func(r io.Reader) []byte { return ctidh1024.GeneratePrivateKey(r).Bytes() },
		func(b []byte) ([]byte, error) {
			k, err := loadSK1024(b)
			if err != nil {
				return nil, err
			}
			return ctidh1024.DerivePublicKey(k).Bytes(), nil
		},
		func(s, p []byte) ([]byte, error) {
			k, err := loadSK1024(s)
			if err != nil {
				return nil, err
			}
			pk, err := loadPK1024(p)
			if err != nil {
				return nil, err
			}
			return ctidh1024.DeriveSecret(k, pk), nil
		},
	},
	2048: {
		func(r io.Reader) []byte { return ctidh2048.GeneratePrivateKey(r).Bytes() },
		func(b []byte) ([]byte, error) {
			k, err := loadSK2048(b)
			if err != nil {
				return nil, err
			}
			return ctidh2048.DerivePublicKey(k).Bytes(), nil
		},
		func(s, p []byte) ([]byte, error) {
			k, err := loadSK2048(s)
			if err != nil {
				return nil, err
			}
			pk, err := loadPK2048(p)
			if err != nil {
				return nil, err
			}
			return ctidh2048.DeriveSecret(k, pk), nil
		},
	},
}

var sizes = []int{511, 512, 1024, 2048}

var fields = []string{
	"size", "rng",
	"seed_a", "sk_a", "pk_a",
	"seed_b", "sk_b", "pk_b",
	"seed_f", "sk_f",
	"ss_a_b", "blind_f_a",
}

func readKAT(t *testing.T, size int) map[string]string {
	t.Helper()
	f, err := os.Open(fmt.Sprintf("ctidh%d.kat", size))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	e := map[string]string{}
	s := bufio.NewScanner(f)
	s.Buffer(nil, 1<<20)
	for s.Scan() {
		k, v, ok := strings.Cut(s.Text(), " = ")
		if !ok {
			t.Fatalf("bad line %q", s.Text())
		}
		e[k] = v
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	if e["size"] != strconv.Itoa(size) || e["rng"] != "splitmix64" {
		t.Fatalf("bad header: %q %q", e["size"], e["rng"])
	}
	return e
}

func unhex(s string) []byte {
	b, err := hex.DecodeString(s)
	if err != nil {
		return nil
	}
	return b
}

func result(b []byte, err error) string {
	if err != nil {
		return "error: " + err.Error()
	}
	return hex.EncodeToString(b)
}

func keygen(t *testing.T, sc scheme, e map[string]string, n string) string {
	t.Helper()
	s, err := strconv.ParseUint(e["seed_"+n], 10, 64)
	if err != nil {
		t.Fatal(err)
	}
	return hex.EncodeToString(sc.gen(&splitmix64{state: s}))
}

func generate(
	t *testing.T, sc scheme, e map[string]string,
) map[string]string {
	t.Helper()
	g := map[string]string{}
	for _, k := range []string{"size", "rng", "seed_a", "seed_b", "seed_f"} {
		g[k] = e[k]
	}
	for _, n := range []string{"a", "b", "f"} {
		g["sk_"+n] = keygen(t, sc, e, n)
	}
	g["pk_a"] = result(sc.pub(unhex(g["sk_a"])))
	g["pk_b"] = result(sc.pub(unhex(g["sk_b"])))
	g["ss_a_b"] = result(sc.action(unhex(g["sk_a"]), unhex(g["pk_b"])))
	g["blind_f_a"] = result(sc.action(unhex(g["sk_f"]), unhex(g["pk_a"])))
	return g
}

func check(t *testing.T, sc scheme, e map[string]string) {
	t.Helper()
	got := map[string]string{}
	for _, n := range []string{"a", "b", "f"} {
		got["sk_"+n] = keygen(t, sc, e, n)
	}
	got["pk_a"] = result(sc.pub(unhex(e["sk_a"])))
	got["pk_b"] = result(sc.pub(unhex(e["sk_b"])))
	got["ss_a_b"] = result(sc.action(unhex(e["sk_a"]), unhex(e["pk_b"])))
	got["blind_f_a"] = result(sc.action(unhex(e["sk_f"]), unhex(e["pk_a"])))
	for _, k := range fields {
		if v, ok := got[k]; ok && v != e[k] {
			t.Errorf("%s:\n got %s\nwant %s", k, v, e[k])
		}
	}
}

func TestKAT(t *testing.T) {
	dir := os.Getenv("HIGHCTIDH_KAT_OUT")
	for _, size := range sizes {
		size := size
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			sc := schemes[size]
			e := readKAT(t, size)
			if dir == "" {
				check(t, sc, e)
				return
			}
			g := generate(t, sc, e)
			var buf bytes.Buffer
			for _, k := range fields {
				fmt.Fprintf(&buf, "%s = %s\n", k, g[k])
			}
			name := filepath.Join(dir, fmt.Sprintf("ctidh%d.kat", size))
			if err := os.WriteFile(name, buf.Bytes(), 0o644); err != nil {
				t.Fatal(err)
			}
		})
	}
}
