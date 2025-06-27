package app

import (
	"context"
	"fmt"
	"tcp-chat/common"
	"tcp-chat/transport"
)

type Client struct {
	id       uint64
	name     string
	conn     transport.Connection
	incoming chan common.Message
	outgoing chan common.Message
}

func NewClient(id uint64, name string, conn transport.Connection) *Client {
	return &Client{
		id: id,
		name:     name,
		conn:     conn,
		incoming: make(chan common.Message, 32),
		outgoing: make(chan common.Message, 32),
	}
}

func (c *Client) Start(ctx context.Context, errCh chan error) error {
	go c.readLoop(ctx, errCh)
	go c.writeLoop(ctx, errCh)

	err := common.SendMessage(c.conn, common.Message{
		Type:    common.MessageNameIssuing,
		Content: c.name,
	})

	if err != nil {
		return fmt.Errorf("error issuing name to server: %w", err)
	}

	return err
}

func (c *Client) readLoop(ctx context.Context, errCh chan<- error) {
	for {
		msg, err := common.ReceiveMessage(c.conn)
		if err != nil {
			errCh <- fmt.Errorf("readLoop: error reading from server: %w", err)
			close(c.incoming)
			return
		}

		select {
		case <-ctx.Done():
			return
		case c.incoming <- msg:
		}
	}
}

func (c *Client) writeLoop(ctx context.Context, errCh chan<- error) {
	for {
		select {
		case <-ctx.Done():
			return
		case msg, ok := <-c.outgoing:
			if !ok {
				return
			}

			if err := c.SendMessageTCP(msg.To, msg.Content); err != nil {
				errCh <- fmt.Errorf("writeLoop: error sending to server: %w", err)
				return
			}

		}
	}
}

func (c *Client) SendUserMessage(to uint64, content string) {
	c.outgoing <- common.Message{
		FromName: c.name,
		From:     c.id,
		To:       to,
		Type:     common.MessageText,
		Content:  content,
	}
}

func (c *Client) SendMessageTCP(to uint64, content string) error {
	msg := common.Message{
		FromName: c.name,
		From:     c.id,
		To:       to,
		Type:     common.MessageText,
		Content:  content,
	}

	return common.SendMessage(c.conn, msg)
}

func (c *Client) IncomingMessages() <-chan common.Message {
	return c.incoming
}

func (c *Client) UserID() uint64 {
	return c.id
}

func (c *Client) UserName() string {
	return c.name
}
