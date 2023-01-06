package validate

import (
	"fmt"
	"net/mail"
	"regexp"
)

var (
	isValidateUsername = regexp.MustCompile(`^[0-9a-zA-Z_]+$`).MatchString
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
	if err := ValidateLen(value, 3, 30); err != nil {
		return err
	}

	if !isValidateUsername(value) {
		return fmt.Errorf("只能包含 字母、数字或下划线")
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
	genders := []int32{1, 2, 3}
	for i := range genders {
		if value == genders[i] {
			return nil
		}
	}
	return fmt.Errorf("非法的数据")
}

func NotEmpty(value string) error {
	if "" == value {
		return fmt.Errorf("not empty")
	}
	return nil
}
