package kat_test

import (
	"bufio"
	"encoding/hex"
	"fmt"
	"os"
	"strconv"
	"strings"
	"testing"
)

type csidhCase struct {
	line int
	kind string
	e, f []byte
	a    string
}

func exponents(s string) ([]byte, error) {
	var r []byte
	for _, w := range strings.Split(s, ",") {
		v, err := strconv.Atoi(w)
		if err != nil || v < -128 || v > 127 {
			return nil, fmt.Errorf("bad exponent %q", w)
		}
		r = append(r, byte(int8(v)))
	}
	return r, nil
}

func readCSIDH(t *testing.T, size int) []csidhCase {
	t.Helper()
	f, err := os.Open(fmt.Sprintf("csidh-reference-%d.txt", size))
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	var cases []csidhCase
	s := bufio.NewScanner(f)
	s.Buffer(nil, 1<<20)
	for n := 1; s.Scan(); n++ {
		w := strings.Fields(s.Text())
		c := csidhCase{line: n, kind: w[0]}
		for _, kv := range w[1:] {
			k, v, ok := strings.Cut(kv, "=")
			if !ok {
				t.Fatalf("line %d: bad field %q", n, kv)
			}
			switch k {
			case "e":
				c.e, err = exponents(v)
			case "f":
				c.f, err = exponents(v)
			case "A":
				c.a = v
			default:
				t.Fatalf("line %d: bad field %q", n, kv)
			}
			if err != nil {
				t.Fatalf("line %d: %v", n, err)
			}
		}
		if c.e == nil || c.a == "" || (c.kind == "ss") != (c.f != nil) {
			t.Fatalf("line %d: bad case", n)
		}
		cases = append(cases, c)
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	return cases
}

func TestCSIDHReference(t *testing.T) {
	for _, size := range []int{512, 1024} {
		size := size
		t.Run(strconv.Itoa(size), func(t *testing.T) {
			sc := schemes[size]
			for _, c := range readCSIDH(t, size) {
				pk, err := sc.pub(c.e)
				if err == nil && c.kind == "ss" {
					pk, err = sc.action(c.f, pk)
				}
				if err != nil {
					t.Errorf("line %d: %v", c.line, err)
					continue
				}
				if got := hex.EncodeToString(pk); got != c.a {
					t.Errorf("line %d %s:\n got %s\nwant %s",
						c.line, c.kind, got, c.a)
				}
			}
		})
	}
}
