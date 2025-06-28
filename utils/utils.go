package utils

import "fmt"

func Exists(id uint64, set map[uint64]struct{}) bool {
	_, ex := set[id]
	return ex
}

func MakeChatKey(userA, userB uint64) string {
	if userA > userB {
		userA, userB = userB, userA
	}
	return fmt.Sprintf("%d_%d", userA, userB)
}
