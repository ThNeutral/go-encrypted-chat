package server

import (
	"chat/shared/ecdh"
	"log"
	"net"
)

type Server struct {
	listener *net.TCPListener
}

func New(addr *net.TCPAddr) (*Server, error) {
	listener, err := net.ListenTCP("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Server{
		listener: listener,
	}, nil
}

func (s *Server) Listen() {
	for {
		conn, err := s.listener.Accept()
		if err != nil {
			log.Println("Failed to accept connection")
			continue
		}

		ecdh, err := ecdh.NewX25519()
		if err != nil {
			log.Println("Failed to create ECDH instance")
			conn.Close()
			continue
		}

		handler := &Handler{
			conn: conn,
			ecdh: ecdh,
		}

		log.Println("Client " + conn.LocalAddr().String() + " has disconnected")
		go handler.Start()
	}
}
