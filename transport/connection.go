package transport

import "net"

type Connection interface {
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Close() error
	RemoteAddr() net.Addr
}

type ConnectionListener interface {
	Accept() (Connection, error)
	Close() error
}
