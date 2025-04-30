package main

import (
	"chat/client"
	"chat/server"
	"chat/shared/config"
	"fmt"
	"os"
)

func main() {
	inst := config.Instance()
	err := inst.Init(os.Args)
	if err != nil {
		fmt.Printf("Failed to init config: %v\n", err.Error())
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
