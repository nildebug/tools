package utils

import "math/rand"

// RandInt
//
//	@Description: 根据最大值,最小值获取随机数
//	@param min
//	@param max
//	@return int
func RandInt[T int32 | int64 | int](min, max T) int {
	return rand.Intn(int(max-min)) + int(min)
}
