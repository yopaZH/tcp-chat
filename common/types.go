package common

import "net"

type User struct {
	Id uint64
	Name string
	Conn net.Conn
}

type Message struct {
	From uint64
	To uint64
	Content string
}