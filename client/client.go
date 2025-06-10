package main

import "net"

type Client struct {
	Conn   net.Conn
	UserId uint64
	Name   string
}

func NewClient(addr string) (*Client, error) {
	conn, err := net.Dial("tcp", addr)
	if err != nil {
		return nil, err
	}

	return &Client{Conn: conn}, nil
}
