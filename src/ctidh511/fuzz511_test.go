package ctidh511

import (
	"bytes"
	"testing"
	"unsafe"

	"github.com/stretchr/testify/require"
)

func FuzzKeys(f *testing.F) {
	priv, pub := GenerateKeyPairWithRNG(&testRNG{limit: -1})
	for op := uint8(0); op < 4; op++ {
		f.Add(pub.Bytes(), priv.Bytes(), op)
	}
	f.Add(make([]byte, PublicKeySize), make([]byte, PrivateKeySize), uint8(0))
	f.Add(bytes.Repeat([]byte{0xff}, PublicKeySize), bytes.Repeat([]byte{0x7f}, PrivateKeySize), uint8(1))
	f.Add([]byte{7}, []byte{1}, uint8(2))
	f.Fuzz(func(t *testing.T, pubBytes, privBytes []byte, op uint8) {
		pk := NewEmptyPublicKey()
		pubErr := pk.FromBytes(pubBytes)
		sk := NewEmptyPrivateKey()
		privErr := sk.FromBytes(privBytes)
		_ = new(PublicKey).UnmarshalText(pubBytes)
		_ = new(PrivateKey).UnmarshalText(privBytes)

		k, err := NewPublicKeyChecked(pubBytes)
		require.Equal(t, pubErr, err)
		if pubErr != nil {
			require.Nil(t, k)
			require.Equal(t, make([]byte, PublicKeySize), pk.Bytes())
			if len(pubBytes) == PublicKeySize {
				raw := new(PublicKey)
				copy(unsafe.Slice((*byte)(unsafe.Pointer(&raw.publicKey)), PublicKeySize), pubBytes)
				shared, err := GroupActionChecked(sk, raw)
				require.ErrorIs(t, err, ErrCTIDH)
				require.Nil(t, shared)
			}
		} else {
			require.Equal(t, pubBytes, pk.Bytes())
			require.True(t, pk.Equal(k))
			text, err := pk.MarshalText()
			require.NoError(t, err)
			k2 := new(PublicKey)
			require.NoError(t, k2.UnmarshalText(text))
			require.Equal(t, pubBytes, k2.Bytes())
		}
		if privErr != nil {
			require.Equal(t, make([]byte, PrivateKeySize), sk.Bytes())
		} else {
			require.Equal(t, privBytes, sk.Bytes())
			text, err := sk.MarshalText()
			require.NoError(t, err)
			k2 := new(PrivateKey)
			require.NoError(t, k2.UnmarshalText(text))
			require.True(t, sk.Equal(k2))
		}
		if pubErr != nil || privErr != nil {
			return
		}

		switch op % 4 {
		case 0:
			secret, err := DeriveSecretChecked(sk, pk)
			require.NoError(t, err)
			require.Len(t, secret, PublicKeySize)
		case 1:
			shared, err := GroupActionChecked(sk, pk)
			require.NoError(t, err)
			require.True(t, shared.validated)
		case 2:
			blinded, err := Blind(sk, pk)
			require.NoError(t, err)
			require.True(t, blinded.validated)
		case 3:
			require.NoError(t, pk.Blind(sk))
			require.True(t, pk.validated)
		}
	})
}
