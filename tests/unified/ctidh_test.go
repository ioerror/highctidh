package unified

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

    "codeberg.org/vula/highctidh/unified/ctidh511"
    "codeberg.org/vula/highctidh/unified/ctidh512"
    "codeberg.org/vula/highctidh/unified/ctidh1024"
    "codeberg.org/vula/highctidh/unified/ctidh2048"
)

func TestNIKE511(t *testing.T) {
	ctidh := ctidh511.New()

	alicePrivate := ctidh.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh.DerivePublicKey(bobPrivate)

	ss1 := ctidh.DH(alicePrivate, bobPublic)
	ss2 := ctidh.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}

func TestNIKE512(t *testing.T) {
	ctidh := ctidh512.New()

	alicePrivate := ctidh.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh.DerivePublicKey(bobPrivate)

	ss1 := ctidh.DH(alicePrivate, bobPublic)
	ss2 := ctidh.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}

func TestNIKE1024(t *testing.T) {
	ctidh := ctidh1024.New()

	alicePrivate := ctidh.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh.DerivePublicKey(bobPrivate)

	ss1 := ctidh.DH(alicePrivate, bobPublic)
	ss2 := ctidh.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}

func TestNIKE2048(t *testing.T) {
	ctidh := ctidh2048.New()

	alicePrivate := ctidh.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh.DerivePublicKey(bobPrivate)

	ss1 := ctidh.DH(alicePrivate, bobPublic)
	ss2 := ctidh.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}
