package client

import (
	"chat/shared/channel"
	"chat/shared/ecdh"
	"net"
)

type Client struct {
	conn          *net.UDPConn
	ecdh          *ecdh.ECDH
	peerPublicKey []byte
}

func New(target *net.UDPAddr) (*Client, error) {
	ecdh, err := ecdh.NewX25519()
	if err != nil {
		return nil, err
	}

	client, err := net.DialUDP("udp", nil, target)
	if err != nil {
		return nil, err
	}

	return &Client{
		ecdh:          ecdh,
		conn:          client,
		peerPublicKey: nil,
	}, nil
}

func (c *Client) Send(payload []byte) error {
	return channel.Send(c.conn, channel.MESSAGE_TYPE_ECDH_CLIENT_PUB, payload)
}
