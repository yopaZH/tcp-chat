package common

import (
	"tcp-chat/transport"
)

type ChatConnection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type User struct {
	Id   uint64
	Name string
	Conn transport.Connection

	Send chan Message
}

type Message struct {
	From    uint64
	To      uint64
	Type    MessageType
	Content string
}
