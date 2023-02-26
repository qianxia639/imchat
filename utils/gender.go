package utils

// 性别
const (
	Man    int16 = iota + 1 // 男
	Woman                   // 女
	Unknow                  // 未知
)

func IsSupportedGender(gender int16) bool {
	switch gender {
	case Man, Woman, Unknow:
		return true
	}
	return false
}
