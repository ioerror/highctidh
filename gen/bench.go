package main

var BenchmarkTests = `
// +build {{.Name}}
// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkPublicKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privKey, publicKey := Generate{{.Name}}KeyPair()

		publicKeyBytes := publicKey.Bytes()

		publicKey2 := new({{.Name}}PublicKey)
		err := publicKey2.FromBytes(publicKeyBytes)
		require.NoError(b, err)

		publicKey2Bytes := publicKey2.Bytes()
		publicKey3 := Derive{{.Name}}PublicKey(privKey)
		publicKey3Bytes := publicKey3.Bytes()

		require.Equal(b, publicKeyBytes, publicKey2Bytes)
		require.Equal(b, publicKey3Bytes, publicKeyBytes)
	}
}

func BenchmarkPrivateKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privateKey, _ := Generate{{.Name}}KeyPair()
		privateKeyBytes := privateKey.Bytes()

		privateKey2 := new({{.Name}}PrivateKey)
		privateKey2.FromBytes(privateKeyBytes)
		privateKey2Bytes := privateKey2.Bytes()

		require.Equal(b, privateKeyBytes, privateKey2Bytes)
	}
}

func BenchmarkNIKE(b *testing.B) {
	for n := 0; n < b.N; n++ {
		alicePrivate, alicePublic := Generate{{.Name}}KeyPair()
		bobPrivate, bobPublic := Generate{{.Name}}KeyPair()

		bobSharedBytes := DeriveSecret{{.Name}}(bobPrivate, alicePublic)
		aliceSharedBytes := DeriveSecret{{.Name}}(alicePrivate, bobPublic)

		require.Equal(b, bobSharedBytes, aliceSharedBytes)
	}
}

func BenchmarkDeriveSecret(b *testing.B) {
	alicePrivate, alicePublic := Generate{{.Name}}KeyPair()
	bobPrivate, bobPublic := Generate{{.Name}}KeyPair()

	var aliceSharedBytes []byte
	for n := 0; n < b.N; n++ {
		aliceSharedBytes = DeriveSecret{{.Name}}(alicePrivate, bobPublic)
	}

	bobSharedBytes := DeriveSecret{{.Name}}(bobPrivate, alicePublic)
	require.Equal(b, bobSharedBytes, aliceSharedBytes)
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = Generate{{.Name}}KeyPair()
	}
}

`
