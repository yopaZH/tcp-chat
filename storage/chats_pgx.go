package storage

import (
	"context"
	"fmt"
	"sync"
	"tcp-chat/utils"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXChatRepo struct {
	db         *pgxpool.Pool
	chatCache  map[string]uint64
	cacheMutex sync.RWMutex
}

func NewPGXChatRepo(db *pgxpool.Pool) *PGXChatRepo {
	return &PGXChatRepo{
		db:        db,
		chatCache: make(map[string]uint64),
	}
}

// CreateChat — только для приватных чатов из двух пользователей.
func (r *PGXChatRepo) CreateChat(userIDs []uint64) (uint64, error) {
	if len(userIDs) != 2 {
		return 0, fmt.Errorf("only 2 users allowed in private chat")
	}

	userA, userB := userIDs[0], userIDs[1]
	chatKey := utils.MakeChatKey(userA, userB)

	// транзакция
	tx, err := r.db.Begin(context.Background())
	if err != nil {
		return 0, err
	}
	defer tx.Rollback(context.Background())

	var chatID uint64
	err = tx.QueryRow(context.Background(),
		`INSERT INTO chats (is_group, chat_key) VALUES (false, $1) RETURNING id`,
		chatKey).Scan(&chatID)
	if err != nil {
		return 0, err
	}

	_, err = tx.Exec(context.Background(),
		`INSERT INTO chat_members (chat_id, user_id) VALUES ($1, $2), ($1, $3)`,
		chatID, userA, userB)
	if err != nil {
		return 0, err
	}

	if err := tx.Commit(context.Background()); err != nil {
		return 0, err
	}

	// добавляем в кэш
	r.cacheMutex.Lock()
	r.chatCache[chatKey] = chatID
	r.cacheMutex.Unlock()

	return chatID, nil
}

func (r *PGXChatRepo) GetUserChats(userID uint64) ([]uint64, error) {
	rows, err := r.db.Query(context.Background(),
		`SELECT chat_id FROM chat_members WHERE user_id = $1`, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var chatIDs []uint64
	for rows.Next() {
		var chatID uint64
		if err := rows.Scan(&chatID); err != nil {
			return nil, err
		}
		chatIDs = append(chatIDs, chatID)
	}

	return chatIDs, nil
}

func (r *PGXChatRepo) GetUserCompanions(userID uint64) ([]uint64, error) {
	const query = `
		SELECT DISTINCT cm2.user_id
		FROM chat_members cm1
		JOIN chat_members cm2 ON cm1.chat_id = cm2.chat_id
		JOIN chats c ON c.id = cm1.chat_id
		WHERE cm1.user_id = $1 AND cm2.user_id != $1 AND c.is_group = FALSE
	`

	rows, err := r.db.Query(context.Background(), query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var companions []uint64
	for rows.Next() {
		var uid uint64
		if err := rows.Scan(&uid); err != nil {
			return nil, err
		}
		companions = append(companions, uid)
	}

	return companions, nil
}

func (r *PGXChatRepo) GetOrCreateChat(userA, userB uint64) (uint64, error) {
	chatKey := utils.MakeChatKey(userA, userB)

	// проверка в кэше
	r.cacheMutex.RLock()
	if id, ok := r.chatCache[chatKey]; ok {
		r.cacheMutex.RUnlock()
		return id, nil
	}
	r.cacheMutex.RUnlock()

	// чекаем БД
	var chatID uint64
	err := r.db.QueryRow(context.Background(),
		`SELECT id FROM chats WHERE chat_key = $1`, chatKey).Scan(&chatID)
	if err == nil {
		// найден — кешируем
		r.cacheMutex.Lock()
		r.chatCache[chatKey] = chatID
		r.cacheMutex.Unlock()
		return chatID, nil
	}

	// не найден — создаём
	return r.CreateChat([]uint64{userA, userB})
}
