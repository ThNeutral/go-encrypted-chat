package ecdh_test

import (
	"chat/shared/ecdh"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestECDH(t *testing.T) {
	var alice, bob ecdh.ECDH
	var err error

	t.Run("should successfully generate ECDH instance", func(t *testing.T) {
		alice, err = ecdh.NewX25519()
		require.NoError(t, err)

		bob, err = ecdh.NewX25519()
		require.NoError(t, err)
	})

	alicePub := alice.PublicKey()
	bobPub := bob.PublicKey()

	t.Run("GenerateSharedSecret should fail if secret was not yet generated", func(t *testing.T) {
		_, err = alice.SharedSecret()
		require.Error(t, err)

		_, err = bob.SharedSecret()
		require.Error(t, err)
	})

	t.Run("should successfully go through ECDH process", func(t *testing.T) {
		err = alice.GenerateSharedSecret(bobPub)
		require.NoError(t, err)

		err = bob.GenerateSharedSecret(alicePub)
		require.NoError(t, err)

		aliceShared, err := alice.SharedSecret()
		require.NoError(t, err)

		bobShared, err := bob.SharedSecret()
		require.NoError(t, err)

		require.Equal(t, aliceShared, bobShared)
	})

	t.Run("GenerateSharedSecret should fail if secret was already generated", func(t *testing.T) {
		err = alice.GenerateSharedSecret(bobPub)
		require.Error(t, err)

		err = bob.GenerateSharedSecret(alicePub)
		require.Error(t, err)
	})
}
