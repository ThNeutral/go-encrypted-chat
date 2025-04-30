package tcp_test

import (
	"chat/shared/tcp"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestMessageAndType(t *testing.T) {
	t.Run("it successfully validates valid types", func(t *testing.T) {
		validTypes := tcp.ValidMessageTypes
		for _, mt := range validTypes {
			err := mt.Validate()
			require.NoError(t, err)
		}
	})

	t.Run("it fails on validation of invalid type", func(t *testing.T) {
		invalidType := tcp.INVALID
		err := invalidType.Validate()
		require.Error(t, err)
	})

	t.Run("it successfully serializes and deserializes message", func(t *testing.T) {
		original := tcp.Message{
			Type:    tcp.CLIENT_HELLO,
			Payload: []byte("payload"),
		}
		serialized := original.Serialize()
		deserialized, err := tcp.Deserialize(serialized)
		require.NoError(t, err)

		require.Equal(t, deserialized.Type, original.Type)
		require.Equal(t, deserialized.Payload, original.Payload)
	})

	t.Run("fails to deserialize empty bytes", func(t *testing.T) {
		_, err := tcp.Deserialize([]byte{})
		require.Error(t, err)
	})

	t.Run("fails to serialize invalid message type", func(t *testing.T) {
		data := []byte{tcp.INVALID.Byte(), 0xAA}
		_, err := tcp.Deserialize(data)
		require.Error(t, err)
	})
}
