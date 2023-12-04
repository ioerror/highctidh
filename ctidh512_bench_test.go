//go:build Ctidh512
// +build Ctidh512

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkPublicKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privKey, publicKey := GenerateCtidh512KeyPair()

		publicKeyBytes := publicKey.Bytes()

		publicKey2 := new(Ctidh512PublicKey)
		err := publicKey2.FromBytes(publicKeyBytes)
		require.NoError(b, err)

		publicKey2Bytes := publicKey2.Bytes()
		publicKey3 := DeriveCtidh512PublicKey(privKey)
		publicKey3Bytes := publicKey3.Bytes()

		require.Equal(b, publicKeyBytes, publicKey2Bytes)
		require.Equal(b, publicKey3Bytes, publicKeyBytes)
	}
}

func BenchmarkPrivateKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privateKey, _ := GenerateCtidh512KeyPair()
		privateKeyBytes := privateKey.Bytes()

		privateKey2 := new(Ctidh512PrivateKey)
		privateKey2.FromBytes(privateKeyBytes)
		privateKey2Bytes := privateKey2.Bytes()

		require.Equal(b, privateKeyBytes, privateKey2Bytes)
	}
}

func BenchmarkNIKE(b *testing.B) {
	for n := 0; n < b.N; n++ {
		alicePrivate, alicePublic := GenerateCtidh512KeyPair()
		bobPrivate, bobPublic := GenerateCtidh512KeyPair()

		bobSharedBytes := DeriveSecretCtidh512(bobPrivate, alicePublic)
		aliceSharedBytes := DeriveSecretCtidh512(alicePrivate, bobPublic)

		require.Equal(b, bobSharedBytes, aliceSharedBytes)
	}
}

func BenchmarkDeriveSecret(b *testing.B) {
	alicePrivate, alicePublic := GenerateCtidh512KeyPair()
	bobPrivate, bobPublic := GenerateCtidh512KeyPair()

	var aliceSharedBytes []byte
	for n := 0; n < b.N; n++ {
		aliceSharedBytes = DeriveSecretCtidh512(alicePrivate, bobPublic)
	}

	bobSharedBytes := DeriveSecretCtidh512(bobPrivate, alicePublic)
	require.Equal(b, bobSharedBytes, aliceSharedBytes)
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = GenerateCtidh512KeyPair()
	}
}
