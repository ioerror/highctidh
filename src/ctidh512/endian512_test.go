package ctidh512

import (
	"encoding/binary"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func rawLimbs(k *PublicKey) []uint64 {
	return unsafe.Slice((*uint64)(unsafe.Pointer(&k.publicKey)), PublicKeySize/8)
}

func setRawPublicKey(k *PublicKey, b []byte) {
	limbs := rawLimbs(k)
	for i := range limbs {
		limbs[i] = binary.LittleEndian.Uint64(b[8*i:])
	}
}

func requireLittleEndian(t *testing.T, k *PublicKey, b []byte) {
	t.Helper()
	require.Len(t, b, PublicKeySize)
	for i, limb := range rawLimbs(k) {
		require.Equal(t, binary.LittleEndian.Uint64(b[8*i:]), limb, "limb %d", i)
	}
}

func TestPublicKeyLittleEndian(t *testing.T) {
	priv, pub := GenerateKeyPair()
	requireLittleEndian(t, pub, pub.Bytes())
	bin, err := pub.MarshalBinary()
	require.NoError(t, err)
	requireLittleEndian(t, pub, bin)

	loaded := new(PublicKey)
	require.NoError(t, loaded.FromBytes(pub.Bytes()))
	require.Equal(t, rawLimbs(pub), rawLimbs(loaded))
	loaded = new(PublicKey)
	require.NoError(t, loaded.UnmarshalBinary(bin))
	require.Equal(t, rawLimbs(pub), rawLimbs(loaded))
	text, err := pub.MarshalText()
	require.NoError(t, err)
	loaded = new(PublicKey)
	require.NoError(t, loaded.UnmarshalText(text))
	require.Equal(t, rawLimbs(pub), rawLimbs(loaded))

	priv2, pub2 := GenerateKeyPair()
	shared, err := GroupActionChecked(priv, pub2)
	require.NoError(t, err)
	secret, err := DeriveSecretChecked(priv, pub2)
	require.NoError(t, err)
	requireLittleEndian(t, shared, secret)
	requireLittleEndian(t, shared, DeriveSecret(priv2, pub))
	requireLittleEndian(t, shared, priv.DeriveSecret(pub2))
}
