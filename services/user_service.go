package services

import (
	"fmt"
	"tcp-chat/common"
	"tcp-chat/session"
	"tcp-chat/storage"
	"tcp-chat/transport"
)

type UserService interface {
	Register(id uint64, name string, conn transport.Connection) (*common.User, *session.Session, error)
	Remove(id uint64)
	GetUser(id uint64) (*common.User, *session.Session, error)
	SendMessage(to uint64, msg common.Message) error
	Close() error
}

type userService struct {
	repo    storage.UserRepository
	manager session.Manager
}

func NewUserService(repo storage.UserRepository, manager session.Manager) *userService {
	return &userService{
		repo:    repo,
		manager: manager,
	}
}

func (s *userService) Register(id uint64, name string, conn transport.Connection) (*common.User, *session.Session, error) {
	user := common.User{
		Id:   id,
		Name: name,
	}

	if _, err := s.repo.AddUser(user); err != nil {
		return nil, nil, fmt.Errorf("add user: %w", err)
	}

	sess := s.manager.Add(id, conn)

	return &user, sess, nil
}

func (s *userService) Remove(id uint64) {
	_ = s.repo.RemoveUser(id) // можно игнорить ошибку, если не нашли — всё ок

	s.manager.Remove(id)
}

func (s *userService) GetUser(id uint64) (*common.User, *session.Session, error) {
	user, err := s.repo.GetUser(id)
	if err != nil {
		return nil, nil, fmt.Errorf("get user: %w", err)
	}

	sess, ok := s.manager.Get(id)
	if !ok {
		return nil, nil, fmt.Errorf("session for user %d not found", id)
	}

	return user, sess, nil
}

func (s *userService) Disconnect(id uint64) error {
	// Удаляем пользователя из репозитория (модель)
	if err := s.repo.RemoveUser(id); err != nil {
		// Можно логировать, но не прерывать — пусть идём дальше
		fmt.Printf("warning: failed to remove user from repo: %v\n", err)
	}

	s.manager.Remove(id)
	/*
		if err := s.manager.Remove(id); err != nil {
			return fmt.Errorf("failed to remove session for user %d: %w", id, err)
		}
	*/
	return nil
}

func (s *userService) SendMessage(to uint64, msg common.Message) error {
	return s.manager.Broadcast(msg)
}

func (s *userService) Close() error {
	if err := s.repo.Close(); err != nil {
		return fmt.Errorf("failed to close repository: %w", err)
	}

	s.manager.Shutdown()

	return nil
}
