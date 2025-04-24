package channel

import (
	"encoding/binary"
	"errors"
	"hash/crc32"
)

type MessageType uint8

const (
	MESSAGE_TYPE_ECDH_CLIENT_PUB   MessageType = 0
	MESSAGE_TYPE_ECDH_SERVER_PUB   MessageType = 1
	MESSAGE_TYPE_VERIFY_CLIENT_KEY MessageType = 2
	MESSAGE_TYPE_VERIFY_SERVER_KEY MessageType = 3

	MESSAGE_TYPE_INVALID MessageType = 255
)

const (
	PACKET_HEADERS_LENGTH = 4 + 2 + 2 + 4 + 4
	MAX_PAYLOAD_SIZE      = 1000
	MAX_PACKET_SIZE       = PACKET_HEADERS_LENGTH + MAX_PAYLOAD_SIZE
)

type Packet struct {
	MessageID      uint32
	TotalFragments uint16
	FragmentIndex  uint16
	PayloadLength  uint32
	Checksum       uint32
	Payload        []byte
}

func (p Packet) Serialize() ([]byte, int) {
	messageLength := PACKET_HEADERS_LENGTH + len(p.Payload)
	message := make([]byte, messageLength)

	binary.BigEndian.PutUint32(message[0:4], p.MessageID)
	binary.BigEndian.PutUint16(message[4:6], p.TotalFragments)
	binary.BigEndian.PutUint16(message[6:8], p.FragmentIndex)
	binary.BigEndian.PutUint32(message[8:12], p.PayloadLength)
	binary.BigEndian.PutUint32(message[12:16], p.Checksum)
	copy(message[16:], p.Payload)

	return message, messageLength
}

func Deserialize(data []byte) (Packet, error) {
	if len(data) < 16 {
		return Packet{}, errors.New("data too short to be a valid Packet")
	}

	p := Packet{}
	p.MessageID = binary.BigEndian.Uint32(data[0:4])
	p.TotalFragments = binary.BigEndian.Uint16(data[4:6])
	p.FragmentIndex = binary.BigEndian.Uint16(data[6:8])
	p.PayloadLength = binary.BigEndian.Uint32(data[8:12])
	p.Checksum = binary.BigEndian.Uint32(data[12:16])

	if len(data) != int(16+p.PayloadLength) {
		return Packet{}, errors.New("data length does not match payload length")
	}

	payload := data[16 : 16+int(p.PayloadLength)]
	if len(payload) == 0 {
		return Packet{}, errors.New("payload is empty")
	}
	p.Payload = make([]byte, len(payload))
	copy(p.Payload, payload)

	calculatedChecksum := crc32.ChecksumIEEE(p.Payload)
	if calculatedChecksum != p.Checksum {
		return Packet{}, errors.New("checksum mismatch")
	}

	return p, nil
}
