package service

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"crypto/subtle"
	"encoding/base64"
	"fmt"
	"strings"
	"time"

	userEntity "github.com/kiosk404/airi-go/backend/modules/foundation/user/domain/entity"
	"github.com/kiosk404/airi-go/backend/modules/foundation/user/infra/repo/gorm_gen/model"
	"github.com/kiosk404/airi-go/backend/pkg/json"
	"github.com/kiosk404/airi-go/backend/types/consts"
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
type Session struct {
	ID        int64     `json:"id"`         // Session unique device identifier
	CreatedAt time.Time `json:"created_at"` // creation time
	ExpiresAt time.Time `json:"expires_at"` // expiration time
}

// The key used for signing (in practice you should read from the configuration or use environment variables)
var hmacSecret = []byte("openairi-session-hmac-key")

// Generate a secure session key
func generateSessionKey(sessionID int64) (string, error) {
	// Create the default session structure (without the user ID, which will be set in the Login method)
	session := Session{
		ID:        sessionID,
		CreatedAt: time.Now(),
		ExpiresAt: time.Now().Add(consts.DefaultSessionDuration),
	}

	// Serialize session data
	sessionData, err := json.Marshal(session)
	if err != nil {
		return "", err
	}

	// Calculate HMAC signatures to ensure integrity
	h := hmac.New(sha256.New, hmacSecret)
	h.Write(sessionData)
	signature := h.Sum(nil)

	// Combining session data and signatures
	finalData := append(sessionData, signature...)

	// Base64 encoding final result
	return base64.RawURLEncoding.EncodeToString(finalData), nil
}

// Verify the validity of the session key
func verifySessionKey(sessionKey string) (*Session, error) {
	// Decode session data
	data, err := base64.RawURLEncoding.DecodeString(sessionKey)
	if err != nil {
		return nil, fmt.Errorf("invalid session format: %w", err)
	}

	// Make sure the data is long enough to include at least session data and signatures
	if len(data) < 32 { // Simple inspection should actually be more rigorous
		return nil, fmt.Errorf("session data too short")
	}

	// Separating session data and signatures
	sessionData := data[:len(data)-32] // Assume the signature is 32 bytes
	signature := data[len(data)-32:]

	// verify signature
	h := hmac.New(sha256.New, hmacSecret)
	h.Write(sessionData)
	expectedSignature := h.Sum(nil)

	if !hmac.Equal(signature, expectedSignature) {
		return nil, fmt.Errorf("invalid session signature")
	}

	// Parsing session data
	var session Session
	if err := json.Unmarshal(sessionData, &session); err != nil {
		return nil, fmt.Errorf("invalid session data: %w", err)
	}

	// Check if the session has expired
	if time.Now().After(session.ExpiresAt) {
		return nil, fmt.Errorf("session expired")
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
