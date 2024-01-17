package interfaces

import (
	"io"
    "codeberg.org/vula/highctidh/ctidh511"
    "codeberg.org/vula/highctidh/ctidh512"
    "codeberg.org/vula/highctidh/ctidh1024"
    "codeberg.org/vula/highctidh/ctidh2048"
)

// CTIDH interface allows the usage of the CTIDH NIKE and it's key types.
type CTIDH interface {

	// Name returns the string name uniquely identifying the particular CTIDH variant.
	Name() string

	// PrivateKeyFromBytes restores a PrivateKey instance from bytes in the canonical
	// (little-endian) representation.
	PrivateKeyFromBytes(b []byte) PrivateKey

	// PublicKeyFromBytes restores a public_key instance from bytes in the canonical
	// (little-endian) representation.
	PublicKeyFromBytes(b []byte) PublicKey

	// PrivateKeyFromHex restores a private_key instance from hex in the canonical
	// (little-endian) representation.
	PrivateKeyFromHex(h string) PrivateKey

	// PublicKeyFromHex restores a public_key instance from hex in the canonical
	// (little-endian) representation
	PublicKeyFromHex(h string) PublicKey

	// Validate returns true if key is valid, otherwise returns false.
	Validate(pubkey PublicKey) bool

	GenerateSecretKey(rng io.Reader) PrivateKey

	DerivePublicKey(privKey PrivateKey) PublicKey

	DH(privKey PrivateKey, pubKey PublicKey) PublicKey

	Blind(blindingFactor PrivateKey, pubKey PublicKey) PublicKey
}

type PrivateKey interface {
	Bytes() []byte

	FromBytes(b []byte)

	FromHex(h string)

	DerivePublicKey() PublicKey
}

type PublicKey interface {
	Bytes() []byte

	FromBytes(b []byte)

	FromHex(h string)
}
