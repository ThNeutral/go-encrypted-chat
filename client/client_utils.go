package client

import (
	"chat/shared/aes"
	"chat/shared/tcp"
)

func (c *Client) writeBytes(payload []byte) error {
	return tcp.WriteAll(c.conn, payload)
}

func (c *Client) writeMessage(message tcp.Message) error {
	bytes := message.Serialize()
	return c.writeBytes(bytes)
}

func (c *Client) encryptAndWriteMessage(message tcp.Message) error {
	key, err := c.ecdh.SharedSecret()
	if err != nil {
		return err
	}

	plaintext := message.Serialize()

	cyphertext, err := aes.Encrypt(plaintext, key)
	if err != nil {
		return err
	}

	return c.writeBytes(cyphertext)
}

func (c *Client) readBytes() ([]byte, error) {
	return tcp.ReadAll(c.conn)
}

func (c *Client) readMessage() (tcp.Message, error) {
	var message tcp.Message

	bytes, err := c.readBytes()
	if err != nil {
		return message, err
	}

	return tcp.Deserialize(bytes)
}

func (c *Client) readAndDencryptMessage() (tcp.Message, error) {
	var message tcp.Message

	key, err := c.ecdh.SharedSecret()
	if err != nil {
		return message, err
	}

	cyphertext, err := c.readBytes()
	if err != nil {
		return message, err
	}

	plaintext, err := aes.Decrypt(cyphertext, key)
	if err != nil {
		return message, err
	}

	return tcp.Deserialize(plaintext)
}
