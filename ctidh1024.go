//go:build Ctidh1024
// +build Ctidh1024

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

/*
#include "binding1024.h"
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
	// Ctidh1024PublicKeySize is the size in bytes of the public key.
	Ctidh1024PublicKeySize int

	// Ctidh1024PrivateKeySize is the size in bytes of the private key.
	Ctidh1024PrivateKeySize int
)

// Ctidh1024PublicKey is a public CTIDH key.
type Ctidh1024PublicKey struct {
	publicKey C.public_key
}

// NewEmptyCtidh1024PublicKey returns an uninitialized
// Ctidh1024PublicKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh1024PublicKey() *Ctidh1024PublicKey {
	return new(Ctidh1024PublicKey)
}

// NewCtidh1024PublicKey creates a new public key from
// the given key material or panics if the
// key data is not Ctidh1024PublicKeySize.
func NewCtidh1024PublicKey(key []byte) *Ctidh1024PublicKey {
	k := new(Ctidh1024PublicKey)
	err := k.FromBytes(key)
	if err != nil {
		panic(err)
	}
	return k
}

// String returns a string identifying
// this type as a CTIDH public key.
func (p *Ctidh1024PublicKey) String() string {
	return "Ctidh1024_PublicKey"
}

// Reset resets the Ctidh1024PublicKey to all zeros.
func (p *Ctidh1024PublicKey) Reset() {
	zeros := make([]byte, Ctidh1024PublicKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes returns the Ctidh1024PublicKey as a byte slice.
func (p *Ctidh1024PublicKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.publicKey.A.x.c), C.int(C.UINTBIG_LIMBS*8))
}

// FromBytes loads a Ctidh1024PublicKey from the given byte slice.
func (p *Ctidh1024PublicKey) FromBytes(data []byte) error {
	if len(data) != Ctidh1024PublicKeySize {
		return ErrPublicKeySize
	}

	p.publicKey = *((*C.public_key)(unsafe.Pointer(&data[0])))
	if !C.validate(&p.publicKey) {
		return ErrPublicKeyValidation
	}

	return nil
}

// Equal is a constant time comparison of the two public keys.
func (p *Ctidh1024PublicKey) Equal(publicKey *Ctidh1024PublicKey) bool {
	return hmac.Equal(p.Bytes(), publicKey.Bytes())
}

// Blind performs a blinding operation
// and mutates the public key.
// See notes below about blinding operation with CTIDH.
func (p *Ctidh1024PublicKey) Blind(blindingFactor *Ctidh1024PrivateKey) error {
	blinded, err := BlindCtidh1024(blindingFactor, p)
	if err != nil {
		panic(err)
	}
	p.publicKey = blinded.publicKey
	return nil
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PublicKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PublicKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PublicKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PublicKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// Ctidh1024PrivateKey is a private CTIDH key.
type Ctidh1024PrivateKey struct {
	privateKey C.private_key
}

// NewEmptyCtidh1024PrivateKey returns an uninitialized
// Ctidh1024PrivateKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh1024PrivateKey() *Ctidh1024PrivateKey {
	return new(Ctidh1024PrivateKey)
}

// DeriveSecret derives a shared secret.
func (p *Ctidh1024PrivateKey) DeriveSecret(publicKey *Ctidh1024PublicKey) []byte {
	return DeriveSecretCtidh1024(p, publicKey)
}

// String returns a string identifying
// this type as a CTIDH private key.
func (p *Ctidh1024PrivateKey) String() string {
	return "Ctidh1024_PrivateKey"
}

// Reset resets the Ctidh1024PrivateKey to all zeros.
func (p *Ctidh1024PrivateKey) Reset() {
	zeros := make([]byte, Ctidh1024PrivateKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes serializes Ctidh1024PrivateKey into a byte slice.
func (p *Ctidh1024PrivateKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.privateKey), C.primes_num)
}

// FromBytes loads a Ctidh1024PrivateKey from the given byte slice.
func (p *Ctidh1024PrivateKey) FromBytes(data []byte) error {
	if len(data) != Ctidh1024PrivateKeySize {
		return ErrPrivateKeySize
	}

	p.privateKey = *((*C.private_key)(unsafe.Pointer(&data[0])))
	return nil
}

// Equal is a constant time comparison of the two private keys.
func (p *Ctidh1024PrivateKey) Equal(privateKey *Ctidh1024PrivateKey) bool {
	return hmac.Equal(p.Bytes(), privateKey.Bytes())
}

// Public returns the public key associated
// with the given private key.
func (p *Ctidh1024PrivateKey) Public() *Ctidh1024PublicKey {
	return DeriveCtidh1024PublicKey(p)
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PrivateKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PrivateKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PrivateKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh1024PrivateKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// DeriveCtidh1024PublicKey derives a public key given a private key.
func DeriveCtidh1024PublicKey(privKey *Ctidh1024PrivateKey) *Ctidh1024PublicKey {
	var base C.public_key
	baseKey := new(Ctidh1024PublicKey)
	baseKey.publicKey = base
	return groupActionCtidh1024(privKey, baseKey)
}

// GenerateCtidh1024KeyPair generates a new Ctidh1024 private and then
// attempts to compute the Ctidh1024 public key.
func GenerateCtidh1024KeyPair() (*Ctidh1024PrivateKey, *Ctidh1024PublicKey) {
	privKey := new(Ctidh1024PrivateKey)
	C.csidh_private(&privKey.privateKey)
	return privKey, DeriveCtidh1024PublicKey(privKey)
}

// GenerateCtidh1024PrivateKey uses the given RNG to derive a new private key.
// This can be used to deterministically generate private keys if the
// entropy source is deterministic, for example an HKDF.
func GenerateCtidh1024PrivateKey(rng io.Reader) *Ctidh1024PrivateKey {
	privKey := &Ctidh1024PrivateKey{}
	p := gopointer.Save(rng)
	C.custom_gen_private(p, &privKey.privateKey)
	gopointer.Unref(p)
	return privKey
}

// GenerateCtidh1024KeyPairWithRNG uses the given RNG to derive a new keypair.
func GenerateCtidh1024KeyPairWithRNG(rng io.Reader) (*Ctidh1024PrivateKey, *Ctidh1024PublicKey) {
	privKey := GenerateCtidh1024PrivateKey(rng)
	return privKey, DeriveCtidh1024PublicKey(privKey)
}

func groupActionCtidh1024(privateKey *Ctidh1024PrivateKey, publicKey *Ctidh1024PublicKey) *Ctidh1024PublicKey {
	sharedKey := new(Ctidh1024PublicKey)
	ok := C.csidh(&sharedKey.publicKey, &publicKey.publicKey, &privateKey.privateKey)
	if !ok {
		panic(ErrCTIDH)
	}
	return sharedKey
}

// DeriveSecret derives a shared secret.
func DeriveSecretCtidh1024(privateKey *Ctidh1024PrivateKey, publicKey *Ctidh1024PublicKey) []byte {
	sharedSecret := groupActionCtidh1024(privateKey, publicKey)
	return sharedSecret.Bytes()
}

// Blind performs a blinding operation returning the blinded public key.
func BlindCtidh1024(blindingFactor *Ctidh1024PrivateKey, publicKey *Ctidh1024PublicKey) (*Ctidh1024PublicKey, error) {
	return groupActionCtidh1024(blindingFactor, publicKey), nil
}

func init() {
	if C.BITS != 1024 {
		panic("CTIDH/cgo: C.BITS must match template Bits")
	}
	validateBitSize(C.BITS)
	Ctidh1024PrivateKeySize = C.primes_num
	switch C.BITS {
	case 511:
		Ctidh1024PublicKeySize = 64
	case 512:
		Ctidh1024PublicKeySize = 64
	case 1024:
		Ctidh1024PublicKeySize = 128
	case 2048:
		Ctidh1024PublicKeySize = 256
	}
}
