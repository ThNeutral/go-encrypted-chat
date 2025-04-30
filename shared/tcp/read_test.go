package tcp_test

import (
	"chat/shared/tcp"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestReadAll(t *testing.T) {
	t.Run("successfully reads valid data until EndToken", func(t *testing.T) {
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		go func() {
			client.Write([]byte("hello world" + tcp.EndToken))
			client.Close()
		}()

		data, err := tcp.ReadAll(server)
		require.NoError(t, err)

		require.Equal(t, data, []byte("hello world"))
	})

	t.Run("returns io.EOF if EndToken was not recieved", func(t *testing.T) {
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		go func() {
			client.Write([]byte("incomplete data"))
			client.Close()
		}()

		_, err := tcp.ReadAll(server)
		require.Error(t, err)
	})

	t.Run("handles large data with EndToken at the end", func(t *testing.T) {
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		large := make([]byte, 8000)
		for i := range large {
			large[i] = 'a'
		}

		go func() {
			client.Write(large)
			client.Write([]byte(tcp.EndToken))
			client.Close()
		}()

		data, err := tcp.ReadAll(server)
		require.NoError(t, err)
		require.Equal(t, data, large)
	})
}
