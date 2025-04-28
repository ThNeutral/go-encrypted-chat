package tcp

import (
	"bytes"
	"io"
	"net"
)

func ReadAll(conn net.Conn) ([]byte, error) {
	data := make([]byte, 0)
	buffer := make([]byte, 4096)
	for {
		n, err := conn.Read(buffer)

		if n > 0 {
			data = append(data, buffer[:n]...)
		}

		if err == io.EOF {
			if len(data) >= len(EndToken) && bytes.Equal(data[len(data)-len(EndToken):], []byte(EndToken)) {
				return data[:len(data)-len(EndToken)], io.EOF
			}
			return data, io.EOF
		}
		if err != nil {
			return nil, err
		}

		if len(data) >= len(EndToken) && bytes.Equal(data[len(data)-len(EndToken):], []byte(EndToken)) {
			return data[:len(data)-len(EndToken)], nil
		}
	}
}
