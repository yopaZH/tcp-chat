package app

import (
	"context"
	"tcp-chat/common"
)

type ClientAPI interface {
	Start(ctx context.Context, errCh chan error) error
	SendUserMessage(to uint64, content string)
	IncomingMessages() <-chan common.Message
	UserID() uint64
	UserName() string
}
