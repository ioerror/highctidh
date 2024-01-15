package ctidh512

import (
	"encoding/hex"
	"io"

	"codeberg.org/vula/highctidh/interfaces"
	"codeberg.org/vula/highctidh/schemes"
	"codeberg.org/vula/highctidh/types"
    ctidh "codeberg.org/vula/highctidh/ctidh512"
)

func init() {
	scheme := New()
	err := schemes.Register(scheme.Name(), scheme)
	if err != nil {
		panic(err)
	}
}

func New() interfaces.CTIDH {
	return &types.CTIDH{
		CTIDH: &CTIDH512{},
	}
}

// the below lines use compiler type checking to ensure
// we are correctly implementing the interfaces:
var _ interfaces.CTIDH = (*CTIDH512)(nil)
var _ interfaces.PrivateKey = (*PrivateKey)(nil)
var _ interfaces.PublicKey = (*PublicKey)(nil)

type CTIDH512 struct{}

func (c *CTIDH512) Name() string {
	return "CTIDH512"
}

func (c *CTIDH512) Blind(blindingFactor interfaces.PrivateKey, pubKey interfaces.PublicKey) interfaces.PublicKey {
	out, err := ctidh.Blind(blindingFactor.(*PrivateKey).privKey, pubKey.(*PublicKey).pubKey)
	if err != nil {
		panic(err)
	}
	return &PublicKey{
		pubKey: out,
	}
}

func (c *CTIDH512) DH(privKey interfaces.PrivateKey, pubKey interfaces.PublicKey) interfaces.PublicKey {
	return &PublicKey{
		pubKey: ctidh.GroupAction(privKey.(*PrivateKey).privKey, pubKey.(*PublicKey).pubKey),
	}
}

func (c *CTIDH512) DerivePublicKey(privKey interfaces.PrivateKey) interfaces.PublicKey {
	return &PublicKey{
		pubKey: ctidh.DerivePublicKey(privKey.(*PrivateKey).privKey),
	}
}

func (c *CTIDH512) GenerateSecretKey(rng io.Reader) interfaces.PrivateKey {
	return &PrivateKey{
		privKey: ctidh.GeneratePrivateKey(rng),
	}
}

func (c *CTIDH512) Validate(pubkey interfaces.PublicKey) bool {
	p := ctidh.NewEmptyPublicKey()
	err := p.FromBytes(pubkey.Bytes())
	if err != nil {
		return false
	}
	return true
}

func (c *CTIDH512) PrivateKeyFromBytes(b []byte) interfaces.PrivateKey {
	privKey := ctidh.NewEmptyPrivateKey()
	err := privKey.FromBytes(b)
	if err != nil {
		panic(err)
	}
	return &PrivateKey{
		privKey: privKey,
	}
}

func (c *CTIDH512) PrivateKeyFromHex(h string) interfaces.PrivateKey {
	v, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return c.PrivateKeyFromBytes(v)
}

func (c *CTIDH512) PublicKeyFromBytes(b []byte) interfaces.PublicKey {
	pubKey := ctidh.NewEmptyPublicKey()
	err := pubKey.FromBytes(b)
	if err != nil {
		panic(err)
	}
	return &PublicKey{
		pubKey: pubKey,
	}
}

func (c *CTIDH512) PublicKeyFromHex(h string) interfaces.PublicKey {
	v, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	return c.PublicKeyFromBytes(v)
}

type PrivateKey struct {
	privKey *ctidh.PrivateKey
}

func (p *PrivateKey) Bytes() []byte {
	return p.privKey.Bytes()
}

func (p *PrivateKey) FromBytes(b []byte) {
	err := p.privKey.FromBytes(b)
	if err != nil {
		panic(err)
	}
}

func (p *PrivateKey) FromHex(h string) {
	v, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	p.privKey.FromBytes(v)
}

func (p *PrivateKey) DerivePublicKey() interfaces.PublicKey {
	pubKey := ctidh.DerivePublicKey(p.privKey)
	return &PublicKey{
		pubKey: pubKey,
	}
}

type PublicKey struct {
	pubKey *ctidh.PublicKey
}

func (p *PublicKey) Bytes() []byte {
	return p.pubKey.Bytes()
}

func (p *PublicKey) FromBytes(b []byte) {
	err := p.pubKey.FromBytes(b)
	if err != nil {
		panic(err)
	}
}

func (p *PublicKey) FromHex(h string) {
	v, err := hex.DecodeString(h)
	if err != nil {
		panic(err)
	}
	p.pubKey.FromBytes(v)
}
