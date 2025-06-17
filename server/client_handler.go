package main

import (
	"fmt"
	"tcp-chat/common"
	"tcp-chat/transport"
)

func (s *Server) HandleClient(conn transport.Connection) error {
	defer conn.Close()

	name, err := s.askName(conn)
	if err != nil {
		return fmt.Errorf("couldn't get name from user: %w", err)
	}

	user, err := s.clients.AddUser(common.User{
		Id:   s.pickNewId(),
		Name: name,
		Conn: conn,
	})

	if err != nil {
		return fmt.Errorf("")
	}

	go s.BroadcastMessages(&user)

	for {
		msg, err := common.ReceiveMessage(conn)

		if err != nil {
			go s.notifyUserLeft(user)

			s.clients.RemoveUser(user.Id)

			return fmt.Errorf("user %d disconected: %w", user.Id, err)
		}

		// добавляем им друг друга в список открытых чатов
		if msg.From != common.ServerId || msg.To != common.ServerId {
			chatID, err := s.chats.GetOrCreateChat(msg.From, msg.To)

			if err != nil {
				return fmt.Errorf("error creating or opening chat (id:%d): %w", chatID, err)
			}
		}

		// отправляем сообщение в канал для отправки
		user.Send <- msg
	}
}

func (s *Server) BroadcastMessages(user *common.User) error {
	for msg := range user.Send {
		user, err := s.clients.GetUser(msg.To)
		if err != nil {
			return fmt.Errorf("couldn't get reciever name: %w", err)
		}
		common.SendMessage(user.Conn, msg)
	}

	return nil
}

func (s *Server) askName(conn transport.Connection) (string, error) {
	err := common.SendMessage(conn, common.Message{
		From:    common.ServerId,
		To:      common.ServerId, // this is a bit shitcode, but i don't wanns pass id for no reason
		Type:    common.MessageRequestSenderName,
		Content: ""})

	if err != nil {
		return "", fmt.Errorf("error asking name from user: %w", err)
	}

	nameMessage, err := common.ReceiveMessage(conn)

	if err != nil {
		return "", fmt.Errorf("error getting name from user: %w", err)
	}

	return nameMessage.Content, nil
}

func (s *Server) pickNewId() uint64 {
	s.mutex.Lock()
	s.lastId += 1
	newId := s.lastId
	s.mutex.Unlock()

	return newId
}

func (s *Server) notifyUserLeft(user common.User) error {
	chats, err := s.chats.GetUserChats(user.Id)
	if err != nil {
		return fmt.Errorf("error getting user chats list: %w", err)
	}

	for _, receiverId := range chats {
		user.Send <- common.Message{
			From:    common.ServerId,
			To:      receiverId,
			Type:    common.MessageQuitChat,
			Content: user.Name}
	}

	return nil
}
