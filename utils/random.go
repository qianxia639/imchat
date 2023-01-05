package utils

import (
	"fmt"
	"math/rand"
	"strings"
	"time"
)

const alphabet = "abcdefghijklmnopqrstuvwxyz"

func init() {
	rand.Seed(time.Now().UnixNano())
}

// RandomInt 生成 min 到 max 范围的随机数
func RandomInt(min, max int64) int64 {
	return min + rand.Int63n(max-min+1)
}

// RandomString 生成指定长度的随机字符串
func RandomString(n int) string {
	var sb strings.Builder
	k := len(alphabet)

	for i := 0; i < n; i++ {
		c := alphabet[rand.Intn(k)]
		sb.WriteByte(c)
	}

	return sb.String()
}

// Randomggender 生成随机的性别
func RandomGender() int {
	currencies := []int{1, 2, 3}
	n := len(currencies)
	return currencies[rand.Intn(n)]
}

// RandomEmail 生产随机的Email
func RandomEmail() string {
	return fmt.Sprintf("%s@email.com", RandomString(6))
}
