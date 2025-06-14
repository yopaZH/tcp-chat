package storage

import (
	"fmt"
	"sync"
	"tcp-chat/common"
)

/*
type UserRepository interface {
	AddUser(user common.User) error
	RemoveUser(id uint64) error
	GetUser(id uint64) (common.User, error)
}
*/

type MemoryStorage struct {
	clients map[uint64]common.User
	mutex   sync.RWMutex
}

func NewMemoryStorage() *MemoryStorage {
	var ms MemoryStorage = MemoryStorage{
		clients: make(map[uint64]common.User)}

	return &ms
}

func (ms *MemoryStorage) AddUser(user common.User) error {
	ms.mutex.Lock()
	ms.clients[user.Id] = user
	ms.mutex.Unlock()

	return nil
}

func (ms *MemoryStorage) RemoveUser(id uint64) error {
	_, exists := ms.clients[id]

	if exists {
		ms.mutex.Lock()
		delete(ms.clients, id)
		ms.mutex.Unlock()
		return nil
	} else {
		return fmt.Errorf("no user with id:%d to delete", id)
	}
}

func (ms *MemoryStorage) GetUser(id uint64) (*common.User, error) {
	ms.mutex.RLock()
	user, exists := ms.clients[id]
	ms.mutex.RUnlock()

	if exists {
		return &user, nil
	} else {
		return nil, fmt.Errorf("no user with index %d found", id)
	}
}

// просто заглушка для соответсвия интерфейсу
func (ms *MemoryStorage) Close() error {
	return nil
}
