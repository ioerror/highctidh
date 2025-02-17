// DO NOT EDIT: generated code, see gen/main.go

package ctidh511

import (
	"crypto/rand"
	"github.com/stretchr/testify/require"
	"sync"
	"testing"
)

func TestBlindingOperation(t *testing.T) {
	t.Logf("out: a")
	mixPrivateKey, mixPublicKey := GenerateKeyPair()
	t.Logf("out: b")
	clientPrivateKey, clientPublicKey := GenerateKeyPair()

	t.Logf("out: c")
	blindingFactor := GeneratePrivateKey(rand.Reader)
	t.Logf("out: d")
	value1, err := Blind(blindingFactor, NewPublicKey(DeriveSecret(clientPrivateKey, mixPublicKey)))
	t.Logf("out: e")
	require.NoError(t, err)
	blinded, err := Blind(blindingFactor, clientPublicKey)
	t.Logf("out: f")
	require.NoError(t, err)
	value2 := DeriveSecret(mixPrivateKey, blinded)
	t.Logf("out: g")

	require.Equal(t, value1.Bytes(), value2)
}

func TestBlindingOperationNoRandReader(t *testing.T) {
	mixPrivateKey, mixPublicKey := GenerateKeyPair()
	clientPrivateKey, clientPublicKey := GenerateKeyPair()

	blindingFactor, _ := GenerateKeyPair()
	value1, err := Blind(blindingFactor, NewPublicKey(DeriveSecret(clientPrivateKey, mixPublicKey)))
	require.NoError(t, err)
	blinded, err := Blind(blindingFactor, clientPublicKey)
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

func TestGenerateKeyPair(t *testing.T) {
	for i := 0; i < 16; i++ {
		privateKey, publicKey := GenerateKeyPair()
		zeros := make([]byte, PublicKeySize)
		require.NotEqual(t, privateKey.Bytes(), zeros)
		require.NotEqual(t, publicKey.Bytes(), zeros)
	}
}

func TestCorruptStack(t *testing.T) {
	errCh := make(chan error, 10)
	wg := new(sync.WaitGroup)
	for i := 0; i < 10; i++ {
		wg.Add(1)
		go func() {
			go func() {
				foo := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
				t.Logf("stack: %s", foo)
			}()

			_, _ = GenerateKeyPair()

			go func() {
				foo := []byte("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa")
				t.Logf("stack: %s", foo)
			}()
			errCh <- nil
			wg.Done()
		}()
	}
	wg.Wait()
	close(errCh)

	for e := range errCh {
		require.NoError(t, e)
	}
	wg.Add(1)
	go func() {
		t.Log("last call")
		wg.Done()
	}()
	wg.Wait()
}

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

func TestPrivateKeyByteMarshaling(t *testing.T) {
	privateKey, _ := GenerateKeyPair()
	privateKeyBytes := privateKey.Bytes()

	privateKey2 := new(PrivateKey)
	privateKey2.FromBytes(privateKeyBytes)
	privateKey2Bytes := privateKey2.Bytes()

	require.Equal(t, privateKeyBytes, privateKey2Bytes)
}

func TestNIKE(t *testing.T) {
	alicePrivate, alicePublic := GenerateKeyPair()
	bobPrivate, bobPublic := GenerateKeyPair()
	bobSharedBytes := DeriveSecret(bobPrivate, alicePublic)
	aliceSharedBytes := DeriveSecret(alicePrivate, bobPublic)
	require.Equal(t, bobSharedBytes, aliceSharedBytes)
}
