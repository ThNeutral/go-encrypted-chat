package server

import (
	"chat/shared/ecdh"
	"chat/shared/tcp"
	"io"
	"log"
	"net"
)

type Handler struct {
	conn net.Conn
	ecdh ecdh.ECDH
}

func (h *Handler) Start() {
	if err := h.diffieHellman(); err != nil {
		log.Printf("Encountered error during ECDH: %v\n", err)
		h.cleanup()
		return
	}

	err := h.listen()
	if err == io.EOF {
		log.Println("Client " + h.conn.LocalAddr().String() + " has disconnected")
		h.cleanup()
		return
	} else if err != nil {
		log.Printf("Listener has failed: %v\n", err)
		h.cleanup()
		return
	}

	h.cleanup()
}
func (h *Handler) diffieHellman() error {
	clientHello, err := h.readMessage()
	if err != nil {
		return err
	}
	if clientHello.Type != tcp.CLIENT_HELLO {
		return tcp.UnexpectedMessageType(tcp.CLIENT_HELLO, clientHello.Type)
	}

	clientPub := clientHello.Payload

	serverHello := tcp.Message{
		Type:    tcp.SERVER_HELLO,
		Payload: h.ecdh.PublicKey(),
	}
	err = h.writeMessage(serverHello)
	if err != nil {
		return err
	}

	err = h.ecdh.GenerateSharedSecret(clientPub)
	if err != nil {
		return err
	}

	return nil
}

func (h *Handler) listen() error {
	for {
		message, err := h.readAndDecryptMessage()
		if err != nil {
			return err
		}

		log.Println(string(message.Payload))
	}
}

func (h *Handler) cleanup() {
	h.conn.Close()
}
