package main

var BindingTests = `
// +build {{.Name}}
// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func Test{{.Name}}BlindingOperation(t *testing.T) {
	mixPrivateKey, mixPublicKey := Generate{{.Name}}KeyPair()
	clientPrivateKey, clientPublicKey := Generate{{.Name}}KeyPair()

	blindingFactor := Generate{{.Name}}PrivateKey(rand.Reader)
	value1, err := Blind{{.Name}}(blindingFactor, New{{.Name}}PublicKey(DeriveSecret{{.Name}}(clientPrivateKey, mixPublicKey)))
	require.NoError(t, err)
	blinded, err := Blind{{.Name}}(blindingFactor, clientPublicKey)
	require.NoError(t, err)
	value2 := DeriveSecret{{.Name}}(mixPrivateKey, blinded)

	require.Equal(t, value1.Bytes(), value2)
}

func TestGenerate{{.Name}}KeyPairWithRNG(t *testing.T) {
	privateKey, publicKey := Generate{{.Name}}KeyPairWithRNG(rand.Reader)
	zeros := make([]byte, {{.Name}}PublicKeySize)
	require.NotEqual(t, privateKey.Bytes(), zeros)
	require.NotEqual(t, publicKey.Bytes(), zeros)
}

func Test{{.Name}}PublicKeyReset(t *testing.T) {
	zeros := make([]byte, {{.Name}}PublicKeySize)
	_, publicKey := Generate{{.Name}}KeyPair()
	require.NotEqual(t, publicKey.Bytes(), zeros)

	publicKey.Reset()
	require.Equal(t, publicKey.Bytes(), zeros)
}

func Test{{.Name}}PrivateKeyReset(t *testing.T) {
	zeros := make([]byte, {{.Name}}PrivateKeySize)
	privateKey, _ := Generate{{.Name}}KeyPair()
	require.NotEqual(t, privateKey.Bytes(), zeros)

	privateKey.Reset()
	require.Equal(t, privateKey.Bytes(), zeros)
}

func Test{{.Name}}PublicKeyMarshaling(t *testing.T) {
	privKey, publicKey := Generate{{.Name}}KeyPair()
	publicKeyBytes := publicKey.Bytes()

	publicKey2 := new({{.Name}}PublicKey)
	err := publicKey2.FromBytes(publicKeyBytes)
	require.NoError(t, err)

	publicKey2Bytes := publicKey2.Bytes()

	publicKey3 := Derive{{.Name}}PublicKey(privKey)
	publicKey3Bytes := publicKey3.Bytes()

	require.Equal(t, publicKeyBytes, publicKey2Bytes)
	require.Equal(t, publicKey3Bytes, publicKeyBytes)
}

func Test{{.Name}}PrivateKeyBytesing(t *testing.T) {
	privateKey, _ := Generate{{.Name}}KeyPair()
	privateKeyBytes := privateKey.Bytes()

	privateKey2 := new({{.Name}}PrivateKey)
	privateKey2.FromBytes(privateKeyBytes)
	privateKey2Bytes := privateKey2.Bytes()

	require.Equal(t, privateKeyBytes, privateKey2Bytes)
}

func Test{{.Name}}NIKE(t *testing.T) {
	alicePrivate, alicePublic := Generate{{.Name}}KeyPair()
	bobPrivate, bobPublic := Generate{{.Name}}KeyPair()
	bobSharedBytes := DeriveSecret{{.Name}}(bobPrivate, alicePublic)
	aliceSharedBytes := DeriveSecret{{.Name}}(alicePrivate, bobPublic)
	require.Equal(t, bobSharedBytes, aliceSharedBytes)
}

`
