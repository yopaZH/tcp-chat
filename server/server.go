package main

import (
	//"errors"
	"fmt"
	"sync"
	"tcp-chat/common"
	"tcp-chat/storage"
	"tcp-chat/transport"
)

type Server struct {
	listener  transport.ConnectionListener
	clients   storage.UserRepository
	broadcast chan common.Message

	mutex sync.Mutex

	lastId uint64
}

func NewServer(port string) (*Server, error) {
	listener, err := transport.NewTCPListener(port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	var s Server = Server{
		listener:  listener,
		clients:   storage.NewMemoryStorage(),
		broadcast: make(chan common.Message),
		mutex:     sync.Mutex{},
		lastId:    1,
	}

	// резервируем id 0 под сервер
	s.clients.AddUser(common.User{Id: 0, Name: "server", Conn: nil, ChatsWith: nil})

	return &s, nil
}

func (s *Server) StartServer() {
	// рассылаем все приходящие сообщения
	go s.BroadcastMessages()

	for {
		conn, err := s.listener.Accept()

		if err != nil {
			fmt.Println("error connecting to client:", err)
		}

		go s.HandleClient(conn)
	}
}

func (s *Server) Shutdown() error {
	err := s.listener.Close()
	if err != nil {
		return err
	}

	err = s.clients.Close()
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) BroadcastMessages() {
	for msg := range s.broadcast {
		s.mutex.Lock()
		user, err := s.clients.GetUser(msg.To)
		if err != nil {
		}
		common.SendMessage(user.Conn, msg)

		s.mutex.Unlock()
	}
}
