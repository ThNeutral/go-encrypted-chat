package tcp

import "net"

func WriteAll(conn net.Conn, data []byte) error {
	data = append(data, []byte(EndToken)...)
	total := 0
	for total < len(data) {
		n, err := conn.Write(data[total:])
		if err != nil {
			return err
		}
		total += n
	}
	return nil
}
