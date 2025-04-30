package config

import (
	"errors"
	"flag"
	"net"

	"github.com/google/uuid"
)

type AppType string

const (
	SERVER AppType = "server"
	CLIENT AppType = "client"
)

type Config interface {
	Init(args []string) error
	Type() AppType
	Addr() *net.TCPAddr
}

func Instance() Config {
	if inst == nil {
		inst = NewConfig().(*impl)
	}
	return inst
}

var inst *impl

// Only for testing
func NewConfig() Config {
	return &impl{
		fs: flag.NewFlagSet(uuid.New().String(), flag.ContinueOnError),
	}
}

type impl struct {
	t    AppType
	addr *net.TCPAddr
	fs   *flag.FlagSet
}

func (i *impl) Init(args []string) error {
	serverFlag := i.fs.Bool("server", false, "Run as server")
	clientFlag := i.fs.Bool("client", false, "Run as client")
	addrFlag := i.fs.String("addr", "", "Addr for the server to listen on or for client to connect to")

	err := i.fs.Parse(args)
	if err != nil {
		return err
	}

	if *serverFlag && *clientFlag {
		return errors.New("cannot run app as server and client at the same time")
	}

	if *serverFlag {
		i.t = SERVER
	} else if *clientFlag {
		i.t = CLIENT
	} else {
		return errors.New("must provide app type as flag. Available app type flags are --server OR --client")
	}

	if *addrFlag == "" {
		return errors.New("address must be provided using --addr=host:port")
	}

	addr, err := net.ResolveTCPAddr("tcp", *addrFlag)
	if err != nil {
		return err
	} else {
		i.addr = addr
	}

	return nil
}
func (i *impl) Type() AppType {
	return i.t
}

func (i *impl) Addr() *net.TCPAddr {
	return i.addr
}
