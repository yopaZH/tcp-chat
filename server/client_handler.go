package main

import (
	"bufio"
	"fmt"
	"net"
	"sync"
	"tcp-chat/common"
)

func BroadcastMessages (clients *map[uint64]common.User, broadcast chan common.Message, mutex *sync.Mutex) {
	for msg := range broadcast {
		mutex.Lock()
		//clients[msg.To].Conn.Write([]byte(msg.Content + "\n"))
        common.SendMessage(clients[msg.To].Conn, msg)
        mutex.Unlock()
	}
}

// go HandleClient(conn, &clients, s.broadcast)

func (s *Server) HandleClient(conn net.Conn) {
    defer conn.Close()

    reader := bufio.NewReader(conn)

    newMessage, error := common.ReceiveMessage(conn)

	s.AddUser(name, conn)

	if err != nil {
		mutex.Lock()
		delete(*clients, user.Id)
		mutex.Unlock()

		broadcast <- fmt.Sprintf("%s вышел из чата", name)
		return
    }

}

func getName(u common.User) (string, error) {
	common.SendMessage(u.Conn, common.Message{From: common.ServerId, To: u.Id, Content: "Enter your name"})
	nameMessage, err := common.ReceiveMessage(u.Conn)

	if err != nil {
		return "", fmt.Errorf("couldn't get a name: %w", err)
	}

	return nameMessage.Content, nil
}