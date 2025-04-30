package config

import (
	"errors"
	"flag"
	"net"
)

type AppType string

const (
	SERVER AppType = "server"
	CLIENT AppType = "client"
)

type Config interface {
	Init() []error
	Type() AppType
	Addr() *net.TCPAddr
}

func Instance() Config {
	return &inst
}

var inst impl

type impl struct {
	t    AppType
	addr *net.TCPAddr
}

func (i *impl) Init() []error {
	serverFlag := flag.Bool("server", false, "Run as server")
	clientFlag := flag.Bool("client", false, "Run as client")
	addrFlag := flag.String("addr", "", "Addr for the server to listen on or for client to connect to")

	flag.Parse()

	var errs []error
	if *serverFlag && *clientFlag {
		errs = append(errs, errors.New("cannot run app as server and client at the same time"))
	}

	if *serverFlag {
		i.t = SERVER
	} else if *clientFlag {
		i.t = CLIENT
	} else {
		errs = append(errs, errors.New("must provide app type as flag. Available app type flags are --server OR --client"))
	}

	if *addrFlag == "" {
		errs = append(errs, errors.New("address must be provided using --addr=host:port"))
	} else {
		addr, err := net.ResolveTCPAddr("tcp", *addrFlag)
		if err != nil {
			errs = append(errs, err)
		} else {
			i.addr = addr
		}
	}

	if len(errs) != 0 {
		return errs
	}

	return nil
}

func (i *impl) Type() AppType {
	return i.t
}

func (i *impl) Addr() *net.TCPAddr {
	return i.addr
}
