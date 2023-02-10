package validate

import (
	"fmt"
	"net/mail"
	"reflect"
	"regexp"
)

var (
	isValidateUsername = regexp.MustCompile(`^[0-9a-zA-Z_]{4,20}$`).MatchString
	// isValidateNickname = regexp.MustCompile(`^[0-9a-zA-Z_?!\\s]+$`).MatchString
)

func ValidateLen(value string, min, max int) error {
	n := len(value)
	if n < min || n > max {
		return fmt.Errorf("字符长度在 %d-%d 之间", min, max)
	}
	return nil
}

func ValidateUsername(value string) error {
	if !isValidateUsername(value) {
		return fmt.Errorf("只能包含 字母、数字或下划线,且字符长度在4-20之间")
	}
	return nil
}

func ValidateEmail(value string) error {
	if err := ValidateLen(value, 6, 30); err != nil {
		return err
	}

	if _, err := mail.ParseAddress(value); err != nil {
		return fmt.Errorf("无效的邮箱地址")
	}
	return nil
}

func ValidateGender(value int32) error {
	var genders = map[int32]int{
		1: 1,
		2: 2,
		3: 3,
	}
	for range genders {
		if _, ok := genders[value]; !ok {
			return fmt.Errorf("非法的数据")
		}
	}
	return nil
}

func IsEmpty(value any) bool {
	v := reflect.ValueOf(value)

	switch v.Kind() {
	case reflect.String, reflect.Array, reflect.Slice, reflect.Map:
		return v.Len() == 0
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	case reflect.Bool:
		return !v.Bool()
	case reflect.Interface, reflect.Ptr:
		return v.IsNil()
	}
	return reflect.DeepEqual(v.Interface(), reflect.Zero(v.Type()).Interface())
}
