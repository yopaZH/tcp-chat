package storage

import "tcp-chat/common"

type UserRepository interface {
	AddUser(user common.User) (common.User, error)
	RemoveUser(id uint64) error
	GetUser(id uint64) (*common.User, error)
	Close() error
}
