package storage

import (
	"sort"
	"sync"
	"tcp-chat/utils"
)

type MemoryChatRepo struct {
	mu         sync.RWMutex
	nextChatID uint64
	chats      map[uint64][]uint64 // chatID -> участники
	chatIndex  map[string]uint64   // chatKey -> chatID
}

func NewMemoryChatRepo() *MemoryChatRepo {
	return &MemoryChatRepo{
		nextChatID: 1,
		chats:      make(map[uint64][]uint64),
		chatIndex:  make(map[string]uint64),
	}
}

func (r *MemoryChatRepo) CreateChat(userIDs []uint64) (uint64, error) {
	r.mu.Lock()
	defer r.mu.Unlock()

	chatID := r.nextChatID
	r.nextChatID++

	r.chats[chatID] = append([]uint64(nil), userIDs...) // копия слайса

	if len(userIDs) == 2 {
		key := utils.MakeChatKey(userIDs[0], userIDs[1])
		r.chatIndex[key] = chatID
	}

	return chatID, nil
}

func (r *MemoryChatRepo) GetUserChats(userID uint64) ([]uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	var result []uint64
	for chatID, members := range r.chats {
		for _, uid := range members {
			if uid == userID {
				result = append(result, chatID)
				break
			}
		}
	}
	return result, nil
}

func (r *MemoryChatRepo) GetOrCreateChat(userA, userB uint64) (uint64, error) {
	key := utils.MakeChatKey(userA, userB)

	r.mu.RLock()
	if id, ok := r.chatIndex[key]; ok {
		r.mu.RUnlock()
		return id, nil
	}
	r.mu.RUnlock()

	return r.CreateChat([]uint64{userA, userB})
}

func (r *MemoryChatRepo) GetUserCompanions(userID uint64) ([]uint64, error) {
	r.mu.RLock()
	defer r.mu.RUnlock()

	companionSet := make(map[uint64]struct{})

	for _, members := range r.chats {
		isParticipant := false
		for _, uid := range members {
			if uid == userID {
				isParticipant = true
				break
			}
		}
		if isParticipant {
			for _, uid := range members {
				if uid != userID {
					companionSet[uid] = struct{}{}
				}
			}
		}
	}

	var companions []uint64
	for uid := range companionSet {
		companions = append(companions, uid)
	}
	sort.Slice(companions, func(i, j int) bool { return companions[i] < companions[j] })
	return companions, nil
}
