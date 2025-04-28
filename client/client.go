package client

import (
	"chat/shared/ecdh"
	"chat/shared/tcp"
	"net"
)

type Client struct {
	conn net.Conn
	ecdh *ecdh.ECDH
}

func New(target *net.TCPAddr) (*Client, error) {
	ecdh, err := ecdh.NewX25519()
	if err != nil {
		return nil, err
	}

	client, err := net.DialTCP("tcp", nil, target)
	if err != nil {
		return nil, err
	}

	return &Client{
		ecdh: ecdh,
		conn: client,
	}, nil
}

func (c *Client) Initialize() error {
	err := c.diffieHellman()
	if err != nil {
		return err
	}

	return nil
}

func (c *Client) SendMessage(message tcp.Message) error {
	return c.encryptAndWriteMessage(message)
}

func (c *Client) diffieHellman() error {
	clientHello := tcp.Message{
		Type:    tcp.CLIENT_HELLO,
		Payload: c.ecdh.PublicKey(),
	}

	err := c.writeMessage(clientHello)
	if err != nil {
		return err
	}

	serverHello, err := c.readMessage()
	if err != nil {
		return err
	}
	if serverHello.Type != tcp.SERVER_HELLO {
		return tcp.UnexpectedMessageType(tcp.SERVER_HELLO, serverHello.Type)
	}

	err = c.ecdh.GenerateSharedSecret(serverHello.Payload)
	if err != nil {
		return err
	}

	return nil
}
