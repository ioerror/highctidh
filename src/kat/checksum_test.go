package kat_test

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"math/bits"
	"os"
	"strconv"
	"strings"
	"testing"

	"codeberg.org/vula/highctidh/src/ctidh1024"
	"codeberg.org/vula/highctidh/src/ctidh2048"
	"codeberg.org/vula/highctidh/src/ctidh511"
	"codeberg.org/vula/highctidh/src/ctidh512"
)

const (
	loopsSmall = 64
	loopsBig   = 512
	guardLen   = 16
)

func chachaBlock(out []byte, key []byte, counter uint64) {
	var s, x [16]uint32
	s[0], s[1], s[2], s[3] = 0x61707865, 0x3320646e, 0x79622d32, 0x6b206574
	for i := 0; i < 8; i++ {
		s[4+i] = binary.LittleEndian.Uint32(key[4*i:])
	}
	s[12] = uint32(counter)
	s[13] = uint32(counter >> 32)
	x = s
	qr := func(a, b, c, d int) {
		x[a] += x[b]
		x[d] = bits.RotateLeft32(x[d]^x[a], 16)
		x[c] += x[d]
		x[b] = bits.RotateLeft32(x[b]^x[c], 12)
		x[a] += x[b]
		x[d] = bits.RotateLeft32(x[d]^x[a], 8)
		x[c] += x[d]
		x[b] = bits.RotateLeft32(x[b]^x[c], 7)
	}
	for i := 0; i < 10; i++ {
		qr(0, 4, 8, 12)
		qr(1, 5, 9, 13)
		qr(2, 6, 10, 14)
		qr(3, 7, 11, 15)
		qr(0, 5, 10, 15)
		qr(1, 6, 11, 12)
		qr(2, 7, 8, 13)
		qr(3, 4, 9, 14)
	}
	for i := range x {
		binary.LittleEndian.PutUint32(out[4*i:], x[i]+s[i])
	}
}

type knownRandom struct {
	key  [32]byte
	pool []byte
}

func (r *knownRandom) Read(p []byte) (int, error) {
	for i := range p {
		if len(r.pool) == 0 {
			stream := make([]byte, 768)
			for j := 0; j < 12; j++ {
				chachaBlock(stream[64*j:], r.key[:], uint64(j))
			}
			copy(r.key[:], stream[:32])
			r.pool = stream[32:]
		}
		p[i] = r.pool[0]
		r.pool = r.pool[1:]
	}
	return len(p), nil
}

func salsaCore(out []byte, in []byte, key []byte) {
	var s, x [16]uint32
	s[0], s[5], s[10], s[15] = 0x61707865, 0x3320646e, 0x79622d32, 0x6b206574
	for i := 0; i < 4; i++ {
		s[1+i] = binary.LittleEndian.Uint32(key[4*i:])
		s[6+i] = binary.LittleEndian.Uint32(in[4*i:])
		s[11+i] = binary.LittleEndian.Uint32(key[16+4*i:])
	}
	x = s
	qr := func(a, b, c, d int) {
		x[b] ^= bits.RotateLeft32(x[a]+x[d], 7)
		x[c] ^= bits.RotateLeft32(x[b]+x[a], 9)
		x[d] ^= bits.RotateLeft32(x[c]+x[b], 13)
		x[a] ^= bits.RotateLeft32(x[d]+x[c], 18)
	}
	for i := 0; i < 10; i++ {
		qr(0, 4, 8, 12)
		qr(5, 9, 13, 1)
		qr(10, 14, 2, 6)
		qr(15, 3, 7, 11)
		qr(0, 1, 2, 3)
		qr(5, 6, 7, 4)
		qr(10, 11, 8, 9)
		qr(15, 12, 13, 14)
	}
	for i := range x {
		binary.LittleEndian.PutUint32(out[4*i:], x[i]+s[i])
	}
}

type checksum struct {
	state [64]byte
}

func (c *checksum) add(x []byte) {
	var next [64]byte
	for ; len(x) >= 16; x = x[16:] {
		salsaCore(next[:], x[:16], c.state[:32])
		c.state = next
	}
	var last [16]byte
	copy(last[:], x)
	last[len(x)] = 1
	c.state[0] ^= 1
	salsaCore(next[:], last[:], c.state[:32])
	c.state = next
}

func (c *checksum) hex() string {
	return hex.EncodeToString(c.state[:32])
}

type dhScheme struct {
	skLen, pkLen, ssLen int
	keypair             func(pk, sk []byte, rng io.Reader) error
	dh                  func(ss, pk, sk []byte) error
}

type guarded struct {
	mem []byte
	n   int
}

func newGuarded(n int) *guarded {
	return &guarded{mem: make([]byte, guardLen+n+guardLen), n: n}
}

func (g *guarded) at(l int) []byte {
	return g.mem[guardLen : guardLen+l]
}

type pattern struct {
	n uint64
}

func (p *pattern) fill(b []byte) {
	for i := range b {
		p.n++
		b[i] = byte(p.n*0x9e3779b97f4a7c15>>56) ^ 0xa5
	}
}

func (p *pattern) frame(g *guarded, l int) []byte {
	p.fill(g.mem[:guardLen])
	p.fill(g.mem[guardLen+l : guardLen+l+guardLen])
	return append([]byte(nil), g.mem[:guardLen+l+guardLen]...)
}

func sameFrame(g *guarded, l int, want []byte) bool {
	return bytes.Equal(g.mem[:guardLen+l+guardLen], want)
}

func sameGuards(g *guarded, l int, want []byte) bool {
	return bytes.Equal(g.mem[:guardLen], want[:guardLen]) &&
		bytes.Equal(g.mem[guardLen+l:guardLen+l+guardLen],
			want[guardLen+l:])
}

func tryDH(s dhScheme, loops int) (string, error) {
	n := s.skLen
	if s.pkLen > n {
		n = s.pkLen
	}
	if s.ssLen > n {
		n = s.ssLen
	}
	a, b, c, d, e, f := newGuarded(n), newGuarded(n), newGuarded(n),
		newGuarded(n), newGuarded(n), newGuarded(n)
	a2, c2, d2, e2 := newGuarded(n), newGuarded(n), newGuarded(n),
		newGuarded(n)
	rng := new(knownRandom)
	var pat pattern
	var sum checksum

	keypair := func(pk, sk *guarded) error {
		wpk := pat.frame(pk, s.pkLen)
		wsk := pat.frame(sk, s.skLen)
		if err := s.keypair(pk.at(s.pkLen), sk.at(s.skLen), rng); err != nil {
			return err
		}
		if !sameGuards(pk, s.pkLen, wpk) || !sameGuards(sk, s.skLen, wsk) {
			return errors.New("keypair writes outside its outputs")
		}
		sum.add(pk.at(s.pkLen))
		sum.add(sk.at(s.skLen))
		return nil
	}

	agree := func(out, out2, pk, pk2, sk, sk2 *guarded) error {
		wout := pat.frame(out, s.ssLen)
		wpk := pat.frame(pk, s.pkLen)
		wsk := pat.frame(sk, s.skLen)
		err := s.dh(out.at(s.ssLen), pk.at(s.pkLen), sk.at(s.skLen))
		if err != nil {
			return err
		}
		if !sameGuards(out, s.ssLen, wout) {
			return errors.New("dh writes outside its output")
		}
		if !sameFrame(pk, s.pkLen, wpk) || !sameFrame(sk, s.skLen, wsk) {
			return errors.New("dh alters its inputs")
		}
		sum.add(out.at(s.ssLen))
		want := append([]byte(nil), out.at(s.ssLen)...)

		copy(pk2.at(s.pkLen), pk.at(s.pkLen))
		copy(sk2.at(s.skLen), sk.at(s.skLen))
		pat.frame(out2, s.ssLen)
		err = s.dh(out2.at(s.ssLen), pk2.at(s.pkLen), sk2.at(s.skLen))
		if err != nil {
			return err
		}
		if !bytes.Equal(out2.at(s.ssLen), want) {
			return errors.New("dh is not deterministic")
		}

		err = s.dh(pk2.at(s.ssLen), pk2.at(s.pkLen), sk.at(s.skLen))
		if err != nil {
			return err
		}
		if !bytes.Equal(pk2.at(s.ssLen), want) {
			return errors.New("dh output overlapping its public key differs")
		}
		copy(pk2.at(s.pkLen), pk.at(s.pkLen))

		err = s.dh(sk2.at(s.ssLen), pk.at(s.pkLen), sk2.at(s.skLen))
		if err != nil {
			return err
		}
		if !bytes.Equal(sk2.at(s.ssLen), want) {
			return errors.New("dh output overlapping its secret key differs")
		}
		copy(sk2.at(s.skLen), sk.at(s.skLen))
		return nil
	}

	for i := 0; i < loops; i++ {
		if err := keypair(c, a); err != nil {
			return "", err
		}
		if err := keypair(d, b); err != nil {
			return "", err
		}
		if err := agree(e, e2, d, d2, a, a2); err != nil {
			return "", err
		}
		if err := agree(f, e2, c, c2, b, a2); err != nil {
			return "", err
		}
		if !bytes.Equal(e.at(s.ssLen), f.at(s.ssLen)) {
			return "", fmt.Errorf("loop %d: shared secrets differ", i)
		}
	}
	return sum.hex(), nil
}

func ctidhScheme(
	skLen, pkLen int,
	gen func(io.Reader) []byte,
	pub func([]byte) ([]byte, error),
	action func(sk, pk []byte) ([]byte, error),
) dhScheme {
	return dhScheme{
		skLen: skLen, pkLen: pkLen, ssLen: pkLen,
		keypair: func(pk, sk []byte, rng io.Reader) error {
			s := gen(rng)
			p, err := pub(s)
			if err != nil {
				return err
			}
			copy(sk, s)
			copy(pk, p)
			return nil
		},
		dh: func(ss, pk, sk []byte) error {
			r, err := action(append([]byte(nil), sk...),
				append([]byte(nil), pk...))
			if err != nil {
				return err
			}
			copy(ss, r)
			return nil
		},
	}
}

func checksumScheme(size int) dhScheme {
	sc := schemes[size]
	var skLen, pkLen int
	switch size {
	case 511:
		skLen, pkLen = ctidh511.PrivateKeySize, ctidh511.PublicKeySize
	case 512:
		skLen, pkLen = ctidh512.PrivateKeySize, ctidh512.PublicKeySize
	case 1024:
		skLen, pkLen = ctidh1024.PrivateKeySize, ctidh1024.PublicKeySize
	case 2048:
		skLen, pkLen = ctidh2048.PrivateKeySize, ctidh2048.PublicKeySize
	}
	return ctidhScheme(skLen, pkLen, sc.gen, sc.pub, sc.action)
}

func readChecksums(t *testing.T) map[string]string {
	t.Helper()
	f, err := os.Open("checksums")
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	m := map[string]string{}
	s := bufio.NewScanner(f)
	for s.Scan() {
		w := strings.Fields(s.Text())
		if len(w) != 3 {
			t.Fatalf("bad line %q", s.Text())
		}
		m[w[0]+" "+w[1]] = w[2]
	}
	if err := s.Err(); err != nil {
		t.Fatal(err)
	}
	return m
}

func checksumSizes(t *testing.T) []int {
	t.Helper()
	v := os.Getenv("HIGHCTIDH_CHECKSUM_SIZES")
	if v == "" {
		return []int{511, 512}
	}
	var r []int
	for _, w := range strings.Fields(v) {
		n, err := strconv.Atoi(w)
		if err != nil || schemes[n].gen == nil {
			t.Fatalf("bad HIGHCTIDH_CHECKSUM_SIZES %q", v)
		}
		r = append(r, n)
	}
	return r
}

func checksumKinds() []string {
	if os.Getenv("HIGHCTIDH_CHECKSUM_BIG") != "" {
		return []string{"small", "big"}
	}
	return []string{"small"}
}

func loopsFor(kind string) int {
	if kind == "big" {
		return loopsBig
	}
	return loopsSmall
}

func TestChecksums(t *testing.T) {
	want := readChecksums(t)
	for _, size := range checksumSizes(t) {
		for _, kind := range checksumKinds() {
			size, kind := size, kind
			t.Run(fmt.Sprintf("%d/%s", size, kind), func(t *testing.T) {
				key := fmt.Sprintf("%d %s", size, kind)
				w, ok := want[key]
				if !ok {
					t.Fatalf("no committed checksum for %s", key)
				}
				got, err := tryDH(checksumScheme(size), loopsFor(kind))
				if err != nil {
					t.Fatal(err)
				}
				t.Logf("%s %s", key, got)
				if got != w {
					t.Errorf("got %s want %s", got, w)
				}
			})
		}
	}
}
