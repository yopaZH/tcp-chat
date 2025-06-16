package storage

import (
	"context"
	"fmt"
	"sync"
	"tcp-chat/common"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXUserRepo struct {
	db     *pgxpool.Pool
	cache  map[uint64]common.User
	mutex  sync.RWMutex
}

// создание нового репозитория с подключением к БД и пустым кэшем
func NewPGXUserRepo(connStr string) (*PGXUserRepo, error) {
	pool, err := pgxpool.New(context.Background(), connStr)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to DB: %w", err)
	}

	return &PGXUserRepo{
		db:    pool,
		cache: make(map[uint64]common.User),
	}, nil
}

func (r *PGXUserRepo) AddUser(user common.User) (common.User, error) {
	// сохраняем в базу (Conn и Send не сохраняются, т.к. это runtime-инфо)
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users (id, name) VALUES ($1, $2) ON CONFLICT (id) DO UPDATE SET name = $2`,
		user.Id, user.Name)
	if err != nil {
		return common.User{}, fmt.Errorf("failed to insert user: %w", err)
	}

	// добавляем в кэш
	r.mutex.Lock()
	r.cache[user.Id] = user
	r.mutex.Unlock()

	return user, nil
}

func (r *PGXUserRepo) RemoveUser(id uint64) error {
	// удаляем из базы
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return fmt.Errorf("failed to delete user: %w", err)
	}

	// удаляем из кэша
	r.mutex.Lock()
	delete(r.cache, id)
	r.mutex.Unlock()

	return nil
}

func (r *PGXUserRepo) GetUser(id uint64) (*common.User, error) {
	// сначала пробуем найти в кэше
	r.mutex.RLock()
	user, ok := r.cache[id]
	r.mutex.RUnlock()
	if ok {
		return &user, nil
	}

	// иначе идём в базу
	row := r.db.QueryRow(context.Background(),
		`SELECT id, name FROM users WHERE id = $1`, id)

	var dbUser common.User
	err := row.Scan(&dbUser.Id, &dbUser.Name)
	if err != nil {
		return nil, fmt.Errorf("user not found in db: %w", err)
	}

	// остальные поля будут нулевыми, т.к. это offline пользователь
	return &dbUser, nil
}

func (r *PGXUserRepo) Close() error {
	r.db.Close()
	return nil
}