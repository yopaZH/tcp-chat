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
			s.notifyUserLeft(currentId)

			s.mutex.Lock()
			delete(s.clients, currentId)
			s.mutex.Unlock()

			return
    	} else {
			// добавляем им друг друга в список открытых чатов
			if (msg.From != common.ServerId || msg.To != common.ServerId) {
				s.updateChatsLists(msg.From, msg.To)
			}

			s.broadcast<- msg
		}
	}
}

func(s *Server) addUser(id uint64, conn net.Conn) error{
	nameCh := make(chan NameResult)
	go s.askName(id, nameCh)
	name := <-nameCh

	if name.Err != nil {
		return fmt.Errorf("error recieving name: %w", name.Err)
	}

	s.mutex.Lock()
	s.clients[s.lastId] = common.User{Id: id, Name: name.Name, Conn: conn}
	s.mutex.Unlock()

	return nil
}

type NameResult struct {
	Name string
	Err  error
}

func(s *Server) askName(id uint64, ch chan<- NameResult) {
	s.broadcast<- common.Message{From: common.ServerId, To: id, Type: common.MessageRequestSenderName, Content: ""}
	nameMessage, err := common.ReceiveMessage(s.clients[id].Conn)

	if err != nil {
		// TODO
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
	s.mutex.Lock() 

	for receiverId := range s.clients[id].ChatsWith {
		s.broadcast <- common.Message{
			From: common.ServerId, 
			To: receiverId, 
			Type: common.MessageQuitChat, 
			Content: s.clients[id].Name}
	}
	
	s.mutex.Unlock()
}

func(s *Server) updateChatsLists(from uint64, to uint64) {
	s.mutex.Lock()
	if !utils.Exists(to, s.clients[from].ChatsWith) {
		s.clients[from].ChatsWith[to] = struct{}{}
		s.clients[to].ChatsWith[from] = struct{}{}
	}
	s.mutex.Unlock()
}