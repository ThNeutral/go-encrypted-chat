package tcp_test

import (
	"chat/shared/tcp"
	"io"
	"net"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestWriteAll(t *testing.T) {
	t.Run("writes full data with EndToken", func(t *testing.T) {
		server, client := net.Pipe()
		defer server.Close()
		defer client.Close()

		originalData := []byte("test message")

		go func() {
			err := tcp.WriteAll(client, originalData)
			require.Error(t, err)
			client.Close()
		}()

		buf := make([]byte, 1024)
		n, err := server.Read(buf)
		if err != nil && err != io.EOF {
			t.Errorf("unexpected read error: %v", err)
		}

		expected := append(originalData, []byte(tcp.EndToken)...)
		require.Equal(t, buf[:n], expected)
	})
}
