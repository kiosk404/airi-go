package service

import (
	"context"
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"errors"
	"fmt"
	"strings"
	"time"

	userEntity "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/pkg/logs"
	"golang.org/x/crypto/argon2"
)

// Argon2id parameter
type argon2Params struct {
	memory      uint32
	iterations  uint32
	parallelism uint8
	saltLength  uint32
	keyLength   uint32
}

// Default Argon2id parameters
var defaultArgon2Params = &argon2Params{
	memory:      64 * 1024, // 64MB
	iterations:  3,
	parallelism: 4,
	saltLength:  16,
	keyLength:   32,
}

// Hashing passwords using the Argon id algorithm
func hashPassword(password string) (string, error) {
	p := defaultArgon2Params

	// Generate random salt values
	salt := make([]byte, p.saltLength)
	_, err := rand.Read(salt)
	if err != nil {
		return "", err
	}

	// Calculate the hash value using the Argon id algorithm
	hash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// Encoding to base64 format
	b64Salt := base64.RawStdEncoding.EncodeToString(salt)
	b64Hash := base64.RawStdEncoding.EncodeToString(hash)

	// Format: $argon2id $v = 19 $m = 65536, t = 3, p = 4 $< salt > $< hash >
	encoded := fmt.Sprintf("$argon2id$v=19$m=%d,t=%d,p=%d$%s$%s",
		p.memory, p.iterations, p.parallelism, b64Salt, b64Hash)

	return encoded, nil
}

// Verify that the passwords match
func verifyPassword(password, encodedHash string) (bool, error) {
	// Parse the encoded hash string
	parts := strings.Split(encodedHash, "$")
	if len(parts) != 6 {
		return false, fmt.Errorf("invalid hash format")
	}

	var p argon2Params
	_, err := fmt.Sscanf(parts[3], "m=%d,t=%d,p=%d", &p.memory, &p.iterations, &p.parallelism)
	if err != nil {
		return false, err
	}

	salt, err := base64.RawStdEncoding.DecodeString(parts[4])
	if err != nil {
		return false, err
	}
	p.saltLength = uint32(len(salt))

	decodedHash, err := base64.RawStdEncoding.DecodeString(parts[5])
	if err != nil {
		return false, err
	}
	p.keyLength = uint32(len(decodedHash))

	// Calculate the hash value using the same parameters and salt values
	computedHash := argon2.IDKey(
		[]byte(password),
		salt,
		p.iterations,
		p.memory,
		p.parallelism,
		p.keyLength,
	)

	// Compare the calculated hash value with the stored hash value
	return subtle.ConstantTimeCompare(decodedHash, computedHash) == 1, nil
}

// Session structure, which contains session information

type sessionServiceImpl struct{}

func NewSessionService() ISessionService {
	return &sessionServiceImpl{}
}

func (s sessionServiceImpl) GenerateSessionKey(ctx context.Context, session *userEntity.Session) (string, error) {
	// 设置会话的创建时间和过期时间
	session.CreatedAt = time.Now()
	session.ExpiresAt = time.Now().Add(userEntity.SessionExpires)

	// 序列化会话数据
	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	// 计算HMAC签名以确保完整性
	h := hmac.New(sha256.New, userEntity.HMACSecret)
	h.Write(sessionData)
	signature := h.Sum(nil)

	// 组合会话数据和签名
	finalData := append(sessionData, signature...)

	// Base64编码最终结果
	return base64.RawURLEncoding.EncodeToString(finalData), nil
}

func (s sessionServiceImpl) ValidateSession(ctx context.Context, sessionID string) (*userEntity.Session, error) {
	logs.Debug("sessionID: %s", sessionID)

	// 解码会话数据
	data, err := base64.RawURLEncoding.DecodeString(sessionID)
	if err != nil {
		return nil, fmt.Errorf("invalid session format: %w, data:%s", err, sessionID)
	}

	// 确保数据长够长，至少包含会话数据和签名
	if len(data) < 32 { // 简单检查，实际应该更严格
		return nil, errors.New("session data too short")
	}

	// 分离会话数据和签名
	sessionData := data[:len(data)-32] // 假设签名是32字节
	signature := data[len(data)-32:]

	// 验证签名
	h := hmac.New(sha256.New, userEntity.HMACSecret)
	h.Write(sessionData)
	expectedSignature := h.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return nil, errors.New("invalid session signature")
	}

	// 解析会话数据
	var session userEntity.Session
	if err := json.Unmarshal(sessionData, &session); err != nil {
		return nil, fmt.Errorf("invalid session data: %w", err)
	}

	// 检查会话是否过期
	if time.Now().After(session.ExpiresAt) {
		return nil, errors.New("session expired")
	}

	return &session, nil
}

func userPo2Do(model *model.User, iconURL string) *userEntity.User {
	return &userEntity.User{
		UserID:       model.ID,
		Name:         model.Name,
		UniqueName:   model.UniqueName,
		Account:      model.Account,
		Description:  model.Description,
		IconURI:      model.IconURI,
		IconURL:      iconURL,
		UserVerified: model.UserVerified,
		Locale:       model.Locale,
		SessionKey:   model.SessionKey,
		CreatedAt:    model.CreatedAt,
		UpdatedAt:    model.UpdatedAt,
	}
}
