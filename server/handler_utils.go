package server

import (
	"chat/shared/tcp"
)

func (h *Handler) readBytes() ([]byte, error) {
	return tcp.ReadAll(h.conn)
}

func (h *Handler) readMessage() (tcp.Message, error) {
	var message tcp.Message

	bytes, err := h.readBytes()
	if err != nil {
		return message, err
	}

	return tcp.Deserialize(bytes)
}

func (h *Handler) readAndDecryptMessage() (tcp.Message, error) {
	var message tcp.Message

	key, err := h.ecdh.SharedSecret()
	if err != nil {
		return message, err
	}

	cyphertext, err := h.readBytes()
	if err != nil {
		return message, err
	}

	plaintext, err := tcp.Decrypt(cyphertext, key)
	if err != nil {
		return message, err
	}

	return tcp.Deserialize(plaintext)
}

func (h *Handler) writeBytes(data []byte) error {
	return tcp.WriteAll(h.conn, data)
}

func (h *Handler) writeMessage(message tcp.Message) error {
	bytes := message.Serialize()
	return h.writeBytes(bytes)
}

func (h *Handler) encryptAndWriteMessage(message tcp.Message) error {
	key, err := h.ecdh.SharedSecret()
	if err != nil {
		return err
	}

	plaintext := message.Serialize()

	cyphertext, err := tcp.Encrypt(plaintext, key)
	if err != nil {
		return err
	}

	return h.writeBytes(cyphertext)
}
