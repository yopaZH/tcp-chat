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

	// резервируем id 0 под сервер
	s.AddUser("server", nil)

	// рассылаем все пришедшие сообщения
	go s.BroadcastMessages()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			fmt.Println("error connecting to client: ", err)
		}

		go s.HandleClient(conn)
	}
}

func (s *Server)BroadcastMessages() {
	for msg := range s.broadcast {
		s.mutex.Lock()
		common.SendMessage(s.clients[msg.To].Conn, msg)

		s.mutex.Unlock()
	}
}