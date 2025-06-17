package storage

import (
	"context"
	"sync"
	"tcp-chat/common"

	"github.com/jackc/pgx/v5/pgxpool"
)

type PGXUserRepository struct {
	db    *pgxpool.Pool
	cache map[uint64]common.User // in-memory кэш с соединениями
	mutex sync.RWMutex
}

func NewPGXUserRepository(db *pgxpool.Pool) *PGXUserRepository {
	return &PGXUserRepository{
		db:    db,
		cache: make(map[uint64]common.User),
	}
}

func (r *PGXUserRepository) AddUser(user common.User) (common.User, error) {
	_, err := r.db.Exec(context.Background(),
		`INSERT INTO users (id, name) VALUES ($1, $2)
		 ON CONFLICT (id) DO UPDATE SET name = EXCLUDED.name`,
		user.Id, user.Name)
	if err != nil {
		return common.User{}, err
	}

	// Кладем в кэш с соединением и каналом
	r.mutex.Lock()
	r.cache[user.Id] = user
	r.mutex.Unlock()

	return user, nil
}

func (r *PGXUserRepository) RemoveUser(id uint64) error {
	_, err := r.db.Exec(context.Background(),
		`DELETE FROM users WHERE id = $1`, id)
	if err != nil {
		return err
	}

	r.mutex.Lock()
	delete(r.cache, id)
	r.mutex.Unlock()

	return nil
}

func (r *PGXUserRepository) GetUser(id uint64) (*common.User, error) {
	r.mutex.RLock()
	user, ok := r.cache[id]
	r.mutex.RUnlock()

	if ok {
		return &user, nil
	}

	row := r.db.QueryRow(context.Background(), `SELECT id, name FROM users WHERE id = $1`, id)

	var dbUser common.User
	err := row.Scan(&dbUser.Id, &dbUser.Name)
	if err != nil {
		return nil, err
	}

	return &dbUser, nil
}

func (r *PGXUserRepository) Close() error {
	r.db.Close()
	return nil
}
