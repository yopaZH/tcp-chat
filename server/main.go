package main

import (
	"fmt"
	//"tcp-chat/server"
)

func main() {
	port := ":8080"
	server, err := NewServer(port)

	if err != nil {
		fmt.Printf("error starting server: %e", err)
	} else {
		fmt.Println("server is up on port:", port)
	}

	defer server.listener.Close()

	server.StartServer()
}
