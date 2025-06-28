package main

import (
	"fmt"
	"tcp-chat/config"
)

func main() {
	cfg, err := config.LoadConfig("config/config.yaml")
	if err != nil {
		fmt.Printf("failed to load config: %v", err)
	}

	server, err := NewServer(cfg.Server.Address)
	defer func() {
		err = server.Shutdown()
		if err != nil {
			fmt.Println("error shuting down the server: %w", err)
		}
	}()

	if err != nil {
		fmt.Printf("error starting server: %e", err)
	} else {
		fmt.Println("server is up on port:", cfg.Server.Address)
	}

	server.StartServer()
}
