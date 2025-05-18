package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"time"
)

type sessionPayload struct {
	UserID string `json:"uid"`
	Exp    int64  `json:"exp"`
	Nonce  string `json:"n"`
}

func GenerateSessionToken(userID string, ttl time.Duration) (string, error) {
	key := []byte(os.Getenv("SESSION_ENC_KEY"))
	if len(key) != 32 {
		return "", errors.New("SESSION_ENC_KEY must be 32 bytes")
	}

	payload := sessionPayload{
		UserID: userID,
		Exp:    time.Now().Add(ttl).Unix(),
	}
	rawNonce := make([]byte, 12)
	if _, err := rand.Read(rawNonce); err != nil {
		return "", fmt.Errorf("GenerateSessionToken: rand.Read: %w", err)
	}
	payload.Nonce = base64.URLEncoding.EncodeToString(rawNonce)

	plain, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("GenerateSessionToken: json.Marshal: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("GenerateSessionToken: NewCipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("GenerateSessionToken: NewGCM: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("GenerateSessionToken: nonce read: %w", err)
	}
	ct := gcm.Seal(nonce, nonce, plain, nil)

	return base64.URLEncoding.EncodeToString(ct), nil
}

func VerifySessionToken(token string) (string, error) {
	var p sessionPayload

	key := []byte(os.Getenv("SESSION_ENC_KEY"))
	if len(key) != 32 {
		return "", errors.New("SESSION_ENC_KEY must be 32 bytes")
	}

	ct, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		return "", fmt.Errorf("VerifySessionToken: decode: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("VerifySessionToken: NewCipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("VerifySessionToken: NewGCM: %w", err)
	}
	if len(ct) < gcm.NonceSize() {
		return "", errors.New("VerifySessionToken: ciphertext too short")
	}
	nonce, ciphertext := ct[:gcm.NonceSize()], ct[gcm.NonceSize():]
	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("VerifySessionToken: decrypt: %w", err)
	}

	if err := json.Unmarshal(plain, &p); err != nil {
		return "", fmt.Errorf("VerifySessionToken: unmarshal: %w", err)
	}
	if time.Now().Unix() > p.Exp {
		return "", errors.New("VerifySessionToken: token expired")
	}
	return p.UserID, nil
}
