package utils

func exists(set map[uint64]struct{}, id uint64) bool{
	_, ex := set[id]
	return ex
}