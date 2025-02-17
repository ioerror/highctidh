package ctidh511

/*
#cgo pkg-config: libhighctidh_511
#define BITS 511
#include "binding511.h"
#include <highctidh.h>

void custom_gen_private(void *const context, highctidh_511_private_key *priv) {
  highctidh_511_csidh_private_withrng(priv, (uintptr_t)context, fillrandom_511_custom);
}
*/
import "C"
import (
	"crypto/hmac"
	"encoding/base64"
	"fmt"
	"io"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

var (
	// PublicKeySize is the size in bytes of the public key.
	PublicKeySize int

	// PrivateKeySize is the size in bytes of the private key.
	PrivateKeySize int
)

// PublicKey is a public CTIDH key.
type PublicKey struct {
	publicKey C.highctidh_511_public_key
}

// NewEmptyPublicKey returns an uninitialized
// PublicKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyPublicKey() *PublicKey {
	return new(PublicKey)
}

// NewPublicKey creates a new public key from
// the given key material or panics if the
// key data is not PublicKeySize.
func NewPublicKey(key []byte) *PublicKey {
	k := new(PublicKey)
	err := k.FromBytes(key)
	if err != nil {
		panic(err)
	}
	return k
}

// String returns a string identifying
// this type as a CTIDH public key.
func (p *PublicKey) String() string {
	return "Ctidh511_PublicKey"
}

// Reset resets the PublicKey to all zeros.
func (p *PublicKey) Reset() {
	zeros := make([]byte, PublicKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes returns the PublicKey as a byte slice.
func (p *PublicKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.publicKey.A.x.c), C.int(((C.BITS+63)/64)*8))
}

// FromBytes loads a PublicKey from the given byte slice.
func (p *PublicKey) FromBytes(data []byte) error {
	if len(data) != PublicKeySize {
		return ErrPublicKeySize
	}

	p.publicKey = *((*C.highctidh_511_public_key)(unsafe.Pointer(&data[0])))
	if !C.highctidh_511_validate(&p.publicKey) {
		return ErrPublicKeyValidation
	}

	return nil
}

// Equal is a constant time comparison of the two public keys.
func (p *PublicKey) Equal(publicKey *PublicKey) bool {
	return hmac.Equal(p.Bytes(), publicKey.Bytes())
}

// Blind performs a blinding operation
// and mutates the public key.
// See notes below about blinding operation with CTIDH.
func (p *PublicKey) Blind(blindingFactor *PrivateKey) error {
	blinded, err := Blind(blindingFactor, p)
	if err != nil {
		panic(err)
	}
	p.publicKey = blinded.publicKey
	return nil
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PublicKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PublicKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PublicKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PublicKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// PrivateKey is a private CTIDH key.
type PrivateKey struct {
	privateKey C.highctidh_511_private_key
}

// NewEmptyPrivateKey returns an uninitialized
// PrivateKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmptyPrivateKey() *PrivateKey {
	return new(PrivateKey)
}

// DeriveSecret derives a shared secret.
func (p *PrivateKey) DeriveSecret(publicKey *PublicKey) []byte {
	return DeriveSecret(p, publicKey)
}

// String returns a string identifying
// this type as a CTIDH private key.
func (p *PrivateKey) String() string {
	return "Ctidh511_PrivateKey"
}

// Reset resets the PrivateKey to all zeros.
func (p *PrivateKey) Reset() {
	zeros := make([]byte, PrivateKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes serializes PrivateKey into a byte slice.
func (p *PrivateKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.privateKey), C.sizeof_highctidh_511_private_key)
}

// FromBytes loads a PrivateKey from the given byte slice.
func (p *PrivateKey) FromBytes(data []byte) error {
	if len(data) != PrivateKeySize {
		return ErrPrivateKeySize
	}

	p.privateKey = *((*C.highctidh_511_private_key)(unsafe.Pointer(&data[0])))
	return nil
}

// Equal is a constant time comparison of the two private keys.
func (p *PrivateKey) Equal(privateKey *PrivateKey) bool {
	return hmac.Equal(p.Bytes(), privateKey.Bytes())
}

// Public returns the public key associated
// with the given private key.
func (p *PrivateKey) Public() *PublicKey {
	return DerivePublicKey(p)
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PrivateKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PrivateKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PrivateKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *PrivateKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// DerivePublicKey derives a public key given a private key.
func DerivePublicKey(privKey *PrivateKey) *PublicKey {
	fmt.Println("out: ba")
	var base C.highctidh_511_public_key = C.highctidh_511_base
	fmt.Println("out: bb")
	baseKey := new(PublicKey)
	fmt.Println("out: bc")
	baseKey.publicKey = base
	fmt.Println("out: bd")
	return GroupAction(privKey, baseKey)
}

// GenerateKeyPair generates a new Ctidh511 private and then
// attempts to compute the Ctidh511 public key.
func GenerateKeyPair() (*PrivateKey, *PublicKey) {
	fmt.Println("out: aa")
	privKey := new(PrivateKey)
	fmt.Println("out: ab")
	C.highctidh_511_csidh_private(&privKey.privateKey)
	fmt.Println("out: ac")
	return privKey, DerivePublicKey(privKey)
}

// GeneratePrivateKey uses the given RNG to derive a new private key.
// This can be used to deterministically generate private keys if the
// entropy source is deterministic, for example an HKDF.
func GeneratePrivateKey(rng io.Reader) *PrivateKey {
	privKey := &PrivateKey{}
// XXX: ? 	privKey := new(PrivateKey)
	p := gopointer.Save(rng)
	C.custom_gen_private(p, &privKey.privateKey)
	gopointer.Unref(p)
	return privKey
}

// GenerateKeyPairWithRNG uses the given RNG to derive a new keypair.
func GenerateKeyPairWithRNG(rng io.Reader) (*PrivateKey, *PublicKey) {
	privKey := GeneratePrivateKey(rng)
	return privKey, DerivePublicKey(privKey)
}

func GroupAction(privateKey *PrivateKey, publicKey *PublicKey) *PublicKey {
	fmt.Println("out: ca")
	sharedKey := new(PublicKey)
// XXX: sharedKey likely needs a correct C.highctidh_511_base
	var base C.highctidh_511_public_key = C.highctidh_511_base
	sharedKey.publicKey = base
	fmt.Println("out: cb")
	ok := C.highctidh_511_csidh(&sharedKey.publicKey, &publicKey.publicKey, &privateKey.privateKey)
	fmt.Println("out: cc")
	if !ok {
		panic(ErrCTIDH)
	}
	fmt.Println("out: cd")
	return sharedKey
}

// DeriveSecret derives a shared secret.
func DeriveSecret(privateKey *PrivateKey, publicKey *PublicKey) []byte {
	sharedSecret := GroupAction(privateKey, publicKey)
	return sharedSecret.Bytes()
}

// Blind performs a blinding operation returning the blinded public key.
func Blind(blindingFactor *PrivateKey, publicKey *PublicKey) (*PublicKey, error) {
	return GroupAction(blindingFactor, publicKey), nil
}

func init() {
	if C.BITS != 511 {
		panic("CTIDH/cgo: C.BITS must match template Bits")
	}
	validateBitSize(C.BITS)
	PublicKeySize = 64
	PrivateKeySize = C.sizeof_highctidh_511_private_key
}
