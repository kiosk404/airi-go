package entity

import (
	"time"
)

const (
	SessionKey     = "session_key"
	SessionExpires = 7 * 24 * time.Hour
)

// HMACSecret 用于签名的密钥（在实际应用中应从配置中读取或使用环境变量）
var HMACSecret = []byte("airi-go-session-hmac-key")

type Session struct {
	UserID    int64
	SessionID int64
	Locale    string
	CreatedAt time.Time
	ExpiresAt time.Time
}
