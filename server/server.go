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
	go BroadcastMessages(&s.clients, s.broadcast, &s.mutex)

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			fmt.Println("error connecting to client: ", err)
		}

		go HandleClient(conn, &clients, s.broadcast)
	}
}

func(s *Server) AddUser(name string, conn net.Conn) {
	s.mutex.Lock()

	s.lastId += 1
	newUser := common.User{Id: s.lastId, Name: name, Conn: conn}
	s.clients[s.lastId] = newUser
	
	s.mutex.Unlock()
}