package common

type ChatConnection interface {
	Read([]byte) (int, error)
	Write([]byte) (int, error)
	Close() error
}

type User struct {
	Id   uint64
	Name string
}

type Message struct {
	FromName string
	From     uint64
	To       uint64
	Type     MessageType
	Content  string
}
