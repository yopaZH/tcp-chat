package common

const ServerId = 0

type MessageType int

const (
	MessageText MessageType = iota
	MessageRequestSenderName
	MessageRequestRecieverName
	MessageQuitChat
	MessageSystemNotice
	MessageError
)
