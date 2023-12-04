//go:build Ctidh511
// +build Ctidh511

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

/*
#include "binding511.h"
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
	// Ctidh511PublicKeySize is the size in bytes of the public key.
	Ctidh511PublicKeySize int

	// Ctidh511PrivateKeySize is the size in bytes of the private key.
	Ctidh511PrivateKeySize int
)

// Ctidh511PublicKey is a public CTIDH key.
type Ctidh511PublicKey struct {
	publicKey C.public_key
}

// NewEmptyCtidh511PublicKey returns an uninitialized
// Ctidh511PublicKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh511PublicKey() *Ctidh511PublicKey {
	return new(Ctidh511PublicKey)
}

// NewCtidh511PublicKey creates a new public key from
// the given key material or panics if the
// key data is not Ctidh511PublicKeySize.
func NewCtidh511PublicKey(key []byte) *Ctidh511PublicKey {
	k := new(Ctidh511PublicKey)
	err := k.FromBytes(key)
	if err != nil {
		panic(err)
	}
	return k
}

// String returns a string identifying
// this type as a CTIDH public key.
func (p *Ctidh511PublicKey) String() string {
	return "Ctidh511_PublicKey"
}

// Reset resets the Ctidh511PublicKey to all zeros.
func (p *Ctidh511PublicKey) Reset() {
	zeros := make([]byte, Ctidh511PublicKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes returns the Ctidh511PublicKey as a byte slice.
func (p *Ctidh511PublicKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.publicKey.A.x.c), C.int(C.UINTBIG_LIMBS*8))
}

// FromBytes loads a Ctidh511PublicKey from the given byte slice.
func (p *Ctidh511PublicKey) FromBytes(data []byte) error {
	if len(data) != Ctidh511PublicKeySize {
		return ErrPublicKeySize
	}

	p.publicKey = *((*C.public_key)(unsafe.Pointer(&data[0])))
	if !C.validate(&p.publicKey) {
		return ErrPublicKeyValidation
	}

	return nil
}

// Equal is a constant time comparison of the two public keys.
func (p *Ctidh511PublicKey) Equal(publicKey *Ctidh511PublicKey) bool {
	return hmac.Equal(p.Bytes(), publicKey.Bytes())
}

// Blind performs a blinding operation
// and mutates the public key.
// See notes below about blinding operation with CTIDH.
func (p *Ctidh511PublicKey) Blind(blindingFactor *Ctidh511PrivateKey) error {
	blinded, err := BlindCtidh511(blindingFactor, p)
	if err != nil {
		panic(err)
	}
	p.publicKey = blinded.publicKey
	return nil
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PublicKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PublicKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PublicKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PublicKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// Ctidh511PrivateKey is a private CTIDH key.
type Ctidh511PrivateKey struct {
	privateKey C.private_key
}

// NewEmptyCtidh511PrivateKey returns an uninitialized
// Ctidh511PrivateKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyCtidh511PrivateKey() *Ctidh511PrivateKey {
	return new(Ctidh511PrivateKey)
}

// DeriveSecret derives a shared secret.
func (p *Ctidh511PrivateKey) DeriveSecret(publicKey *Ctidh511PublicKey) []byte {
	return DeriveSecretCtidh511(p, publicKey)
}

// String returns a string identifying
// this type as a CTIDH private key.
func (p *Ctidh511PrivateKey) String() string {
	return "Ctidh511_PrivateKey"
}

// Reset resets the Ctidh511PrivateKey to all zeros.
func (p *Ctidh511PrivateKey) Reset() {
	zeros := make([]byte, Ctidh511PrivateKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes serializes Ctidh511PrivateKey into a byte slice.
func (p *Ctidh511PrivateKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.privateKey), C.primes_num)
}

// FromBytes loads a Ctidh511PrivateKey from the given byte slice.
func (p *Ctidh511PrivateKey) FromBytes(data []byte) error {
	if len(data) != Ctidh511PrivateKeySize {
		return ErrPrivateKeySize
	}

	p.privateKey = *((*C.private_key)(unsafe.Pointer(&data[0])))
	return nil
}

// Equal is a constant time comparison of the two private keys.
func (p *Ctidh511PrivateKey) Equal(privateKey *Ctidh511PrivateKey) bool {
	return hmac.Equal(p.Bytes(), privateKey.Bytes())
}

// Public returns the public key associated
// with the given private key.
func (p *Ctidh511PrivateKey) Public() *Ctidh511PublicKey {
	return DeriveCtidh511PublicKey(p)
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PrivateKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PrivateKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PrivateKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *Ctidh511PrivateKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// DeriveCtidh511PublicKey derives a public key given a private key.
func DeriveCtidh511PublicKey(privKey *Ctidh511PrivateKey) *Ctidh511PublicKey {
	var base C.public_key
	baseKey := new(Ctidh511PublicKey)
	baseKey.publicKey = base
	return groupActionCtidh511(privKey, baseKey)
}

// GenerateCtidh511KeyPair generates a new Ctidh511 private and then
// attempts to compute the Ctidh511 public key.
func GenerateCtidh511KeyPair() (*Ctidh511PrivateKey, *Ctidh511PublicKey) {
	privKey := new(Ctidh511PrivateKey)
	C.csidh_private(&privKey.privateKey)
	return privKey, DeriveCtidh511PublicKey(privKey)
}

// GenerateCtidh511PrivateKey uses the given RNG to derive a new private key.
// This can be used to deterministically generate private keys if the
// entropy source is deterministic, for example an HKDF.
func GenerateCtidh511PrivateKey(rng io.Reader) *Ctidh511PrivateKey {
	privKey := &Ctidh511PrivateKey{}
	p := gopointer.Save(rng)
	C.custom_gen_private(p, &privKey.privateKey)
	gopointer.Unref(p)
	return privKey
}

// GenerateCtidh511KeyPairWithRNG uses the given RNG to derive a new keypair.
func GenerateCtidh511KeyPairWithRNG(rng io.Reader) (*Ctidh511PrivateKey, *Ctidh511PublicKey) {
	privKey := GenerateCtidh511PrivateKey(rng)
	return privKey, DeriveCtidh511PublicKey(privKey)
}

func groupActionCtidh511(privateKey *Ctidh511PrivateKey, publicKey *Ctidh511PublicKey) *Ctidh511PublicKey {
	sharedKey := new(Ctidh511PublicKey)
	ok := C.csidh(&sharedKey.publicKey, &publicKey.publicKey, &privateKey.privateKey)
	if !ok {
		panic(ErrCTIDH)
	}
	return sharedKey
}

// DeriveSecret derives a shared secret.
func DeriveSecretCtidh511(privateKey *Ctidh511PrivateKey, publicKey *Ctidh511PublicKey) []byte {
	sharedSecret := groupActionCtidh511(privateKey, publicKey)
	return sharedSecret.Bytes()
}

// Blind performs a blinding operation returning the blinded public key.
func BlindCtidh511(blindingFactor *Ctidh511PrivateKey, publicKey *Ctidh511PublicKey) (*Ctidh511PublicKey, error) {
	return groupActionCtidh511(blindingFactor, publicKey), nil
}

func init() {
	if C.BITS != 511 {
		panic("CTIDH/cgo: C.BITS must match template Bits")
	}
	validateBitSize(C.BITS)
	Ctidh511PrivateKeySize = C.primes_num
	switch C.BITS {
	case 511:
		Ctidh511PublicKeySize = 64
	case 512:
		Ctidh511PublicKeySize = 64
	case 1024:
		Ctidh511PublicKeySize = 128
	case 2048:
		Ctidh511PublicKeySize = 256
	}
}
