package main

import (
	"fmt"
	"tcp-chat/common"
	"tcp-chat/transport"
)

func (s *Server) HandleClient(conn transport.Connection) error {
	defer conn.Close()

	id := s.pickNewId()

	nameMessage, err := common.ReceiveMessage(conn)
	if err != nil {
		return fmt.Errorf("error getting name from user: %w", err)
	}

	user, sess, err := s.userService.Register(id, nameMessage.FromName, conn)

	if err != nil {
		return fmt.Errorf("failed to add user: %w", err)
	}

	// issuing id to user
	common.SendMessage(conn, common.Message{
		From: common.ServerId,
		To:   user.Id,
		Type: common.MessageIdIssuing,
	})

	go s.BroadcastMessages(id)

	for {
		msg, err := common.ReceiveMessage(conn)

		if err != nil {
			go s.notifyUserLeft(id)

			s.userService.Remove(id)

			return fmt.Errorf("user %d disconected: %w", user.Id, err)
		}

		// добавляем им друг друга в список открытых чатов
		if msg.From != common.ServerId || msg.To != common.ServerId {
			chatID, err := s.chats.GetOrCreateChat(msg.From, msg.To)

			if err != nil {
				return fmt.Errorf("error creating or opening chat (id:%d): %w", chatID, err)
			}
		}

		// отправляем сообщение в канал на отправку
		sess.Send <- msg
	}
}

func (s *Server) BroadcastMessages(id uint64) error {
	_, sess, err := s.userService.GetUser(id)

	if err != nil {
		return fmt.Errorf("failed to get sender from userService: %w", err)
	}

	for msg := range sess.Send {
		_, sessTo, err := s.userService.GetUser(msg.To)
		if err != nil {
			return fmt.Errorf("failed to get reciever from userServicee: %w", err)
		}
		common.SendMessage(sessTo.Conn, msg)
	}

	return nil
}

func (s *Server) pickNewId() uint64 {
	s.mutex.Lock()
	s.lastId += 1
	newId := s.lastId
	s.mutex.Unlock()

	return newId
}

func (s *Server) notifyUserLeft(id uint64) error {
	user, sess, err := s.userService.GetUser(id)
	if err != nil {
		return nil
	}

	chats, err := s.chats.GetUserChats(id)
	if err != nil {
		return err
	}

	for _, receiverId := range chats {
		sess.Send <- common.Message{
			FromName: user.Name,
			From:     common.ServerId,
			To:       receiverId,
			Type:     common.MessageQuitChat,
			Content:  "",
		}
	}

	return nil
}
