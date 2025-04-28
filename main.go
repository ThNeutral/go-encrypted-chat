package main

import (
	"chat/client"
	"chat/server"
	"chat/shared/tcp"
	"flag"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"
)

func runServer(port int) error {
	addr := &net.TCPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	}

	s, err := server.New(addr)
	if err != nil {
		return err
	}

	log.Println("Listening on :" + strconv.Itoa(addr.Port))
	s.Listen()

	return nil
}

func runClient(target string) error {
	addr, err := net.ResolveTCPAddr("tcp", target)
	if err != nil {
		return err
	}

	c, err := client.New(addr)
	if err != nil {
		return err
	}

	payload := "This is message that I want to send via TCP."
	err = c.Initialize()
	if err != nil {
		return err
	}

	err = c.SendMessage(tcp.Message{
		Type:    tcp.TESTING,
		Payload: []byte(payload),
	})
	if err != nil {
		return err
	}

	return nil
}

func main() {
	serverFlag := flag.Bool("server", false, "Run as server")
	portFlag := flag.Int("port", 0, "Port for the server to listen on")

	clientFlag := flag.Bool("client", false, "Run as client")
	targetFlag := flag.String("target", "", "Target host:port for client to connect to")

	flag.Parse()

	if *serverFlag && *clientFlag {
		fmt.Println("Cannot use --server and --client at the same time")
		os.Exit(1)
	}

	if *serverFlag {
		if *portFlag == 0 {
			fmt.Println("--port must be specified when running as server")
			os.Exit(1)
		}

		err := runServer(*portFlag)

		if err != nil {
			fmt.Println(err)
		}
	} else if *clientFlag {
		if *targetFlag == "" {
			fmt.Println("--target must be specified when running as client")
			os.Exit(1)
		}

		err := runClient(*targetFlag)

		if err != nil {
			fmt.Println(err)
		}
	} else {
		fmt.Println("Must specify either --server or --client")
		os.Exit(1)
	}
}
