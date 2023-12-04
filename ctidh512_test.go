//go:build Ctidh512
// +build Ctidh512

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtidh512BlindingOperation(t *testing.T) {
	mixPrivateKey, mixPublicKey := GenerateCtidh512KeyPair()
	clientPrivateKey, clientPublicKey := GenerateCtidh512KeyPair()

	blindingFactor := GenerateCtidh512PrivateKey(rand.Reader)
	value1, err := BlindCtidh512(blindingFactor, NewCtidh512PublicKey(DeriveSecretCtidh512(clientPrivateKey, mixPublicKey)))
	require.NoError(t, err)
	blinded, err := BlindCtidh512(blindingFactor, clientPublicKey)
	require.NoError(t, err)
	value2 := DeriveSecretCtidh512(mixPrivateKey, blinded)

	require.Equal(t, value1.Bytes(), value2)
}

func TestGenerateCtidh512KeyPairWithRNG(t *testing.T) {
	privateKey, publicKey := GenerateCtidh512KeyPairWithRNG(rand.Reader)
	zeros := make([]byte, Ctidh512PublicKeySize)
	require.NotEqual(t, privateKey.Bytes(), zeros)
	require.NotEqual(t, publicKey.Bytes(), zeros)
}

func TestCtidh512PublicKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh512PublicKeySize)
	_, publicKey := GenerateCtidh512KeyPair()
	require.NotEqual(t, publicKey.Bytes(), zeros)

	publicKey.Reset()
	require.Equal(t, publicKey.Bytes(), zeros)
}

func TestCtidh512PrivateKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh512PrivateKeySize)
	privateKey, _ := GenerateCtidh512KeyPair()
	require.NotEqual(t, privateKey.Bytes(), zeros)

	privateKey.Reset()
	require.Equal(t, privateKey.Bytes(), zeros)
}

func TestCtidh512PublicKeyMarshaling(t *testing.T) {
	privKey, publicKey := GenerateCtidh512KeyPair()
	publicKeyBytes := publicKey.Bytes()

	publicKey2 := new(Ctidh512PublicKey)
	err := publicKey2.FromBytes(publicKeyBytes)
	require.NoError(t, err)

	publicKey2Bytes := publicKey2.Bytes()

	publicKey3 := DeriveCtidh512PublicKey(privKey)
	publicKey3Bytes := publicKey3.Bytes()

	require.Equal(t, publicKeyBytes, publicKey2Bytes)
	require.Equal(t, publicKey3Bytes, publicKeyBytes)
}

func TestCtidh512PrivateKeyBytesing(t *testing.T) {
	privateKey, _ := GenerateCtidh512KeyPair()
	privateKeyBytes := privateKey.Bytes()

	privateKey2 := new(Ctidh512PrivateKey)
	privateKey2.FromBytes(privateKeyBytes)
	privateKey2Bytes := privateKey2.Bytes()

	require.Equal(t, privateKeyBytes, privateKey2Bytes)
}

func TestCtidh512NIKE(t *testing.T) {
	alicePrivate, alicePublic := GenerateCtidh512KeyPair()
	bobPrivate, bobPublic := GenerateCtidh512KeyPair()
	bobSharedBytes := DeriveSecretCtidh512(bobPrivate, alicePublic)
	aliceSharedBytes := DeriveSecretCtidh512(alicePrivate, bobPublic)
	require.Equal(t, bobSharedBytes, aliceSharedBytes)
}
