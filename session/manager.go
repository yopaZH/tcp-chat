package session

import (
	"fmt"
	"sync"
	"tcp-chat/common"
	"tcp-chat/transport"
)

type Manager interface {
	Add(userID uint64, conn transport.Connection) *Session
	Get(userID uint64) (*Session, bool)
	Remove(userID uint64)
	Broadcast(msg common.Message) error
	Shutdown()
}

type Session struct {
	UserID uint64
	Conn   transport.Connection
	Send   chan common.Message
}

// InMemoryManager — управляет активными сессиями
type InMemoryManager struct {
	mu       sync.RWMutex
	sessions map[uint64]*Session
}

// NewInMemoryManager создает новый менеджер
func NewInMemoryManager() *InMemoryManager {
	return &InMemoryManager{
		sessions: make(map[uint64]*Session),
	}
}

// Add добавляет новую сессию (вызывается после подключения)
func (sm *InMemoryManager) Add(userID uint64, conn transport.Connection) *Session {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session := &Session{
		UserID: userID,
		Conn:   conn,
		Send:   make(chan common.Message, 32),
	}
	sm.sessions[userID] = session

	return session
}

// Get возвращает сессию по userID
func (sm *InMemoryManager) Get(userID uint64) (*Session, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, ok := sm.sessions[userID]
	return session, ok
}

// Remove удаляет сессию (и закрывает канал)
func (sm *InMemoryManager) Remove(userID uint64) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	session, ok := sm.sessions[userID]
	if ok {
		close(session.Send)
		session.Conn.Close()
		delete(sm.sessions, userID)
	}
}

// Broadcast отправляет сообщение конкретному пользователю
func (sm *InMemoryManager) Broadcast(msg common.Message) error {
	session, ok := sm.Get(msg.To)
	if !ok {
		return fmt.Errorf("session for user %d not found", msg.To)
	}

	select {
	case session.Send <- msg:
		return nil
	default:
		return fmt.Errorf("send buffer full for user %d", msg.To)
	}
}

// Shutdown останавливает все сессии (например при выходе из сервера)
func (sm *InMemoryManager) Shutdown() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for id, session := range sm.sessions {
		close(session.Send)
		session.Conn.Close()
		delete(sm.sessions, id)
	}
}
