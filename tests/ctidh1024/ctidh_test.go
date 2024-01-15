package ctidh1024

import (
	"crypto/rand"
	"testing"

	"github.com/stretchr/testify/require"

	"codeberg.org/vula/highctidh/unified/ctidh1024"
)

func TestNIKE(t *testing.T) {
	ctidh := ctidh1024.New()

	alicePrivate := ctidh.GenerateSecretKey(rand.Reader)
	alicePublic := ctidh.DerivePublicKey(alicePrivate)

	bobPrivate := ctidh.GenerateSecretKey(rand.Reader)
	bobPublic := ctidh.DerivePublicKey(bobPrivate)

	ss1 := ctidh.DH(alicePrivate, bobPublic)
	ss2 := ctidh.DH(bobPrivate, alicePublic)

	require.Equal(t, ss1.Bytes(), ss2.Bytes())
}