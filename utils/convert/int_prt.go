package convert

import "time"

// GetValuePrt GetValuePrt[T int | int32 | int64 | bool | float32 | float64]
//
//	@Description: 获取基础值的指针
//	@param v
//	@return *T
func GetValuePrt[T int | int32 | int64 | bool | float32 | float64 | uint | time.Time](v T) *T {
	return &v
}
