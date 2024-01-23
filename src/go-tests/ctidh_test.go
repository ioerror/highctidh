package tests

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

    ctidh511 "codeberg.org/vula/highctidh/src/ctidh511"
    ctidh512 "codeberg.org/vula/highctidh/src/ctidh512"
    ctidh1024 "codeberg.org/vula/highctidh/src/ctidh1024"
    ctidh2048 "codeberg.org/vula/highctidh/src/ctidh2048"
)

func TestNIKE511(t *testing.T) {
	alicePrivate := ctidh511.GeneratePrivateKey(rand.Reader)
	alicePublic := ctidh511.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh511.GeneratePrivateKey(rand.Reader)
	bobPublic := ctidh511.DerivePublicKey(bobPrivate)

	ss1 := ctidh511.DeriveSecret(alicePrivate, bobPublic)
	ss2 := ctidh511.DeriveSecret(bobPrivate, alicePublic)

	require.Equal(t, ss1, ss2)
}

func TestNIKE512(t *testing.T) {
	alicePrivate := ctidh512.GeneratePrivateKey(rand.Reader)
	alicePublic := ctidh512.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh512.GeneratePrivateKey(rand.Reader)
	bobPublic := ctidh512.DerivePublicKey(bobPrivate)

	ss1 := ctidh512.DeriveSecret(alicePrivate, bobPublic)
	ss2 := ctidh512.DeriveSecret(bobPrivate, alicePublic)

	require.Equal(t, ss1, ss2)
}

func TestNIKE1024(t *testing.T) {
	alicePrivate := ctidh1024.GeneratePrivateKey(rand.Reader)
	alicePublic := ctidh1024.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh1024.GeneratePrivateKey(rand.Reader)
	bobPublic := ctidh1024.DerivePublicKey(bobPrivate)

	ss1 := ctidh1024.DeriveSecret(alicePrivate, bobPublic)
	ss2 := ctidh1024.DeriveSecret(bobPrivate, alicePublic)

	require.Equal(t, ss1, ss2)
}

func TestNIKE2048(t *testing.T) {
	alicePrivate := ctidh2048.GeneratePrivateKey(rand.Reader)
	alicePublic := ctidh2048.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh2048.GeneratePrivateKey(rand.Reader)
	bobPublic := ctidh2048.DerivePublicKey(bobPrivate)

	ss1 := ctidh2048.DeriveSecret(alicePrivate, bobPublic)
	ss2 := ctidh2048.DeriveSecret(bobPrivate, alicePublic)

	require.Equal(t, ss1, ss2)
}
