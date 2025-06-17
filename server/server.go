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
	listener transport.ConnectionListener
	clients  storage.UserRepository
	chats    storage.ChatRepository

	mutex sync.Mutex

	lastId uint64
}

func NewServer(port string) (*Server, error) {
	listener, err := transport.NewTCPListener(port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	var s Server = Server{
		listener: listener,
		clients:  storage.NewInMemoryUserRepo(),
		chats:    storage.NewMemoryChatRepo(),
		mutex:    sync.Mutex{},
		lastId:   1,
	}

	// резервируем id 0 под сервер
	_, err = s.clients.AddUser(common.User{
		Id:   0,
		Name: "server",
		Conn: nil})

	if err != nil {
		return nil, fmt.Errorf("error reserving user (id:0) for server: %w", err)
	}

	return &s, nil
}

func (s *Server) StartServer() error {
	for {
		conn, err := s.listener.Accept()

		if err != nil {
			return fmt.Errorf("error connecting to client: %w", err)
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
