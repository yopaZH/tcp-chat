package common

import "net"

type User struct {
	Id        uint64
	Name      string
	Conn      net.Conn
	ChatsWith map[uint64]struct{}
}

type Message struct {
	From    uint64
	To      uint64
	Type    MessageType
	Content string
}
