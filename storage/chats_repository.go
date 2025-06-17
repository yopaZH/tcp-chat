package storage

type ChatRepository interface {
	CreateChat(userIDs []uint64) (uint64, error)
	GetUserChats(userID uint64) ([]uint64, error)
	GetUserCompanions(userID uint64) ([]uint64, error)
	GetOrCreateChat(userA, userB uint64) (uint64, error)
}
