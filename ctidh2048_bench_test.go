//go:build Ctidh2048
// +build Ctidh2048

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func BenchmarkPublicKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privKey, publicKey := GenerateCtidh2048KeyPair()

		publicKeyBytes := publicKey.Bytes()

		publicKey2 := new(Ctidh2048PublicKey)
		err := publicKey2.FromBytes(publicKeyBytes)
		require.NoError(b, err)

		publicKey2Bytes := publicKey2.Bytes()
		publicKey3 := DeriveCtidh2048PublicKey(privKey)
		publicKey3Bytes := publicKey3.Bytes()

		require.Equal(b, publicKeyBytes, publicKey2Bytes)
		require.Equal(b, publicKey3Bytes, publicKeyBytes)
	}
}

func BenchmarkPrivateKeySerializing(b *testing.B) {
	for n := 0; n < b.N; n++ {
		privateKey, _ := GenerateCtidh2048KeyPair()
		privateKeyBytes := privateKey.Bytes()

		privateKey2 := new(Ctidh2048PrivateKey)
		privateKey2.FromBytes(privateKeyBytes)
		privateKey2Bytes := privateKey2.Bytes()

		require.Equal(b, privateKeyBytes, privateKey2Bytes)
	}
}

func BenchmarkNIKE(b *testing.B) {
	for n := 0; n < b.N; n++ {
		alicePrivate, alicePublic := GenerateCtidh2048KeyPair()
		bobPrivate, bobPublic := GenerateCtidh2048KeyPair()

		bobSharedBytes := DeriveSecretCtidh2048(bobPrivate, alicePublic)
		aliceSharedBytes := DeriveSecretCtidh2048(alicePrivate, bobPublic)

		require.Equal(b, bobSharedBytes, aliceSharedBytes)
	}
}

func BenchmarkDeriveSecret(b *testing.B) {
	alicePrivate, alicePublic := GenerateCtidh2048KeyPair()
	bobPrivate, bobPublic := GenerateCtidh2048KeyPair()

	var aliceSharedBytes []byte
	for n := 0; n < b.N; n++ {
		aliceSharedBytes = DeriveSecretCtidh2048(alicePrivate, bobPublic)
	}

	bobSharedBytes := DeriveSecretCtidh2048(bobPrivate, alicePublic)
	require.Equal(b, bobSharedBytes, aliceSharedBytes)
}

func BenchmarkGenerateKeyPair(b *testing.B) {
	for n := 0; n < b.N; n++ {
		_, _ = GenerateCtidh2048KeyPair()
	}
}
