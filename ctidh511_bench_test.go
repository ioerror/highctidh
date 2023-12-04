//go:build Ctidh511
// +build Ctidh511

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkPublicKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privKey, publicKey := GenerateCtidh511KeyPair()

		publicKeyBytes := publicKey.Bytes()

		publicKey2 := new(Ctidh511PublicKey)
		err := publicKey2.FromBytes(publicKeyBytes)
		require.NoError(b, err)

		publicKey2Bytes := publicKey2.Bytes()
		publicKey3 := DeriveCtidh511PublicKey(privKey)
		publicKey3Bytes := publicKey3.Bytes()

		require.Equal(b, publicKeyBytes, publicKey2Bytes)
		require.Equal(b, publicKey3Bytes, publicKeyBytes)
	}
}

func BenchmarkPrivateKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privateKey, _ := GenerateCtidh511KeyPair()
		privateKeyBytes := privateKey.Bytes()

		privateKey2 := new(Ctidh511PrivateKey)
		privateKey2.FromBytes(privateKeyBytes)
		privateKey2Bytes := privateKey2.Bytes()

		require.Equal(b, privateKeyBytes, privateKey2Bytes)
	}
}

func BenchmarkNIKE(b *testing.B) {
	for n := 0; n < b.N; n++ {
		alicePrivate, alicePublic := GenerateCtidh511KeyPair()
		bobPrivate, bobPublic := GenerateCtidh511KeyPair()

		bobSharedBytes := DeriveSecretCtidh511(bobPrivate, alicePublic)
		aliceSharedBytes := DeriveSecretCtidh511(alicePrivate, bobPublic)

		require.Equal(b, bobSharedBytes, aliceSharedBytes)
	}
}

func BenchmarkDeriveSecret(b *testing.B) {
	alicePrivate, alicePublic := GenerateCtidh511KeyPair()
	bobPrivate, bobPublic := GenerateCtidh511KeyPair()

	var aliceSharedBytes []byte
	for n := 0; n < b.N; n++ {
		aliceSharedBytes = DeriveSecretCtidh511(alicePrivate, bobPublic)
	}

	bobSharedBytes := DeriveSecretCtidh511(bobPrivate, alicePublic)
	require.Equal(b, bobSharedBytes, aliceSharedBytes)
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = GenerateCtidh511KeyPair()
	}
}
