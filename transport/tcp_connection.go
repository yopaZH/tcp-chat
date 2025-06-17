package transport

import (
	"net"
)

/*
IMPLEMENTS

type ConnectionListener interface {
	Accept() (Connection, error)
	Close() error
}
*/

type tcpListener struct {
	listener net.Listener
}

func NewTCPListener(addr string) (ConnectionListener, error) {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return nil, err
	}
	return &tcpListener{listener: l}, nil
}

func (l *tcpListener) Accept() (Connection, error) {
	conn, err := l.listener.Accept()
	if err != nil {
		return nil, err
	}
	return &tcpConn{conn: conn}, nil
}

func (l *tcpListener) Close() error {
	return l.listener.Close()
}

// обёртка над net.Conn

/*
IMPLEMENTS

type Connection interface {
	Read(p []byte) (int, error)
	Write(p []byte) (int, error)
	Close() error
	RemoteAddr() net.Addr
}
*/

type tcpConn struct {
	conn net.Conn
}

func (c *tcpConn) Read(p []byte) (int, error) {
	return c.conn.Read(p)
}

func (c *tcpConn) Write(p []byte) (int, error) {
	return c.conn.Write(p)
}

func (c *tcpConn) Close() error {
	return c.conn.Close()
}

func (c *tcpConn) RemoteAddr() net.Addr {
	return c.conn.RemoteAddr()
}
