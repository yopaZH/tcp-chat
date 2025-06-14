package main

import (
	"fmt"
	"tcp-chat/common"
	"tcp-chat/transport"
	"tcp-chat/utils"
)

// go HandleClient(conn, &clients, s.broadcast)

func (s *Server) HandleClient(conn transport.Connection) error {
	defer conn.Close()

	currentId := s.pickNewId()

	// TODO: поправить приём имени
	nameCh := make(chan NameResult)
	go s.askName(currentId, nameCh)
	askNameResult := <-nameCh
	if askNameResult.Err != nil {
		return fmt.Errorf("error recieving name: %w", askNameResult.Err)
	}

	s.clients.AddUser(common.User{
		Id:        currentId,
		Name:      askNameResult.Name,
		Conn:      conn,
		ChatsWith: make(map[uint64]struct{}),
	})

	for {
		msg, err := common.ReceiveMessage(conn)

		if err != nil {
			//s.broadcast <- fmt.Sprintf("%s quit chat", name)
			s.notifyUserLeft(currentId)

			s.mutex.Lock()
			s.clients.RemoveUser(currentId)
			s.mutex.Unlock()

			return fmt.Errorf("user %d disconected: %w", currentId, err)
		} else {
			// добавляем им друг друга в список открытых чатов
			if msg.From != common.ServerId || msg.To != common.ServerId {
				s.updateChatsLists(msg.From, msg.To)
			}

			s.broadcast <- msg
		}
	}
}

type NameResult struct {
	Name string
	Err  error
}

func (s *Server) askName(id uint64, ch chan<- NameResult) {
	s.broadcast <- common.Message{
		From:    common.ServerId,
		To:      id,
		Type:    common.MessageRequestSenderName,
		Content: ""}

	user, err := s.clients.GetUser(id)
	nameMessage, err := common.ReceiveMessage(user.Conn)

	if err != nil {
		ch <- NameResult{Name: "", Err: fmt.Errorf("couldn't get a name: %w", err)}
	}

	ch <- NameResult{Name: nameMessage.Content, Err: nil}
}

func (s *Server) pickNewId() uint64 {
	s.mutex.Lock()
	s.lastId += 1
	newId := s.lastId
	s.mutex.Unlock()

	return newId
}

func (s *Server) notifyUserLeft(id uint64) error {
	s.mutex.Lock()

	user, err := s.clients.GetUser(id)
	if err != nil {
		return err
	}

	for receiverId := range user.ChatsWith {
		s.broadcast <- common.Message{
			From:    common.ServerId,
			To:      receiverId,
			Type:    common.MessageQuitChat,
			Content: user.Name}
	}

	s.mutex.Unlock()

	return nil
}

func (s *Server) updateChatsLists(from uint64, to uint64) error {
	s.mutex.Lock()
	fromUser, err := s.clients.GetUser(from)
	if err != nil {
		return err
	}
	toUser, err := s.clients.GetUser(to)
	if err != nil {
		return err
	}

	if !utils.Exists(to, fromUser.ChatsWith) {
		fromUser.ChatsWith[to] = struct{}{}
		toUser.ChatsWith[from] = struct{}{}
	}
	s.mutex.Unlock()

	return nil
}
