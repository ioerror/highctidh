package main

var BindingCode = `
// +build {{.Name}}
// DO NOT EDIT: generated code, see gen/main.go

package ctidh

/*
#include "binding{{.Bits}}.h"
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
	"fmt"
	"io"
	"unsafe"

	gopointer "github.com/mattn/go-pointer"
)

var (
	// {{.Name}}PublicKeySize is the size in bytes of the public key.
	{{.Name}}PublicKeySize int

	// {{.Name}}PrivateKeySize is the size in bytes of the private key.
	{{.Name}}PrivateKeySize int
)

// {{.Name}}PublicKey is a public CTIDH key.
type {{.Name}}PublicKey struct {
	publicKey C.public_key
}

// NewEmpty{{.Name}}PublicKey returns an uninitialized
// {{.Name}}PublicKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmpty{{.Name}}PublicKey() *{{.Name}}PublicKey {
	return new({{.Name}}PublicKey)
}

// New{{.Name}}PublicKey creates a new public key from
// the given key material or panics if the
// key data is not {{.Name}}PublicKeySize.
func New{{.Name}}PublicKey(key []byte) *{{.Name}}PublicKey {
	k := new({{.Name}}PublicKey)
	err := k.FromBytes(key)
	if err != nil {
		panic(err)
	}
	return k
}

// String returns a string identifying
// this type as a CTIDH public key.
func (p *{{.Name}}PublicKey) String() string {
	return "{{.Name}}_PublicKey"
}

// Reset resets the {{.Name}}PublicKey to all zeros.
func (p *{{.Name}}PublicKey) Reset() {
	zeros := make([]byte, {{.Name}}PublicKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes returns the {{.Name}}PublicKey as a byte slice.
func (p *{{.Name}}PublicKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.publicKey.A.x.c), C.int(C.UINTBIG_LIMBS*8))
}

// FromBytes loads a {{.Name}}PublicKey from the given byte slice.
func (p *{{.Name}}PublicKey) FromBytes(data []byte) error {
	if len(data) != {{.Name}}PublicKeySize {
		return ErrPublicKeySize
	}

	p.publicKey = *((*C.public_key)(unsafe.Pointer(&data[0])))
	if !C.validate(&p.publicKey) {
		return ErrPublicKeyValidation
	}

	return nil
}

// Equal is a constant time comparison of the two public keys.
func (p *{{.Name}}PublicKey) Equal(publicKey *{{.Name}}PublicKey) bool {
	return hmac.Equal(p.Bytes(), publicKey.Bytes())
}

// Blind performs a blinding operation
// and mutates the public key.
// See notes below about blinding operation with CTIDH.
func (p *{{.Name}}PublicKey) Blind(blindingFactor *{{.Name}}PrivateKey) error {
	blinded, err := Blind{{.Name}}(blindingFactor, p)
	if err != nil {
		panic(err)
	}
	p.publicKey = blinded.publicKey
	return nil
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PublicKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PublicKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PublicKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PublicKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// {{.Name}}PrivateKey is a private CTIDH key.
type {{.Name}}PrivateKey struct {
	privateKey C.private_key
}

// NewEmpty{{.Name}}PrivateKey returns an uninitialized
// {{.Name}}PrivateKey which is suitable to be loaded
// via some serialization format via FromBytes
// or FromPEMFile methods.
func NewEmpty{{.Name}}PrivateKey() *{{.Name}}PrivateKey {
	return new({{.Name}}PrivateKey)
}

// DeriveSecret derives a shared secret.
func (p *{{.Name}}PrivateKey) DeriveSecret(publicKey *{{.Name}}PublicKey) []byte {
	return DeriveSecret{{.Name}}(p, publicKey)
}

// String returns a string identifying
// this type as a CTIDH private key.
func (p *{{.Name}}PrivateKey) String() string {
	return "{{.Name}}_PrivateKey"
}

// Reset resets the {{.Name}}PrivateKey to all zeros.
func (p *{{.Name}}PrivateKey) Reset() {
	zeros := make([]byte, {{.Name}}PrivateKeySize)
	err := p.FromBytes(zeros)
	if err != nil {
		panic(err)
	}
}

// Bytes serializes {{.Name}}PrivateKey into a byte slice.
func (p *{{.Name}}PrivateKey) Bytes() []byte {
	return C.GoBytes(unsafe.Pointer(&p.privateKey), C.primes_num)
}

// FromBytes loads a {{.Name}}PrivateKey from the given byte slice.
func (p *{{.Name}}PrivateKey) FromBytes(data []byte) error {
	if len(data) != {{.Name}}PrivateKeySize {
		return ErrPrivateKeySize
	}

	p.privateKey = *((*C.private_key)(unsafe.Pointer(&data[0])))
	return nil
}

// Equal is a constant time comparison of the two private keys.
func (p *{{.Name}}PrivateKey) Equal(privateKey *{{.Name}}PrivateKey) bool {
	return hmac.Equal(p.Bytes(), privateKey.Bytes())
}

// Public returns the public key associated
// with the given private key.
func (p *{{.Name}}PrivateKey) Public() *{{.Name}}PublicKey {
	return Derive{{.Name}}PublicKey(p)
}

// MarshalBinary is an implementation of a method on the
// BinaryMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PrivateKey) MarshalBinary() ([]byte, error) {
	return p.Bytes(), nil
}

// UnmarshalBinary is an implementation of a method on the
// BinaryUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PrivateKey) UnmarshalBinary(data []byte) error {
	return p.FromBytes(data)
}

// MarshalText is an implementation of a method on the
// TextMarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PrivateKey) MarshalText() ([]byte, error) {
	return []byte(base64.StdEncoding.EncodeToString(p.Bytes())), nil
}

// UnmarshalText is an implementation of a method on the
// TextUnmarshaler interface defined in https://golang.org/pkg/encoding/
func (p *{{.Name}}PrivateKey) UnmarshalText(data []byte) error {
	raw, err := base64.StdEncoding.DecodeString(string(data))
	if err != nil {
		return err
	}
	return p.FromBytes(raw)
}

// Derive{{.Name}}PublicKey derives a public key given a private key.
func Derive{{.Name}}PublicKey(privKey *{{.Name}}PrivateKey) *{{.Name}}PublicKey {
	var base C.public_key
	baseKey := new({{.Name}}PublicKey)
	baseKey.publicKey = base
	return groupAction{{.Name}}(privKey, baseKey)
}

// Generate{{.Name}}KeyPair generates a new {{.Name}} private and then
// attempts to compute the {{.Name}} public key.
func Generate{{.Name}}KeyPair() (*{{.Name}}PrivateKey, *{{.Name}}PublicKey) {
	privKey := new({{.Name}}PrivateKey)
	C.csidh_private(&privKey.privateKey)
	return privKey, Derive{{.Name}}PublicKey(privKey)
}

// Generate{{.Name}}PrivateKey uses the given RNG to derive a new private key.
// This can be used to deterministically generate private keys if the
// entropy source is deterministic, for example an HKDF.
func Generate{{.Name}}PrivateKey(rng io.Reader) *{{.Name}}PrivateKey {
	privKey := &{{.Name}}PrivateKey{}
	p := gopointer.Save(rng)
	C.custom_gen_private(p, &privKey.privateKey)
	gopointer.Unref(p)
	return privKey
}

// Generate{{.Name}}KeyPairWithRNG uses the given RNG to derive a new keypair.
func Generate{{.Name}}KeyPairWithRNG(rng io.Reader) (*{{.Name}}PrivateKey, *{{.Name}}PublicKey) {
	privKey := Generate{{.Name}}PrivateKey(rng)
	return privKey, Derive{{.Name}}PublicKey(privKey)
}

func groupAction{{.Name}}(privateKey *{{.Name}}PrivateKey, publicKey *{{.Name}}PublicKey) *{{.Name}}PublicKey {
	sharedKey := new({{.Name}}PublicKey)
	ok := C.csidh(&sharedKey.publicKey, &publicKey.publicKey, &privateKey.privateKey)
	if !ok {
		panic(ErrCTIDH)
	}
	return sharedKey
}

// DeriveSecret derives a shared secret.
func DeriveSecret{{.Name}}(privateKey *{{.Name}}PrivateKey, publicKey *{{.Name}}PublicKey) []byte {
	sharedSecret := groupAction{{.Name}}(privateKey, publicKey)
	return sharedSecret.Bytes()
}

// Blind performs a blinding operation returning the blinded public key.
func Blind{{.Name}}(blindingFactor *{{.Name}}PrivateKey, publicKey *{{.Name}}PublicKey) (*{{.Name}}PublicKey, error) {
	return groupAction{{.Name}}(blindingFactor, publicKey), nil
}

func init() {
	if C.BITS != {{.Bits}} {
		panic("CTIDH/cgo: C.BITS must match template Bits")
	}
	validateBitSize(C.BITS)
	{{.Name}}PrivateKeySize = C.primes_num
	switch C.BITS {
	case 511:
		{{.Name}}PublicKeySize = 64
	case 512:
		{{.Name}}PublicKeySize = 64
	case 1024:
		{{.Name}}PublicKeySize = 128
	case 2048:
		{{.Name}}PublicKeySize = 256
	}
}

`
