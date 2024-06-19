package convert

import (
	"strconv"
	"strings"
)

// StringToInt 字符串转Int
func StringToInt[T int | int32 | int64](s string) T {
	s = strings.ReplaceAll(s, ",", "")
	var zero T
	switch any(zero).(type) {
	case int:
		i, err := strconv.Atoi(s)
		if err != nil {
			return 0
		}
		return T(i)
	case int32:
		i, err := strconv.ParseInt(s, 10, 32)
		if err != nil {
			return 0
		}
		return T(i)
	case int64:
		i, err := strconv.ParseInt(s, 10, 64)
		if err != nil {
			return 0
		}
		return T(i)
	default:
		return 0
	}
}
