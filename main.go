package main

import (
	"chat/client"
	"chat/server"
	"chat/shared/config"
	"fmt"
)

func main() {
	inst := config.Instance()
	errs := inst.Init()
	if errs != nil {
		for _, err := range errs {
			fmt.Println(err)
		}
		return
	}

	if inst.Type() == config.SERVER {
		s, err := server.New(inst.Addr())
		if err != nil {
			fmt.Printf("Failed to create server: %v", err)
			return
		}

		s.Listen()
	}

	if inst.Type() == config.CLIENT {
		c, err := client.New(inst.Addr())
		if err != nil {
			fmt.Printf("Failed to create client: %v", err)
			return
		}

		err = c.Initialize()
		if err != nil {
			fmt.Printf("Failed to initialize client: %v", err)
			return
		}
	}
}
