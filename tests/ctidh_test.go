package tests

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

    ctidh511 "codeberg.org/vula/highctidh/ctidh511"
    ctidh512 "codeberg.org/vula/highctidh/ctidh512"
    ctidh1024 "codeberg.org/vula/highctidh/ctidh1024"
    ctidh2048 "codeberg.org/vula/highctidh/ctidh2048"
)

func TestNIKE511(t *testing.T) {
	alicePrivate := ctidh511.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh511.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh511.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh511.DerivePublicKey(bobPrivate)

	ss1 := ctidh511.DH(alicePrivate, bobPublic)
	ss2 := ctidh511.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}

func TestNIKE512(t *testing.T) {
	alicePrivate := ctidh512.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh512.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh512.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh512.DerivePublicKey(bobPrivate)

	ss1 := ctidh512.DH(alicePrivate, bobPublic)
	ss2 := ctidh512.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}

func TestNIKE1024(t *testing.T) {
	alicePrivate := ctidh1024.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh1024.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh1024.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh1024.DerivePublicKey(bobPrivate)

	ss1 := ctidh1024.DH(alicePrivate, bobPublic)
	ss2 := ctidh1024.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}

func TestNIKE2048(t *testing.T) {
	alicePrivate := ctidh2048.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh2048.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh2048.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh2048.DerivePublicKey(bobPrivate)

	ss1 := ctidh2048.DH(alicePrivate, bobPublic)
	ss2 := ctidh2048.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}
