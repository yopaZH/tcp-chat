package main

import (
	"errors"
	"fmt"
	"net"
	"sync"
	"tcp-chat/common"
)

type Server struct {
	listener net.Listener
    clients map[uint64]common.User
    broadcast chan common.Message
    
    mutex sync.Mutex

	lastId uint64
}

func(s *Server) StartServer(port string) {
	var err error
	s.listener, err = net.Listen("tcp", port)
	if err != nil {
		fmt.Println("error starting up server")
	}
	defer s.listener.Close()

	fmt.Println("server is up on port:", port)
	
	s.clients = make(map[uint64]common.User)
	s.broadcast = make(chan common.Message)

	// рассылаем все пришедшие сообщения
	go broadcastMessages(s.clients, s.broadcast, s.mutex)

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			fmt.Println("error connecting to client: ", err)
		}

		go HandleClient(conn)
	}
}