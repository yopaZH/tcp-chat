package main

import (
	"fmt"
	"sync"
	"tcp-chat/common"
	"tcp-chat/services"
	"tcp-chat/session"
	"tcp-chat/storage"
	"tcp-chat/transport"
)

type Server struct {
	listener transport.ConnectionListener

	userService services.UserService
	chats       storage.ChatRepository

	mutex sync.Mutex

	lastId uint64
}

func NewServer(port string) (*Server, error) {
	listener, err := transport.NewTCPListener(port)
	if err != nil {
		return nil, fmt.Errorf("failed to listen: %w", err)
	}

	var s Server = Server{
		listener:    listener,
		userService: services.NewUserService(storage.NewInMemoryUserRepo(), session.NewInMemoryManager()),
		chats:       storage.NewMemoryChatRepo(),
		mutex:       sync.Mutex{},
		lastId:      0,
	}

	// резервируем id 0 под сервер
	_, _, err = s.userService.Register(common.ServerId, "server", nil)

	if err != nil {
		return nil, fmt.Errorf("error reserving user (id:%d) for server: %w", common.ServerId, err)
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

	err = s.userService.Close()
	if err != nil {
		return err
	}

	return nil
}
