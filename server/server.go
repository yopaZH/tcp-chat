package main

import (
	//"errors"
	"fmt"
	"net"
	"sync"
	"tcp-chat/common"
)

type Server struct {
	listener  net.Listener
	clients   map[uint64]common.User
	broadcast chan common.Message

	mutex sync.Mutex

	lastId uint64
}

func NewServer(port string) (*Server, error) {
	listener, err := net.Listen("tcp", port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	var s Server = Server{
		listener:  listener,
		clients:   make(map[uint64]common.User),
		broadcast: make(chan common.Message),
		mutex:     sync.Mutex{},
		lastId:    1,
	}

	// резервируем id 0 под сервер
	s.clients[s.lastId] = common.User{Id: 0, Name: "server", Conn: nil, ChatsWith: nil}

	return &s, nil
}

func (s *Server) StartServer() {
	s.clients = make(map[uint64]common.User)
	s.broadcast = make(chan common.Message)

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

func (s *Server) BroadcastMessages() {
	for msg := range s.broadcast {
		s.mutex.Lock()
		common.SendMessage(s.clients[msg.To].Conn, msg)

		s.mutex.Unlock()
	}
}
