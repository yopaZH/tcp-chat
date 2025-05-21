package main

import (
	"sync"
	"tcp-chat/common"
)

func broadcastMessages (clients map[uint64]common.User, broadcast chan common.Message, mutex *sync.Mutex) {
	for msg := range broadcast {
		mutex.Lock()
		
		clients[msg.To].Conn.Write([]byte(msg.Content + "\n"))
				
	}
}

/*
func broadcastMessages(clients map[net.Conn]string, broadcast chan string) {
    for msg := range broadcast {
        mutex.Lock()
        for conn := range clients {
            conn.Write([]byte(msg + "\n"))
        }
        mutex.Unlock()
    }
}
*/