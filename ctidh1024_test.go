//go:build Ctidh1024
// +build Ctidh1024

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtidh1024BlindingOperation(t *testing.T) {
	mixPrivateKey, mixPublicKey := GenerateCtidh1024KeyPair()
	clientPrivateKey, clientPublicKey := GenerateCtidh1024KeyPair()

	blindingFactor := GenerateCtidh1024PrivateKey(rand.Reader)
	value1, err := BlindCtidh1024(blindingFactor, NewCtidh1024PublicKey(DeriveSecretCtidh1024(clientPrivateKey, mixPublicKey)))
	require.NoError(t, err)
	blinded, err := BlindCtidh1024(blindingFactor, clientPublicKey)
	require.NoError(t, err)
	value2 := DeriveSecretCtidh1024(mixPrivateKey, blinded)

	require.Equal(t, value1.Bytes(), value2)
}

func TestGenerateCtidh1024KeyPairWithRNG(t *testing.T) {
	privateKey, publicKey := GenerateCtidh1024KeyPairWithRNG(rand.Reader)
	zeros := make([]byte, Ctidh1024PublicKeySize)
	require.NotEqual(t, privateKey.Bytes(), zeros)
	require.NotEqual(t, publicKey.Bytes(), zeros)
}

func TestCtidh1024PublicKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh1024PublicKeySize)
	_, publicKey := GenerateCtidh1024KeyPair()
	require.NotEqual(t, publicKey.Bytes(), zeros)

	publicKey.Reset()
	require.Equal(t, publicKey.Bytes(), zeros)
}

func TestCtidh1024PrivateKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh1024PrivateKeySize)
	privateKey, _ := GenerateCtidh1024KeyPair()
	require.NotEqual(t, privateKey.Bytes(), zeros)

	privateKey.Reset()
	require.Equal(t, privateKey.Bytes(), zeros)
}

func TestCtidh1024PublicKeyMarshaling(t *testing.T) {
	privKey, publicKey := GenerateCtidh1024KeyPair()
	publicKeyBytes := publicKey.Bytes()

	publicKey2 := new(Ctidh1024PublicKey)
	err := publicKey2.FromBytes(publicKeyBytes)
	require.NoError(t, err)

	publicKey2Bytes := publicKey2.Bytes()

	publicKey3 := DeriveCtidh1024PublicKey(privKey)
	publicKey3Bytes := publicKey3.Bytes()

	require.Equal(t, publicKeyBytes, publicKey2Bytes)
	require.Equal(t, publicKey3Bytes, publicKeyBytes)
}

func TestCtidh1024PrivateKeyBytesing(t *testing.T) {
	privateKey, _ := GenerateCtidh1024KeyPair()
	privateKeyBytes := privateKey.Bytes()

	privateKey2 := new(Ctidh1024PrivateKey)
	privateKey2.FromBytes(privateKeyBytes)
	privateKey2Bytes := privateKey2.Bytes()

	require.Equal(t, privateKeyBytes, privateKey2Bytes)
}

func TestCtidh1024NIKE(t *testing.T) {
	alicePrivate, alicePublic := GenerateCtidh1024KeyPair()
	bobPrivate, bobPublic := GenerateCtidh1024KeyPair()
	bobSharedBytes := DeriveSecretCtidh1024(bobPrivate, alicePublic)
	aliceSharedBytes := DeriveSecretCtidh1024(alicePrivate, bobPublic)
	require.Equal(t, bobSharedBytes, aliceSharedBytes)
}
