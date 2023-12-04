//go:build Ctidh2048
// +build Ctidh2048

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtidh2048BlindingOperation(t *testing.T) {
	mixPrivateKey, mixPublicKey := GenerateCtidh2048KeyPair()
	clientPrivateKey, clientPublicKey := GenerateCtidh2048KeyPair()

	blindingFactor := GenerateCtidh2048PrivateKey(rand.Reader)
	value1, err := BlindCtidh2048(blindingFactor, NewCtidh2048PublicKey(DeriveSecretCtidh2048(clientPrivateKey, mixPublicKey)))
	require.NoError(t, err)
	blinded, err := BlindCtidh2048(blindingFactor, clientPublicKey)
	require.NoError(t, err)
	value2 := DeriveSecretCtidh2048(mixPrivateKey, blinded)

	require.Equal(t, value1.Bytes(), value2)
}

func TestGenerateCtidh2048KeyPairWithRNG(t *testing.T) {
	privateKey, publicKey := GenerateCtidh2048KeyPairWithRNG(rand.Reader)
	zeros := make([]byte, Ctidh2048PublicKeySize)
	require.NotEqual(t, privateKey.Bytes(), zeros)
	require.NotEqual(t, publicKey.Bytes(), zeros)
}

func TestCtidh2048PublicKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh2048PublicKeySize)
	_, publicKey := GenerateCtidh2048KeyPair()
	require.NotEqual(t, publicKey.Bytes(), zeros)

	publicKey.Reset()
	require.Equal(t, publicKey.Bytes(), zeros)
}

func TestCtidh2048PrivateKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh2048PrivateKeySize)
	privateKey, _ := GenerateCtidh2048KeyPair()
	require.NotEqual(t, privateKey.Bytes(), zeros)

	privateKey.Reset()
	require.Equal(t, privateKey.Bytes(), zeros)
}

func TestCtidh2048PublicKeyMarshaling(t *testing.T) {
	privKey, publicKey := GenerateCtidh2048KeyPair()
	publicKeyBytes := publicKey.Bytes()

	publicKey2 := new(Ctidh2048PublicKey)
	err := publicKey2.FromBytes(publicKeyBytes)
	require.NoError(t, err)

	publicKey2Bytes := publicKey2.Bytes()

	publicKey3 := DeriveCtidh2048PublicKey(privKey)
	publicKey3Bytes := publicKey3.Bytes()

	require.Equal(t, publicKeyBytes, publicKey2Bytes)
	require.Equal(t, publicKey3Bytes, publicKeyBytes)
}

func TestCtidh2048PrivateKeyBytesing(t *testing.T) {
	privateKey, _ := GenerateCtidh2048KeyPair()
	privateKeyBytes := privateKey.Bytes()

	privateKey2 := new(Ctidh2048PrivateKey)
	privateKey2.FromBytes(privateKeyBytes)
	privateKey2Bytes := privateKey2.Bytes()

	require.Equal(t, privateKeyBytes, privateKey2Bytes)
}

func TestCtidh2048NIKE(t *testing.T) {
	alicePrivate, alicePublic := GenerateCtidh2048KeyPair()
	bobPrivate, bobPublic := GenerateCtidh2048KeyPair()
	bobSharedBytes := DeriveSecretCtidh2048(bobPrivate, alicePublic)
	aliceSharedBytes := DeriveSecretCtidh2048(alicePrivate, bobPublic)
	require.Equal(t, bobSharedBytes, aliceSharedBytes)
}
