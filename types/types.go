package types

import (
	"io"

	"codeberg.org/vula/highctidh/interfaces"
)

// CTIDH wrapper type
type CTIDH struct {
	CTIDH interfaces.CTIDH
}

// the below line uses compiler type checking to ensure
// we are correctly implementing the interfaces:
var _ interfaces.CTIDH = (*CTIDH)(nil)

func (c *CTIDH) Name() string {
	return c.CTIDH.Name()
}

func (c *CTIDH) Blind(blindingFactor interfaces.PrivateKey, pubKey interfaces.PublicKey) interfaces.PublicKey {
	return c.CTIDH.Blind(blindingFactor, pubKey)
}

func (c *CTIDH) DH(privKey interfaces.PrivateKey, pubKey interfaces.PublicKey) interfaces.PublicKey {
	return c.CTIDH.DH(privKey, pubKey)
}

func (c *CTIDH) DerivePublicKey(privKey interfaces.PrivateKey) interfaces.PublicKey {
	return c.CTIDH.DerivePublicKey(privKey)
}

func (c *CTIDH) GenerateSecretKey(rng io.Reader) interfaces.PrivateKey {
	return c.CTIDH.GenerateSecretKey(rng)
}

func (c *CTIDH) Validate(pubkey interfaces.PublicKey) bool {
	return c.CTIDH.Validate(pubkey)
}

func (c *CTIDH) PublicKeyFromHex(h string) interfaces.PublicKey {
	return c.CTIDH.PublicKeyFromHex(h)
}

func (c *CTIDH) PrivateKeyFromHex(h string) interfaces.PrivateKey {
	return c.CTIDH.PrivateKeyFromHex(h)
}

func (c *CTIDH) PublicKeyFromBytes(b []byte) interfaces.PublicKey {
	return c.CTIDH.PublicKeyFromBytes(b)
}

func (c *CTIDH) PrivateKeyFromBytes(b []byte) interfaces.PrivateKey {
	return c.CTIDH.PrivateKeyFromBytes(b)
}
