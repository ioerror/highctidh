package ctidh1024_test

import (
	"bytes"
	"crypto/sha256"
	"debug/elf"
	"encoding/binary"
	"encoding/hex"
	"errors"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"testing"

	"codeberg.org/vula/highctidh/src/ctidh1024"
	"codeberg.org/vula/highctidh/src/ctidh2048"
	"codeberg.org/vula/highctidh/src/ctidh511"
	"codeberg.org/vula/highctidh/src/ctidh512"
)

type detReader struct {
	seed string
	ctr  uint64
	buf  []byte
}

func (r *detReader) Read(p []byte) (int, error) {
	for len(r.buf) < len(p) {
		var c [8]byte
		binary.LittleEndian.PutUint64(c[:], r.ctr)
		r.ctr++
		h := sha256.Sum256(append([]byte(r.seed), c[:]...))
		r.buf = append(r.buf, h[:]...)
	}
	n := copy(p, r.buf)
	r.buf = r.buf[n:]
	return n, nil
}

type multiScheme struct {
	size       int
	batchSize  []int
	batchBound []int
	gen        func(io.Reader) []byte
	pair       func(io.Reader) ([]byte, []byte)
	derive     func([]byte) []byte
	load       func([]byte) error
}

var multiSchemes = map[string]multiScheme{
	"ctidh511": {
		ctidh511.PrivateKeySize,
		[]int{2, 3, 4, 4, 5, 5, 5, 5, 5, 7, 7, 8, 7, 6, 1},
		[]int{6, 9, 11, 11, 12, 12, 12, 12, 12, 12, 12, 12, 8, 6, 1},
		func(r io.Reader) []byte { return ctidh511.GeneratePrivateKey(r).Bytes() },
		func(r io.Reader) ([]byte, []byte) {
			k, p := ctidh511.GenerateKeyPairWithRNG(r)
			return k.Bytes(), p.Bytes()
		},
		func(b []byte) []byte {
			k := new(ctidh511.PrivateKey)
			if err := k.FromBytes(b); err != nil {
				panic(err)
			}
			return ctidh511.DerivePublicKey(k).Bytes()
		},
		func(b []byte) error { return new(ctidh511.PrivateKey).FromBytes(b) },
	},
	"ctidh512": {
		ctidh512.PrivateKeySize,
		[]int{2, 3, 4, 4, 5, 5, 6, 7, 7, 8, 8, 6, 8, 1},
		[]int{10, 14, 16, 17, 17, 17, 18, 18, 18, 18, 18, 13, 13, 1},
		func(r io.Reader) []byte { return ctidh512.GeneratePrivateKey(r).Bytes() },
		func(r io.Reader) ([]byte, []byte) {
			k, p := ctidh512.GenerateKeyPairWithRNG(r)
			return k.Bytes(), p.Bytes()
		},
		func(b []byte) []byte {
			k := new(ctidh512.PrivateKey)
			if err := k.FromBytes(b); err != nil {
				panic(err)
			}
			return ctidh512.DerivePublicKey(k).Bytes()
		},
		func(b []byte) error { return new(ctidh512.PrivateKey).FromBytes(b) },
	},
	"ctidh1024": {
		ctidh1024.PrivateKeySize,
		[]int{2, 3, 5, 4, 6, 6, 6, 6, 6, 7, 7, 7, 6, 7, 7, 5, 6, 5, 10, 3, 10, 5, 1},
		[]int{2, 4, 5, 5, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 6, 5, 5, 3, 6, 2, 6, 2, 0},
		func(r io.Reader) []byte { return ctidh1024.GeneratePrivateKey(r).Bytes() },
		func(r io.Reader) ([]byte, []byte) {
			k, p := ctidh1024.GenerateKeyPairWithRNG(r)
			return k.Bytes(), p.Bytes()
		},
		func(b []byte) []byte {
			k := new(ctidh1024.PrivateKey)
			if err := k.FromBytes(b); err != nil {
				panic(err)
			}
			return ctidh1024.DerivePublicKey(k).Bytes()
		},
		func(b []byte) error { return new(ctidh1024.PrivateKey).FromBytes(b) },
	},
	"ctidh2048": {
		ctidh2048.PrivateKeySize,
		[]int{9, 10, 8, 8, 7, 10, 12, 11, 10, 15, 10, 9, 8, 6, 10, 13, 10, 9, 12, 13, 10, 10, 10, 1},
		[]int{1, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 3, 2, 0, 2, 0},
		func(r io.Reader) []byte { return ctidh2048.GeneratePrivateKey(r).Bytes() },
		func(r io.Reader) ([]byte, []byte) {
			k, p := ctidh2048.GenerateKeyPairWithRNG(r)
			return k.Bytes(), p.Bytes()
		},
		func(b []byte) []byte {
			k := new(ctidh2048.PrivateKey)
			if err := k.FromBytes(b); err != nil {
				panic(err)
			}
			return ctidh2048.DerivePublicKey(k).Bytes()
		},
		func(b []byte) error { return new(ctidh2048.PrivateKey).FromBytes(b) },
	},
}

var multiVectors = []struct {
	scheme, seed, privateKey, publicKey string
}{
	{"ctidh511", "hpqc ctidh sampler vector 0", "01fe01fefc02010602020105fe010400ff03fd0004ff00010202fe0001f801ff00010100fe06ff00fbfbff000001030201fe02010001020401ff000100020002fd010000feffff0100ff", "bb3ab77ac542fb421e3532a2b20e82197396954f0c8e8dbd85f78525625e3fe17401976570ed148ddc19d5217ef04eca23ab8c006ab74941c776ac5ec9bab251"},
	{"ctidh511", "hpqc ctidh sampler vector 1", "fffefc01fffbfa0000000500060101f900ff03fefe03fefe00010500ff02fcfffdfa020000fdfffe01ff01fe02fe00040201010100fd000206010000ffff02000001fd00ff0300000101", "48ad73a92800ee1211b16453e58f6c9b38ecd466fbd80c6e10ec98fc9d2fa1cfd8daf3f76d2dbf054b5ff0c779806f2b07600ecca44cc3fcc9522f497067cb50"},
	{"ctidh511", "hpqc ctidh sampler vector 2", "fc020403fffff9fd0003000400ff03fe0005fefb010200fdff00ff04ffff01fff8fe02fb00ff00020103ff0104ff0000fefb03ff0200fe000201fe0202ffffff000000ff010001ff0101", "9f4ba13fd7163b49b88588788e968ebfbf71ae6c0dce69671854e63f5fda8b1e7a4da31ccff448d9a6280f7cd65ea4b6cb96d25a7d91096068ce3d596c16ee23"},
	{"ctidh512", "hpqc ctidh sampler vector 0", "01fd0202f7f9fc01fcfd0402fefdfb01fb0102fd00fafffff502ff0002070202ff0302ff04fa010300feff01fd0500060001ff01ff01000000fe03ff00fef801ffff01fdfe020002fe01", "4e990e55383fdbaa0b6fca6c79ae63c4fec1f2312c1c4fa6494a3cd85a0ce10b5a8df154b09ff3db2ed0e118171a65a2db153e4140f331a0296076761c23e61a"},
	{"ctidh512", "hpqc ctidh sampler vector 1", "fefc01010b0900fefc02f7feffff0306ff0402fc07fe0202ff03ff0305fffdff0106fffb01fdfcff0305ff0102fcfe01fbff00fe01fffc00fe0006050302ff020000000000020901ff01", "fbabe97393e3a245092a1505b69c80f28f1fb88549218429a762d6e5622d635605ef957b3d830f009c9bb12804b1af08d496d883f8d883a543a57f4d94c3582f"},
	{"ctidh512", "hpqc ctidh sampler vector 2", "00070202fffff5fc000000fdf50600f8000304fb0300fffc0102fe080103ff01fe05ff05fd02fcfb000003030103fb0001000100fffe04ff000300fc0400ff000002ff0202000001fdff", "29457ae203550fb532a1888290f3b4ec2233f815fa1d752f724d329ddda688cd91ae7bad857d597531786c83340d96d5c3ddbc2d8b2630612fd623b33daa2b0c"},
	{"ctidh1024", "hpqc ctidh sampler vector 0", "010101fe010100fd01000001fe0000fefffffe0000fc00ff00ff0200000100010001ff01000200010101ff010001ffff0100fe0003000001ff00fe0100ff0000feff0002fd0000000002ff0000fe0100000000fcff020100010102000000ff0100000000020100000001fe00000000000000010001000000000000ff000002000000", "d073660d9d916760c2394acd32baf8c6fe7f07c52655c351805ab171ec12c025bed4d90c846f30147be6107afe8f0eb3b4038981cc7cdfe206cd914c8f81cc24b78eb508dae8f90e11e2b6d790b7614024d241e8ea0cc877c4d7cedaa6884d88087a2914cd56d1adb4a85cf50ae61730025ef288dad706ed491a3d8075ce7308"},
	{"ctidh1024", "hpqc ctidh sampler vector 1", "0002fffffe0000000201030001fffc00010000ff010001020100feffffff0000fffffe000200030000feff00ffff0002ff01000000ffff000400fffffd0000000102010100000000fffe000001feff01000002000202fd000000ff02000000010000feff0000000000000000fd0200010001000001fd0000000001ff000002000000", "2eddfe3001bee83aa46ff0748d6981fa534d65acc8130545cad9df0501508e79cca5163afe4aec23839cf34fd3b707b5bb040ad8c8307be0bbc84cef2444bfc2d15e1745a1caa508b148bfc20a919e64db3830081bb50f2e579cb794795decec9ebda4da28ddfe685905c56c4095ec88b3448484f7d663fc00d9928d8609b404"},
	{"ctidh1024", "hpqc ctidh sampler vector 2", "0001ff02010000ff00fdfe00020100010100010100000000fc0000ff00010002fe000000ff00fe01000101ff0301000000000100000001ff040000ffff0001fffe01000000000200010000fefe0100000001fe0002000200ff0001ffff010001010000ff01000001ff0000ff00feffff00000000fd000000020000ffff000000ff00", "ff71e8b8cb06b08502b59ef10c0e3d9206bed231c5065951f0b0b34417868de1408b948939a9250f03cf88c53af7a6a6b043addfc73a01a1f942fee92ee874f79b8737f262a3e5be3a48a184f40f685b91686099bf8dc7dc050fc6f162a9abe4e99a959639898cd8fb1907a08625d0131ebeb9a5eb477c2e4c83f5c7f2c1a305"},
	{"ctidh2048", "hpqc ctidh sampler vector 0", "0000000000000100000000ffff0001000000000100000000ffff000000fe0000000100000000000100020001ff0000000000ff00000000ffff00000000ff00000000010000ff010000000000ff0000000000020000000000fe00ff000000000000000000000000000001000100000001000000010000ff000000000001ffff000000feff000100000000ff0000ff00010000000001010000000000000001ff000000ff000000fe00000000ff0000000000000000ff00ff00000001000100010000ff00000000000000000000000000010100000000000000000000000001000000000100000000", "3ab2e3c41cb0811d08368d65fd0a7a0f07f45ef86c581003af6ee30e7a448a9867ceacb2eb84104879d386085d269f0fb41e1c80a7d88f0d8edbb047f6c267f08f4645b89b9ac133680d6b7f45ab986e2e16de9d55d01f711bc23222c06c692f35579b8045f919330b53e5b76044e6f6b629597247eb18e81c5e65562c63434463e76fbd75654249df0ae93a1c389a5787f2d764bd63280b90f8cfb34f2a7093a7ebbabc18728f9e29483c5bd271c940b6a4b0dcaa16a0f3018246d0d347c6ef502a2cfe23b3a7a6584dcf834765aef8b265c2fdfcd99912d7d72a429ae5b78b1df4cda8325005f9106f26c353fb500b0a030d68fb78e2c399271931c4788e20"},
	{"ctidh2048", "hpqc ctidh sampler vector 1", "00ff00000000000000000000010000ff00000101010000000001000000ff000001010000010001ff000000010000000001010000ff0000000100000000000001000100000000000000ffffff000100000000000000000000000000000000fe000000ff0001ff00010000000000000000fffe0000000000000000000201000000ff0100ff0000ff000000ff000000000000000000ffffff00000000000100ff000000ff00000001000000ffff0000000000ff000101000000000000ff00000000000000010000000000ff0000000000000100000000000000000000000100000000000000000000", "3628cca6e2115148469e970f5c42f2058067c089ff0bcc60c20e87ae8c491c50619cb530482a1790697b28d6b71dcd2c100c229ab8f8f2c1890da2b5c6a84ff16fe1fc42c09c4289d5ceb9790a529883286a0d44c4a09f0521c5592d5be2c4c29cf7742a3ddef8c7867c966ce390791abf26f69fdf0e71ee18dd8d2e375b8b8d5d6879704541437000db2997b66e69a4477bdbbd2cf68c7982db403ca20428bb493993bb8585327f8f0f1cde780bd3c53fa523766a58fe96e90998632a448e5535c5a74de12ddce15481c2980d7479dd3794047cb8a95dac05115e36e508dfa6220b6f70c1d1230cbd7ad6ce642cca1ff425a5724c5cce726495178e5c35fc16"},
	{"ctidh2048", "hpqc ctidh sampler vector 2", "00000000000000000000ff000000ff00ff00000000ff0000ff00ff0002000100000000000000ff00000000ff000000000000fe0000000100000000000100ff00000000010100000000000100000000ffff00ff00000001000000000000000000ff0100000000010000000100ff00000000fe0000000100010000000001000100ff0001ff00000001ff0000ff000000fe00000000ff00000000000000010001000000000000ffff00ff00000000ff0000000000010000000000000000ff0000000000fe0000000000ff0000ff00000000000000000000000000000000000000000000000000fe00", "ee7bf307f9c05997015e12f481f23674abcc8f5f2e13c5bf04bcb6823750b1741bf3ac25eeeb8c247f0b462d58ca9b3cc00aa1eb7f85a80786ce15aa2eae5d46b6a752404d611ba553a83137309be4442246480b2b649c72b5b68a478a4a94bb2b6d4d7c68def8507b710164dd3ef3cb68c677fe2a0de18718f538aed0ef793a23f9a4abd7c83f97d9b5a9ba709901fdeaae95a23108acf17634cba52fd3d733e28ca5e4870145a8583bd7362bdfde05709c85c27f548fe51c8496003f15eda17c0b8c5a9a067bd5d1604929016251d0c8df9f27d2a5bc3ccc12e1cc3dbf27fb058854878c29c1c8218c5b5b229ef6f22d024b7157a350027d203e20c31a4d1b"},
}

func checkWellFormed(t *testing.T, name string, s multiScheme, key []byte) {
	if len(key) != s.size {
		t.Errorf("%s: key length %d, want %d", name, len(key), s.size)
		return
	}
	pos := 0
	for b, w := range s.batchSize {
		l1 := 0
		for _, e := range key[pos : pos+w] {
			if int8(e) < 0 {
				l1 -= int(int8(e))
			} else {
				l1 += int(int8(e))
			}
		}
		if l1 > s.batchBound[b] {
			t.Errorf("%s: batch %d has L1 norm %d, bound %d", name, b, l1, s.batchBound[b])
		}
		pos += w
	}
	if pos != s.size {
		t.Errorf("%s: batches cover %d exponents, want %d", name, pos, s.size)
	}
}

func TestMultiSizeGeneratePrivateKey(t *testing.T) {
	for _, v := range multiVectors {
		s := multiSchemes[v.scheme]
		want, _ := hex.DecodeString(v.privateKey)
		wantPub, _ := hex.DecodeString(v.publicKey)
		checkWellFormed(t, v.scheme+" vector", s, want)
		if got := s.derive(want); !bytes.Equal(got, wantPub) {
			t.Errorf("%s %q: DerivePublicKey of vector key = %x, want %x", v.scheme, v.seed, got, wantPub)
		}
		got := s.gen(&detReader{seed: v.seed})
		checkWellFormed(t, v.scheme, s, got)
		if !bytes.Equal(got, want) {
			t.Errorf("%s %q: GeneratePrivateKey = %x, want %x", v.scheme, v.seed, got, want)
		}
	}
}

func TestMultiSizeGenerateKeyPairWithRNG(t *testing.T) {
	for _, v := range multiVectors {
		s := multiSchemes[v.scheme]
		priv, pub := s.pair(&detReader{seed: v.seed})
		checkWellFormed(t, v.scheme, s, priv)
		if derived := s.derive(priv); !bytes.Equal(pub, derived) {
			t.Errorf("%s %q: public key %x, DerivePublicKey %x", v.scheme, v.seed, pub, derived)
		}
		if hex.EncodeToString(priv) != v.privateKey {
			t.Errorf("%s %q: GenerateKeyPairWithRNG private key = %x, want %s", v.scheme, v.seed, priv, v.privateKey)
		}
		if hex.EncodeToString(pub) != v.publicKey {
			t.Errorf("%s %q: GenerateKeyPairWithRNG public key = %x, want %s", v.scheme, v.seed, pub, v.publicKey)
		}
	}
}

func TestMultiSizePrivateKeyFromBytesRange(t *testing.T) {
	for name, s := range multiSchemes {
		for _, e := range []int8{127, -128} {
			key := bytes.Repeat([]byte{byte(e)}, s.size)
			if s.load(key) == nil {
				t.Errorf("%s: private key with every exponent %d accepted", name, e)
			}
		}
		pos := 0
		for b, w := range s.batchSize {
			bound := s.batchBound[b]
			key := make([]byte, s.size)
			key[pos] = byte(int8(bound))
			if err := s.load(key); err != nil {
				t.Errorf("%s: batch %d at bound %d rejected: %v", name, b, bound, err)
			}
			key[pos] = byte(int8(-bound))
			if err := s.load(key); err != nil {
				t.Errorf("%s: batch %d at bound -%d rejected: %v", name, b, bound, err)
			}
			key[pos] = byte(int8(bound + 1))
			if s.load(key) == nil {
				t.Errorf("%s: batch %d exponent %d over bound %d accepted", name, b, bound+1, bound)
			}
			key[pos] = byte(int8(-bound - 1))
			if s.load(key) == nil {
				t.Errorf("%s: batch %d exponent %d over bound %d accepted", name, b, -bound-1, bound)
			}
			if w > 1 {
				key[pos] = byte(int8(bound))
				key[pos+w-1] = 0xff
				if s.load(key) == nil {
					t.Errorf("%s: batch %d L1 norm %d over bound %d accepted", name, b, bound+1, bound)
				}
			}
			pos += w
		}
		for i := 0; i < 8; i++ {
			key := s.gen(&detReader{seed: name, ctr: uint64(i) << 32})
			if err := s.load(key); err != nil {
				t.Errorf("%s: generated key %x rejected: %v", name, key, err)
			}
		}
	}
	for _, v := range multiVectors {
		key, _ := hex.DecodeString(v.privateKey)
		if err := multiSchemes[v.scheme].load(key); err != nil {
			t.Errorf("%s %q: vector private key rejected: %v", v.scheme, v.seed, err)
		}
	}
}

func TestMultiSizeDistinctSamplers(t *testing.T) {
	seed := "highctidh multi-size sampler"
	k511 := multiSchemes["ctidh511"].gen(&detReader{seed: seed})
	k512 := multiSchemes["ctidh512"].gen(&detReader{seed: seed})
	k1024 := multiSchemes["ctidh1024"].gen(&detReader{seed: seed})
	k2048 := multiSchemes["ctidh2048"].gen(&detReader{seed: seed})
	if bytes.Equal(k511, k1024[:len(k511)]) {
		t.Errorf("ctidh511 key is a prefix of the ctidh1024 key")
	}
	if bytes.Equal(k512, k1024[:len(k512)]) {
		t.Errorf("ctidh512 key is a prefix of the ctidh1024 key")
	}
	if bytes.Equal(k2048, append(append([]byte{}, k1024...), make([]byte, len(k2048)-len(k1024))...)) {
		t.Errorf("ctidh2048 key is the ctidh1024 key followed by zeros")
	}
}

func TestMultiSizeSamplerSymbols(t *testing.T) {
	exe, err := os.Executable()
	if err != nil {
		t.Skipf("cannot locate test binary: %v", err)
	}
	f, err := elf.Open(exe)
	if err != nil {
		t.Skipf("test binary is not ELF: %v", err)
	}
	defer f.Close()
	syms, err := f.Symbols()
	if errors.Is(err, elf.ErrNoSymbols) {
		gotool, lerr := exec.LookPath("go")
		if lerr != nil {
			t.Skipf("test binary is stripped and go is not in PATH: %v", lerr)
		}
		bin := filepath.Join(t.TempDir(), "multisize.test")
		out, berr := exec.Command(gotool, "test", "-c", "-o", bin, ".").CombinedOutput()
		if berr != nil {
			t.Fatalf("go test -c: %v\n%s", berr, out)
		}
		g, gerr := elf.Open(bin)
		if gerr != nil {
			t.Fatal(gerr)
		}
		defer g.Close()
		syms, err = g.Symbols()
	}
	if err != nil {
		t.Fatal(err)
	}
	found := map[string][]elf.Symbol{}
	for _, s := range syms {
		found[s.Name] = append(found[s.Name], s)
	}
	if n := len(found["custom_gen_private"]); n != 0 {
		t.Errorf("custom_gen_private defined %d times, want none", n)
	}
	for _, name := range []string{
		"highctidh_511_custom_gen_private",
		"highctidh_512_custom_gen_private",
		"highctidh_1024_custom_gen_private",
		"highctidh_2048_custom_gen_private",
	} {
		ss := found[name]
		if len(ss) != 1 {
			t.Errorf("%s defined %d times, want 1", name, len(ss))
			continue
		}
		if elf.ST_BIND(ss[0].Info) == elf.STB_WEAK {
			t.Errorf("%s is weak", name)
		}
	}
}
