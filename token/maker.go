package token

import "time"

// Maker is an interface for managing tokens
type Maker interface {
	// 创建token
	CreateToken(username string, duration time.Duration) (string, error)

	// 校验token
	VerifyToken(token string) (*Payload, error)
}
