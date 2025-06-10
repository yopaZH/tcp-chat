package utils

func Exists(id uint64, set map[uint64]struct{}) bool{
	_, ex := set[id]
	return ex
}