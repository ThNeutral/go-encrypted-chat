package tcp

import (
	"errors"
	"fmt"
)

type MessageType byte

const (
	CLIENT_HELLO MessageType = 0b00000000
	SERVER_HELLO MessageType = 0b00000001

	TESTING MessageType = 0b01111111
	INVALID MessageType = 0b11111111
)

func (mt MessageType) Validate() error {
	if mt == CLIENT_HELLO || mt == SERVER_HELLO || mt == TESTING || mt == INVALID {
		return nil
	}
	return fmt.Errorf("invalid MessageType: %v", mt.Byte())
}

func (mt MessageType) Byte() byte {
	return byte(mt)
}

func UnexpectedMessageType(want, have MessageType) error {
	return fmt.Errorf("unexpected message type. Want %v, Have %v", want, have)
}

type Message struct {
	Type    MessageType
	Payload []byte
}

func (m Message) Serialize() []byte {
	bytes := []byte{m.Type.Byte()}
	bytes = append(bytes, m.Payload...)
	return bytes
}

func Deserialize(bytes []byte) (Message, error) {
	var message Message
	if len(bytes) == 0 {
		return message, errors.New("bytes are not deserializable into Message")
	}

	messageType := MessageType(bytes[0])
	err := messageType.Validate()
	if err != nil {
		return message, err
	}

	message.Type = messageType
	message.Payload = bytes[1:]

	return message, nil
}
