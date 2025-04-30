package aes_test

import (
	"chat/shared/aes"
	"crypto/rand"
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
)

func generateAESKey(size int) ([]byte, error) {
	key := make([]byte, size)
	_, err := rand.Read(key)
	if err != nil {
		return nil, err
	}
	return key, nil
}

func TestEncryption(t *testing.T) {
	t.Run("should successfully encrypt and decrypt plaintext", func(t *testing.T) {
		keyValid, err := generateAESKey(32)
		require.NoError(t, err)

		plaintext := []byte(uuid.New().String())

		cyphertext, err := aes.Encrypt(plaintext, keyValid)
		require.NoError(t, err)

		plaintext2, err := aes.Decrypt(cyphertext, keyValid)
		require.NoError(t, err)

		require.Equal(t, plaintext, plaintext2)
	})

	t.Run("should fail with invalid key size", func(t *testing.T) {
		keyValid, err := generateAESKey(32)
		require.NoError(t, err)

		keyInvalid, err := generateAESKey(65)
		require.NoError(t, err)

		plaintext := []byte(uuid.New().String())

		_, err = aes.Encrypt(plaintext, keyInvalid)
		require.Error(t, err)

		cyphertext, err := aes.Encrypt(plaintext, keyValid)
		require.NoError(t, err)

		_, err = aes.Decrypt(cyphertext, keyInvalid)
		require.Error(t, err)
	})

	t.Run("Encrypt should fail with invalid cyphertext size", func(t *testing.T) {
		keyValid, err := generateAESKey(32)
		require.NoError(t, err)

		_, err = aes.Decrypt([]byte("fasdfsd"), keyValid)
		require.Error(t, err)
	})
}
