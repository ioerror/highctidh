// DO NOT EDIT: generated code, see gen/main.go

package ctidh512

import (
	/*	"crypto/rand" */
	"testing"

	"github.com/stretchr/testify/require"
)

/*
func TestCtidh512BlindingOperation(t *testing.T) {
	mixPrivateKey, mixPublicKey := GenerateKeyPair()
	clientPrivateKey, clientPublicKey := GenerateKeyPair()

	blindingFactor := GeneratePrivateKey(rand.Reader)
	value1, err := BlindCtidh512(blindingFactor, NewPublicKey(DeriveSecret(clientPrivateKey, mixPublicKey)))
	require.NoError(t, err)
	blinded, err := BlindCtidh512(blindingFactor, clientPublicKey)
	require.NoError(t, err)
	value2 := DeriveSecret(mixPrivateKey, blinded)

	require.Equal(t, value1.Bytes(), value2)
}

func TestGenerateKeyPairWithRNG(t *testing.T) {
	privateKey, publicKey := GenerateKeyPairWithRNG(rand.Reader)
	zeros := make([]byte, PublicKeySize)
	require.NotEqual(t, privateKey.Bytes(), zeros)
	require.NotEqual(t, publicKey.Bytes(), zeros)
}
*/

func TestPublicKeyReset(t *testing.T) {
	zeros := make([]byte, PublicKeySize)
	_, publicKey := GenerateKeyPair()
	require.NotEqual(t, publicKey.Bytes(), zeros)

	publicKey.Reset()
	require.Equal(t, publicKey.Bytes(), zeros)
}

func TestPrivateKeyReset(t *testing.T) {
	zeros := make([]byte, PrivateKeySize)
	privateKey, _ := GenerateKeyPair()
	require.NotEqual(t, privateKey.Bytes(), zeros)

	privateKey.Reset()
	require.Equal(t, privateKey.Bytes(), zeros)
}

func TestPublicKeyMarshaling(t *testing.T) {
	privKey, publicKey := GenerateKeyPair()
	publicKeyBytes := publicKey.Bytes()

	publicKey2 := new(PublicKey)
	err := publicKey2.FromBytes(publicKeyBytes)
	require.NoError(t, err)

	publicKey2Bytes := publicKey2.Bytes()

	publicKey3 := DerivePublicKey(privKey)
	publicKey3Bytes := publicKey3.Bytes()

	require.Equal(t, publicKeyBytes, publicKey2Bytes)
	require.Equal(t, publicKey3Bytes, publicKeyBytes)
}

func TestPrivateKeyBytesing(t *testing.T) {
	privateKey, _ := GenerateKeyPair()
	privateKeyBytes := privateKey.Bytes()

	privateKey2 := new(PrivateKey)
	privateKey2.FromBytes(privateKeyBytes)
	privateKey2Bytes := privateKey2.Bytes()

	require.Equal(t, privateKeyBytes, privateKey2Bytes)
}

func TestCtidh512NIKE(t *testing.T) {
	alicePrivate, alicePublic := GenerateKeyPair()
	bobPrivate, bobPublic := GenerateKeyPair()
	bobSharedBytes := DeriveSecret(bobPrivate, alicePublic)
	aliceSharedBytes := DeriveSecret(alicePrivate, bobPublic)
	require.Equal(t, bobSharedBytes, aliceSharedBytes)
}
