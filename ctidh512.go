//go:build Ctidh512
// +build Ctidh512

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

/*
#include "binding512.h"
#include <csidh.h>

extern ctidh_fillrandom fillrandom_custom;

__attribute__((weak))
void custom_gen_private(void *const context, private_key *priv) {
  csidh_private_withrng(priv, (uintptr_t)context, fillrandom_custom);
}
*/
import "C"
import (
	"crypto/hmac"
	"encoding/base64"
	"io"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

var (
	// Ctidh512PublicKeySize is the size in bytes of the public key.
	Ctidh512PublicKeySize int

	// Ctidh512PrivateKeySize is the size in bytes of the private key.
	Ctidh512PrivateKeySize int
)

// Ctidh512PublicKey is a public CTIDH key.
type Ctidh512PublicKey struct {
	publicKey C.public_key
}

// NewEmptyCtidh512PublicKey returns an uninitialized
// Ctidh512PublicKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh512PublicKey() *Ctidh512PublicKey {
	return new(Ctidh512PublicKey)
}

// NewCtidh512PublicKey creates a new public key from
// the given key material or panics if the
// key data is not Ctidh512PublicKeySize.
func NewCtidh512PublicKey(key []byte) *Ctidh512PublicKey {
	k := new(Ctidh512PublicKey)
	err := k.FromBytes(key)
	if err != nil {
		panic(err)
	}
	return k
}

// String returns a string identifying
// this type as a CTIDH public key.
func (p *Ctidh512PublicKey) String() string {
	return "Ctidh512_PublicKey"
}

// Reset resets the Ctidh512PublicKey to all zeros.
func (p *Ctidh512PublicKey) Reset() {
	zeros := make([]byte, Ctidh512PublicKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes returns the Ctidh512PublicKey as a byte slice.
func (p *Ctidh512PublicKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.publicKey.A.x.c), C.int(C.UINTBIG_LIMBS*8))
}

// FromBytes loads a Ctidh512PublicKey from the given byte slice.
func (p *Ctidh512PublicKey) FromBytes(data []byte) error {
	if len(data) != Ctidh512PublicKeySize {
		return ErrPublicKeySize
	}

	p.publicKey = *((*C.public_key)(unsafe.Pointer(&data[0])))
	if !C.validate(&p.publicKey) {
		return ErrPublicKeyValidation
	}

	return nil
}

// Equal is a constant time comparison of the two public keys.
func (p *Ctidh512PublicKey) Equal(publicKey *Ctidh512PublicKey) bool {
	return hmac.Equal(p.Bytes(), publicKey.Bytes())
}

// Blind performs a blinding operation
// and mutates the public key.
// See notes below about blinding operation with CTIDH.
func (p *Ctidh512PublicKey) Blind(blindingFactor *Ctidh512PrivateKey) error {
	blinded, err := BlindCtidh512(blindingFactor, p)
	if err != nil {
		panic(err)
	}
	p.publicKey = blinded.publicKey
	return nil
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PublicKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PublicKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PublicKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PublicKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// Ctidh512PrivateKey is a private CTIDH key.
type Ctidh512PrivateKey struct {
	privateKey C.private_key
}

// NewEmptyCtidh512PrivateKey returns an uninitialized
// Ctidh512PrivateKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh512PrivateKey() *Ctidh512PrivateKey {
	return new(Ctidh512PrivateKey)
}

// DeriveSecret derives a shared secret.
func (p *Ctidh512PrivateKey) DeriveSecret(publicKey *Ctidh512PublicKey) []byte {
	return DeriveSecretCtidh512(p, publicKey)
}

// String returns a string identifying
// this type as a CTIDH private key.
func (p *Ctidh512PrivateKey) String() string {
	return "Ctidh512_PrivateKey"
}

// Reset resets the Ctidh512PrivateKey to all zeros.
func (p *Ctidh512PrivateKey) Reset() {
	zeros := make([]byte, Ctidh512PrivateKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes serializes Ctidh512PrivateKey into a byte slice.
func (p *Ctidh512PrivateKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.privateKey), C.primes_num)
}

// FromBytes loads a Ctidh512PrivateKey from the given byte slice.
func (p *Ctidh512PrivateKey) FromBytes(data []byte) error {
	if len(data) != Ctidh512PrivateKeySize {
		return ErrPrivateKeySize
	}

	p.privateKey = *((*C.private_key)(unsafe.Pointer(&data[0])))
	return nil
}

// Equal is a constant time comparison of the two private keys.
func (p *Ctidh512PrivateKey) Equal(privateKey *Ctidh512PrivateKey) bool {
	return hmac.Equal(p.Bytes(), privateKey.Bytes())
}

// Public returns the public key associated
// with the given private key.
func (p *Ctidh512PrivateKey) Public() *Ctidh512PublicKey {
	return DeriveCtidh512PublicKey(p)
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PrivateKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PrivateKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PrivateKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh512PrivateKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// DeriveCtidh512PublicKey derives a public key given a private key.
func DeriveCtidh512PublicKey(privKey *Ctidh512PrivateKey) *Ctidh512PublicKey {
	var base C.public_key
	baseKey := new(Ctidh512PublicKey)
	baseKey.publicKey = base
	return groupActionCtidh512(privKey, baseKey)
}

// GenerateCtidh512KeyPair generates a new Ctidh512 private and then
// attempts to compute the Ctidh512 public key.
func GenerateCtidh512KeyPair() (*Ctidh512PrivateKey, *Ctidh512PublicKey) {
	privKey := new(Ctidh512PrivateKey)
	C.csidh_private(&privKey.privateKey)
	return privKey, DeriveCtidh512PublicKey(privKey)
}

// GenerateCtidh512PrivateKey uses the given RNG to derive a new private key.
// This can be used to deterministically generate private keys if the
// entropy source is deterministic, for example an HKDF.
func GenerateCtidh512PrivateKey(rng io.Reader) *Ctidh512PrivateKey {
	privKey := &Ctidh512PrivateKey{}
	p := gopointer.Save(rng)
	C.custom_gen_private(p, &privKey.privateKey)
	gopointer.Unref(p)
	return privKey
}

// GenerateCtidh512KeyPairWithRNG uses the given RNG to derive a new keypair.
func GenerateCtidh512KeyPairWithRNG(rng io.Reader) (*Ctidh512PrivateKey, *Ctidh512PublicKey) {
	privKey := GenerateCtidh512PrivateKey(rng)
	return privKey, DeriveCtidh512PublicKey(privKey)
}

func groupActionCtidh512(privateKey *Ctidh512PrivateKey, publicKey *Ctidh512PublicKey) *Ctidh512PublicKey {
	sharedKey := new(Ctidh512PublicKey)
	ok := C.csidh(&sharedKey.publicKey, &publicKey.publicKey, &privateKey.privateKey)
	if !ok {
		panic(ErrCTIDH)
	}
	return sharedKey
}

// DeriveSecret derives a shared secret.
func DeriveSecretCtidh512(privateKey *Ctidh512PrivateKey, publicKey *Ctidh512PublicKey) []byte {
	sharedSecret := groupActionCtidh512(privateKey, publicKey)
	return sharedSecret.Bytes()
}

// Blind performs a blinding operation returning the blinded public key.
func BlindCtidh512(blindingFactor *Ctidh512PrivateKey, publicKey *Ctidh512PublicKey) (*Ctidh512PublicKey, error) {
	return groupActionCtidh512(blindingFactor, publicKey), nil
}

func init() {
	if C.BITS != 512 {
		panic("CTIDH/cgo: C.BITS must match template Bits")
	}
	validateBitSize(C.BITS)
	Ctidh512PrivateKeySize = C.primes_num
	switch C.BITS {
	case 511:
		Ctidh512PublicKeySize = 64
	case 512:
		Ctidh512PublicKeySize = 64
	case 1024:
		Ctidh512PublicKeySize = 128
	case 2048:
		Ctidh512PublicKeySize = 256
	}
}
