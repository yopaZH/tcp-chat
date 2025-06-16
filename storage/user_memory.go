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

type InMemoryUserRepo struct {
	clients map[uint64]common.User
	mutex   sync.RWMutex
}

func NewInMemoryUserRepo() *InMemoryUserRepo {
	var ms InMemoryUserRepo = InMemoryUserRepo{
		clients: make(map[uint64]common.User)}

	return &ms
}

func (r *InMemoryUserRepo) AddUser(user common.User) (common.User, error) {
	r.mutex.Lock()
	r.clients[user.Id] = user
	r.mutex.Unlock()

	return r.clients[user.Id], nil
}

func (r *InMemoryUserRepo) RemoveUser(id uint64) error {
	_, exists := r.clients[id]

	if exists {
		r.mutex.Lock()
		delete(r.clients, id)
		r.mutex.Unlock()
		return nil
	} else {
		return fmt.Errorf("no user with id:%d to delete", id)
	}
}

func (r *InMemoryUserRepo) GetUser(id uint64) (*common.User, error) {
	r.mutex.RLock()
	user, exists := r.clients[id]
	r.mutex.RUnlock()

	if exists {
		return &user, nil
	} else {
		return nil, fmt.Errorf("no user with id %d found", id)
	}
}

// просто заглушка для соответсвия интерфейсу
func (r *InMemoryUserRepo) Close() error {
	return nil
}
