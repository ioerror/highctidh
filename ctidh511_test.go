//go:build Ctidh511
// +build Ctidh511

// DO NOT EDIT: generated code, see gen/main.go

package ctidh

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCtidh511BlindingOperation(t *testing.T) {
	mixPrivateKey, mixPublicKey := GenerateCtidh511KeyPair()
	clientPrivateKey, clientPublicKey := GenerateCtidh511KeyPair()

	blindingFactor := GenerateCtidh511PrivateKey(rand.Reader)
	value1, err := BlindCtidh511(blindingFactor, NewCtidh511PublicKey(DeriveSecretCtidh511(clientPrivateKey, mixPublicKey)))
	require.NoError(t, err)
	blinded, err := BlindCtidh511(blindingFactor, clientPublicKey)
	require.NoError(t, err)
	value2 := DeriveSecretCtidh511(mixPrivateKey, blinded)

	require.Equal(t, value1.Bytes(), value2)
}

func TestGenerateCtidh511KeyPairWithRNG(t *testing.T) {
	privateKey, publicKey := GenerateCtidh511KeyPairWithRNG(rand.Reader)
	zeros := make([]byte, Ctidh511PublicKeySize)
	require.NotEqual(t, privateKey.Bytes(), zeros)
	require.NotEqual(t, publicKey.Bytes(), zeros)
}

func TestCtidh511PublicKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh511PublicKeySize)
	_, publicKey := GenerateCtidh511KeyPair()
	require.NotEqual(t, publicKey.Bytes(), zeros)

	publicKey.Reset()
	require.Equal(t, publicKey.Bytes(), zeros)
}

func TestCtidh511PrivateKeyReset(t *testing.T) {
	zeros := make([]byte, Ctidh511PrivateKeySize)
	privateKey, _ := GenerateCtidh511KeyPair()
	require.NotEqual(t, privateKey.Bytes(), zeros)

	privateKey.Reset()
	require.Equal(t, privateKey.Bytes(), zeros)
}

func TestCtidh511PublicKeyMarshaling(t *testing.T) {
	privKey, publicKey := GenerateCtidh511KeyPair()
	publicKeyBytes := publicKey.Bytes()

	publicKey2 := new(Ctidh511PublicKey)
	err := publicKey2.FromBytes(publicKeyBytes)
	require.NoError(t, err)

	publicKey2Bytes := publicKey2.Bytes()

	publicKey3 := DeriveCtidh511PublicKey(privKey)
	publicKey3Bytes := publicKey3.Bytes()

	require.Equal(t, publicKeyBytes, publicKey2Bytes)
	require.Equal(t, publicKey3Bytes, publicKeyBytes)
}

func TestCtidh511PrivateKeyBytesing(t *testing.T) {
	privateKey, _ := GenerateCtidh511KeyPair()
	privateKeyBytes := privateKey.Bytes()

	privateKey2 := new(Ctidh511PrivateKey)
	privateKey2.FromBytes(privateKeyBytes)
	privateKey2Bytes := privateKey2.Bytes()

	require.Equal(t, privateKeyBytes, privateKey2Bytes)
}

func TestCtidh511NIKE(t *testing.T) {
	alicePrivate, alicePublic := GenerateCtidh511KeyPair()
	bobPrivate, bobPublic := GenerateCtidh511KeyPair()
	bobSharedBytes := DeriveSecretCtidh511(bobPrivate, alicePublic)
	aliceSharedBytes := DeriveSecretCtidh511(alicePrivate, bobPublic)
	require.Equal(t, bobSharedBytes, aliceSharedBytes)
}
