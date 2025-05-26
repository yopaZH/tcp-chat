package common

import (
	"bytes"
	"encoding/binary"
	"encoding/gob"
	"io"
	"net"
)

func init() {
	gob.Register(Message{})
}

func EncodeMessage(msg Message) ([]byte, error) {
	var buf bytes.Buffer
	enc := gob.NewEncoder(&buf)
	if err := enc.Encode(msg); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func DecodeMessage(data []byte) (Message, error) {
	var msg Message
	buf := bytes.NewBuffer(data)
	dec := gob.NewDecoder(buf)
	err := dec.Decode(&msg)
	return msg, err
}

func SendMessage(conn net.Conn, msg Message) error {
	data, err := EncodeMessage(msg)
	if err != nil {
		return err
	}

	length := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return err
	}

	_, err = conn.Write(data)
	return err
}

func ReceiveMessage(conn net.Conn) (Message, error) {
	var msg Message
	var length uint32

	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return msg, err
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return msg, err
	}

	return DecodeMessage(data)
}
