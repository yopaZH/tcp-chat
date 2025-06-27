package common

import (
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"tcp-chat/transport"
)

func EncodeMessage(msg Message) ([]byte, error) {
	return json.Marshal(msg)
}

func DecodeMessage(data []byte) (Message, error) {
	var msg Message
	err := json.Unmarshal(data, &msg)
	return msg, err
}

func SendMessage(conn transport.Connection, msg Message) error {
	if conn == nil {
		return fmt.Errorf("conn is nil")
	}

	data, err := EncodeMessage(msg)
	if err != nil {
		return fmt.Errorf("encode error: %w", err)
	}

	length := uint32(len(data))
	if err := binary.Write(conn, binary.BigEndian, length); err != nil {
		return fmt.Errorf("length write error: %w", err)
	}

	_, err = conn.Write(data)
	return err
}

func ReceiveMessage(conn transport.Connection) (Message, error) {
	var length uint32
	if err := binary.Read(conn, binary.BigEndian, &length); err != nil {
		return Message{}, fmt.Errorf("length read error: %w", err)
	}

	data := make([]byte, length)
	if _, err := io.ReadFull(conn, data); err != nil {
		return Message{}, fmt.Errorf("read full error: %w", err)
	}

	return DecodeMessage(data)
}
