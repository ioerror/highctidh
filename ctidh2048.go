//go:build Ctidh2048
// +build Ctidh2048

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

/*
#include "binding2048.h"
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
	// Ctidh2048PublicKeySize is the size in bytes of the public key.
	Ctidh2048PublicKeySize int

	// Ctidh2048PrivateKeySize is the size in bytes of the private key.
	Ctidh2048PrivateKeySize int
)

// Ctidh2048PublicKey is a public CTIDH key.
type Ctidh2048PublicKey struct {
	publicKey C.public_key
}

// NewEmptyCtidh2048PublicKey returns an uninitialized
// Ctidh2048PublicKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh2048PublicKey() *Ctidh2048PublicKey {
	return new(Ctidh2048PublicKey)
}

// NewCtidh2048PublicKey creates a new public key from
// the given key material or panics if the
// key data is not Ctidh2048PublicKeySize.
func NewCtidh2048PublicKey(key []byte) *Ctidh2048PublicKey {
	k := new(Ctidh2048PublicKey)
	err := k.FromBytes(key)
	if err != nil {
		panic(err)
	}
	return k
}

// String returns a string identifying
// this type as a CTIDH public key.
func (p *Ctidh2048PublicKey) String() string {
	return "Ctidh2048_PublicKey"
}

// Reset resets the Ctidh2048PublicKey to all zeros.
func (p *Ctidh2048PublicKey) Reset() {
	zeros := make([]byte, Ctidh2048PublicKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes returns the Ctidh2048PublicKey as a byte slice.
func (p *Ctidh2048PublicKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.publicKey.A.x.c), C.int(C.UINTBIG_LIMBS*8))
}

// FromBytes loads a Ctidh2048PublicKey from the given byte slice.
func (p *Ctidh2048PublicKey) FromBytes(data []byte) error {
	if len(data) != Ctidh2048PublicKeySize {
		return ErrPublicKeySize
	}

	p.publicKey = *((*C.public_key)(unsafe.Pointer(&data[0])))
	if !C.validate(&p.publicKey) {
		return ErrPublicKeyValidation
	}

	return nil
}

// Equal is a constant time comparison of the two public keys.
func (p *Ctidh2048PublicKey) Equal(publicKey *Ctidh2048PublicKey) bool {
	return hmac.Equal(p.Bytes(), publicKey.Bytes())
}

// Blind performs a blinding operation
// and mutates the public key.
// See notes below about blinding operation with CTIDH.
func (p *Ctidh2048PublicKey) Blind(blindingFactor *Ctidh2048PrivateKey) error {
	blinded, err := BlindCtidh2048(blindingFactor, p)
	if err != nil {
		panic(err)
	}
	p.publicKey = blinded.publicKey
	return nil
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PublicKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PublicKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PublicKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PublicKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// Ctidh2048PrivateKey is a private CTIDH key.
type Ctidh2048PrivateKey struct {
	privateKey C.private_key
}

// NewEmptyCtidh2048PrivateKey returns an uninitialized
// Ctidh2048PrivateKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh2048PrivateKey() *Ctidh2048PrivateKey {
	return new(Ctidh2048PrivateKey)
}

// DeriveSecret derives a shared secret.
func (p *Ctidh2048PrivateKey) DeriveSecret(publicKey *Ctidh2048PublicKey) []byte {
	return DeriveSecretCtidh2048(p, publicKey)
}

// String returns a string identifying
// this type as a CTIDH private key.
func (p *Ctidh2048PrivateKey) String() string {
	return "Ctidh2048_PrivateKey"
}

// Reset resets the Ctidh2048PrivateKey to all zeros.
func (p *Ctidh2048PrivateKey) Reset() {
	zeros := make([]byte, Ctidh2048PrivateKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes serializes Ctidh2048PrivateKey into a byte slice.
func (p *Ctidh2048PrivateKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.privateKey), C.primes_num)
}

// FromBytes loads a Ctidh2048PrivateKey from the given byte slice.
func (p *Ctidh2048PrivateKey) FromBytes(data []byte) error {
	if len(data) != Ctidh2048PrivateKeySize {
		return ErrPrivateKeySize
	}

	p.privateKey = *((*C.private_key)(unsafe.Pointer(&data[0])))
	return nil
}

// Equal is a constant time comparison of the two private keys.
func (p *Ctidh2048PrivateKey) Equal(privateKey *Ctidh2048PrivateKey) bool {
	return hmac.Equal(p.Bytes(), privateKey.Bytes())
}

// Public returns the public key associated
// with the given private key.
func (p *Ctidh2048PrivateKey) Public() *Ctidh2048PublicKey {
	return DeriveCtidh2048PublicKey(p)
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PrivateKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PrivateKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PrivateKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh2048PrivateKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// DeriveCtidh2048PublicKey derives a public key given a private key.
func DeriveCtidh2048PublicKey(privKey *Ctidh2048PrivateKey) *Ctidh2048PublicKey {
	var base C.public_key
	baseKey := new(Ctidh2048PublicKey)
	baseKey.publicKey = base
	return groupActionCtidh2048(privKey, baseKey)
}

// GenerateCtidh2048KeyPair generates a new Ctidh2048 private and then
// attempts to compute the Ctidh2048 public key.
func GenerateCtidh2048KeyPair() (*Ctidh2048PrivateKey, *Ctidh2048PublicKey) {
	privKey := new(Ctidh2048PrivateKey)
	C.csidh_private(&privKey.privateKey)
	return privKey, DeriveCtidh2048PublicKey(privKey)
}

// GenerateCtidh2048PrivateKey uses the given RNG to derive a new private key.
// This can be used to deterministically generate private keys if the
// entropy source is deterministic, for example an HKDF.
func GenerateCtidh2048PrivateKey(rng io.Reader) *Ctidh2048PrivateKey {
	privKey := &Ctidh2048PrivateKey{}
	p := gopointer.Save(rng)
	C.custom_gen_private(p, &privKey.privateKey)
	gopointer.Unref(p)
	return privKey
}

// GenerateCtidh2048KeyPairWithRNG uses the given RNG to derive a new keypair.
func GenerateCtidh2048KeyPairWithRNG(rng io.Reader) (*Ctidh2048PrivateKey, *Ctidh2048PublicKey) {
	privKey := GenerateCtidh2048PrivateKey(rng)
	return privKey, DeriveCtidh2048PublicKey(privKey)
}

func groupActionCtidh2048(privateKey *Ctidh2048PrivateKey, publicKey *Ctidh2048PublicKey) *Ctidh2048PublicKey {
	sharedKey := new(Ctidh2048PublicKey)
	ok := C.csidh(&sharedKey.publicKey, &publicKey.publicKey, &privateKey.privateKey)
	if !ok {
		panic(ErrCTIDH)
	}
	return sharedKey
}

// DeriveSecret derives a shared secret.
func DeriveSecretCtidh2048(privateKey *Ctidh2048PrivateKey, publicKey *Ctidh2048PublicKey) []byte {
	sharedSecret := groupActionCtidh2048(privateKey, publicKey)
	return sharedSecret.Bytes()
}

// Blind performs a blinding operation returning the blinded public key.
func BlindCtidh2048(blindingFactor *Ctidh2048PrivateKey, publicKey *Ctidh2048PublicKey) (*Ctidh2048PublicKey, error) {
	return groupActionCtidh2048(blindingFactor, publicKey), nil
}

func init() {
	if C.BITS != 2048 {
		panic("CTIDH/cgo: C.BITS must match template Bits")
	}
	validateBitSize(C.BITS)
	Ctidh2048PrivateKeySize = C.primes_num
	switch C.BITS {
	case 511:
		Ctidh2048PublicKeySize = 64
	case 512:
		Ctidh2048PublicKeySize = 64
	case 1024:
		Ctidh2048PublicKeySize = 128
	case 2048:
		Ctidh2048PublicKeySize = 256
	}
}
