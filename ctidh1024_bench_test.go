//go:build Ctidh1024
// +build Ctidh1024

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkPublicKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privKey, publicKey := GenerateCtidh1024KeyPair()

		publicKeyBytes := publicKey.Bytes()

		publicKey2 := new(Ctidh1024PublicKey)
		err := publicKey2.FromBytes(publicKeyBytes)
		require.NoError(b, err)

		publicKey2Bytes := publicKey2.Bytes()
		publicKey3 := DeriveCtidh1024PublicKey(privKey)
		publicKey3Bytes := publicKey3.Bytes()

		require.Equal(b, publicKeyBytes, publicKey2Bytes)
		require.Equal(b, publicKey3Bytes, publicKeyBytes)
	}
}

func BenchmarkPrivateKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privateKey, _ := GenerateCtidh1024KeyPair()
		privateKeyBytes := privateKey.Bytes()

		privateKey2 := new(Ctidh1024PrivateKey)
		privateKey2.FromBytes(privateKeyBytes)
		privateKey2Bytes := privateKey2.Bytes()

		require.Equal(b, privateKeyBytes, privateKey2Bytes)
	}
}

func BenchmarkNIKE(b *testing.B) {
	for n := 0; n < b.N; n++ {
		alicePrivate, alicePublic := GenerateCtidh1024KeyPair()
		bobPrivate, bobPublic := GenerateCtidh1024KeyPair()

		bobSharedBytes := DeriveSecretCtidh1024(bobPrivate, alicePublic)
		aliceSharedBytes := DeriveSecretCtidh1024(alicePrivate, bobPublic)

		require.Equal(b, bobSharedBytes, aliceSharedBytes)
	}
}

func BenchmarkDeriveSecret(b *testing.B) {
	alicePrivate, alicePublic := GenerateCtidh1024KeyPair()
	bobPrivate, bobPublic := GenerateCtidh1024KeyPair()

	var aliceSharedBytes []byte
	for n := 0; n < b.N; n++ {
		aliceSharedBytes = DeriveSecretCtidh1024(alicePrivate, bobPublic)
	}

	bobSharedBytes := DeriveSecretCtidh1024(bobPrivate, alicePublic)
	require.Equal(b, bobSharedBytes, aliceSharedBytes)
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = GenerateCtidh1024KeyPair()
	}
}
