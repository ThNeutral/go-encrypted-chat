package channel

import (
	"fmt"
	"hash/crc32"
	"net"
	"sync/atomic"
)

var nextMessageID = uint32(0)

func Send(conn *net.UDPConn, mtype MessageType, bytes []byte) error {
	messageID := atomic.AddUint32(&nextMessageID, 1)

	bytes = append([]byte{byte(mtype)}, bytes...)
	bytesLength := len(bytes)
	totalFragments := uint16((bytesLength + MAX_PAYLOAD_SIZE - 1) / MAX_PAYLOAD_SIZE)

	for i := uint16(0); i < totalFragments; i++ {
		start := int(i) * MAX_PAYLOAD_SIZE
		end := start + MAX_PAYLOAD_SIZE
		if end >= bytesLength {
			end = bytesLength
		}

		payload := bytes[start:end]
		checksum := crc32.ChecksumIEEE(payload)

		packet := Packet{
			MessageID:      messageID,
			TotalFragments: totalFragments,
			FragmentIndex:  i,
			PayloadLength:  uint32(end - start),
			Payload:        payload,
			Checksum:       checksum,
		}

		message, messageLength := packet.Serialize()

		writtenLength, err := conn.Write(message)
		if err != nil {
			return err
		}
		if writtenLength != messageLength {
			return fmt.Errorf("failed to send: expected to send %v bytes, actually sent %v bytes", messageLength, writtenLength)
		}
	}

	return nil
}
