package server

import (
	"chat/shared/channel"
	"fmt"
	"net"
)

type Server struct {
	conn     *net.UDPConn
	receiver *channel.Receiver
}

func New(addr *net.UDPAddr) (*Server, error) {
	conn, err := net.ListenUDP("udp", addr)
	if err != nil {
		return nil, err
	}

	receiver := channel.NewReceiver(conn)
	receiver.OnSuccess(func(mt channel.MessageType, b []byte) {
		fmt.Println(mt)
		fmt.Println(string(b))
	})
	receiver.OnError(func(err error) {
		fmt.Println(err)
	})

	return &Server{
		conn:     conn,
		receiver: receiver,
	}, nil
}

func (s *Server) Listen() {
	s.receiver.Run()
}
