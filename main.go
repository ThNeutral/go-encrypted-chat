package main

import (
	"chat/client"
	"chat/server"
	"flag"
	"fmt"
	"net"
	"os"
)

func runServer(port int) error {
	addr := &net.UDPAddr{
		IP:   net.ParseIP("127.0.0.1"),
		Port: port,
	}

	s, err := server.New(addr)
	if err != nil {
		return err
	}

	s.Listen()

	return nil
}

func runClient(target string) error {
	addr, err := net.ResolveUDPAddr("udp", target)
	if err != nil {
		return err
	}

	c, err := client.New(addr)
	if err != nil {
		return err
	}

	payload := "8Nrnnd9Rp7PJBrC7gjRUWYW5ShazjbcQUJKNrwc3nmjARHIa8XGjr2buka3O8kYSRf3hoN9uRXFKLTTJRNSmoUMiqcjNTJd47HchSn4tl2fczsmHU2uCDtGsRERU1xwJ7nZFHQRx8nwffr0yIOugDsWgT0OZFOek0zaESMKBx7KlGh0fN1j8c5560U0Um5jjTHXLc8gU5STWgtQdKI73TCq1ulK0dziK7ODgoYwy7opItTKrOKK0dg5T73g8hvaq4PzIFzSe0VzqVuMu61nGmezRnW8UxdJnUhqdKnn2djGGwxOCRGYb0dnu0zSGGIiDgsNZxfbJboSmuzZKGLj8I6NjKXts9wMVPOuSbh6WyP0vzkWWP9l5gPSrxS4X3wbvtR4oqesAjPBFFxyDBLl3CjYNb3UYyzojWrK3IutiqoBUfSeKP469L3Hwno390LBhxhbwhBub9zlpmqQiVfABPVjDq68uj0mU9h4g9iWMKpZbOAPxvuqzu5DnaI9S5JhwuwMDuQGNBFnSkcnedfgwa8W8ITkqcsBL4V5ne1ZWOcVgXwo5hFfNfP69AyEUqmtQld4gUNOcUTUoIAFp7TUo75F1W0Mh8K7sjTHoihhsXBwdMfNv5ttrQmFl1KBXeJZr6tVNiyTCAQKST6MXjswolLGucjiUcyuQzWy1c3s5GpGom2vgdLvdDk0bTMg4eDseoHsZhgyeVyXppbXWFKQFM2t8NRtFyPUkI54ItdOXwVqvGyNUfT1QxmXZUn43DD1htLEw8gnHG8mlmJuXW9dSp9VabseKGpld12KH6k4WIIcD5Bp4pfVDF1u1tMJopOUo1tUXt0tvCXUBex78MgykdQOojDo8IMoU2eFhsoUuVGKttihW7mAg2sroen4iZB6wYnk6nTMHtpKNqk2iGR6VhRymFcP45IajP13gKjPUNPIwd1CnbL4onNyUD6LF97q2tELsEUSWjWWk12bTnT4AoIqxaeMNfrnUvIQgkMPjbE1cc01E5354uWkhcGjzjrFJEu11qFJxkjFbMm4GUdgDAlwNSGFzI5X316pvg4vTJeH9pYgsG847inSL1hheKwKrBK3SGu1Gc6RCHDZK3BGEY8dME8MoQJJw7hLFLg1tDn8hE0uelLFtk9mSthtHHiarfrRmUhKETF07w9nJTlEXLmAjgB2sVZ4fYM1t5ZkxXAhkcgVC4Gi1wYEi6DjKEnFovvc6vmjbgPjRXtHVA6CHdpoklObO1msmvGbg84V8yd81kI6rFzVUPMs79da5whrh17cXBE0ZGXSEw6akFnn9mHb0Yf7TpC4cxHfLnCtoAEoaY1L2QWnMXYrUVPSSdXCUfGsRZJvC2JELSAZzbJgWmRUtJPKs0wTmlpu2pa00S8xO4vE7AOgMQA9FL0BFiCxNuiSSj8RzYCwIytJnIWRyGcw6xf6lrG7bU2OH66qXKXxtfUo4GCKhNeWTGmdlYcrjJPakILBawx8OCrxFrqnzSNvyAZgM"
	err = c.Send([]byte(payload))
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
