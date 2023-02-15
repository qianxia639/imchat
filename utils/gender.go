package utils

const (
	Man int16 = iota + 1
	Woman
	Unknow
)

func IsSupportedGender(gender int16) bool {
	switch gender {
	case Man, Woman, Unknow:
		return true
	}
	return false
}
