package main

import (
	"fmt"
	"net"
	"tcp-chat/common"
	"tcp-chat/utils"
)

// go HandleClient(conn, &clients, s.broadcast)

func (s *Server) HandleClient(conn net.Conn) {
    defer conn.Close()

	currentId := s.pickNewId()

	s.addUser(currentId, conn)	

	for {
		msg, err := common.ReceiveMessage(conn)

		if err != nil {
			//s.broadcast <- fmt.Sprintf("%s quit chat", name)
			notifyUserLeft(currentId)

			s.mutex.Lock()
			delete(s.clients, user.Id)
			s.mutex.Unlock()

			return
    	} else {
			// добавляем им обоим друг друга в список открытых чатов
			s.mutex.Lock()
			if !utils.exists(s.clients[msg.From].ChatsWith[msg.To]) {
				s.clients[msg.From].ChatsWith[msg.To] = struct{}{}
				s.clients[msg.To].ChatsWith[msg.From] = struct{}{}
			}
			s.mutex.Unlock()

			s.broadcast<- msg
		}
	}
}

func(s *Server) addUser(id uint64, conn net.Conn) {
	nameCh := make(chan NameResult)

	go s.askName(id, nameCh)

	name := <-nameCh

	s.mutex.Lock()
	newUser := common.User{Id: id, Name: name.Name, Conn: conn}
	s.clients[s.lastId] = newUser
	s.mutex.Unlock()
}

type NameResult struct {
	Name string
	Err  error
}

func(s *Server) askName(id uint64, ch chan<- NameResult) {
	s.broadcast<- common.Message{From: common.ServerId, To: id, Type: common.MessageRequestSenderName, Content: ""}
	nameMessage, err := common.ReceiveMessage(u.Conn)

	if err != nil {
		ch <- NameResult{Name: "", Err: fmt.Errorf("couldn't get a name: %w", err)}
	}

	ch <- NameResult{Name: nameMessage.Content, Err: nil}
}

func(s *Server) pickNewId() uint64{
	s.mutex.Lock()
	s.lastId += 1
	newId := s.lastId
	s.mutex.Unlock()

	return newId
}

func(s *Server) notifyUserLeft(id uint64) {
	s.mutex.Rlock() 
	
	for ()
}