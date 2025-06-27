package common

const ServerId = 0

type MessageType int

const (
	MessageText MessageType = iota
	MessageRequestSenderName
	MessageNameIssuing
	MessageIdIssuing
	MessageQuitChat
	MessageSystemNotice
	MessageError
)
