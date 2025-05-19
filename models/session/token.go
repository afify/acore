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

	"log/slog"

	"github.com/google/uuid"
)

type sessionPayload struct {
	UserID uuid.UUID `json:"uid"`
	Exp    int64     `json:"exp"`
	Nonce  string    `json:"n"`
}

func GenerateSessionToken(userID uuid.UUID, ttl time.Duration) (string, error) {
	key := []byte(os.Getenv("SESSION_ENC_KEY"))
	if len(key) != 32 {
		slog.Error("GenerateSessionToken: invalid SESSION_ENC_KEY length", "length", len(key))
		return "", errors.New("SESSION_ENC_KEY must be 32 bytes")
	}
	slog.Info("GenerateSessionToken: loaded encryption key", "length", len(key))

	payload := sessionPayload{
		UserID: userID,
		Exp:    time.Now().Add(ttl).Unix(),
	}
	slog.Info("GenerateSessionToken: built payload", "userID", userID, "exp", payload.Exp)

	rawNonce := make([]byte, 12)
	if _, err := rand.Read(rawNonce); err != nil {
		slog.Error("GenerateSessionToken: failed to generate raw nonce", "error", err)
		return "", fmt.Errorf("GenerateSessionToken: rand.Read: %w", err)
	}
	payload.Nonce = base64.URLEncoding.EncodeToString(rawNonce)
	slog.Info("GenerateSessionToken: generated payload nonce", "nonce", payload.Nonce)

	plain, err := json.Marshal(payload)
	if err != nil {
		slog.Error("GenerateSessionToken: failed to marshal payload", "error", err)
		return "", fmt.Errorf("GenerateSessionToken: json.Marshal: %w", err)
	}
	slog.Info("GenerateSessionToken: marshaled payload", "size", len(plain))

	block, err := aes.NewCipher(key)
	if err != nil {
		slog.Error("GenerateSessionToken: aes.NewCipher failed", "error", err)
		return "", fmt.Errorf("GenerateSessionToken: NewCipher: %w", err)
	}
	slog.Info("GenerateSessionToken: created AES cipher block")

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("GenerateSessionToken: cipher.NewGCM failed", "error", err)
		return "", fmt.Errorf("GenerateSessionToken: NewGCM: %w", err)
	}
	slog.Info("GenerateSessionToken: initialized GCM", "nonceSize", gcm.NonceSize())

	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		slog.Error("GenerateSessionToken: failed to read GCM nonce", "error", err)
		return "", fmt.Errorf("GenerateSessionToken: nonce read: %w", err)
	}
	slog.Info("GenerateSessionToken: generated GCM nonce")

	ct := gcm.Seal(nonce, nonce, plain, nil)
	slog.Info("GenerateSessionToken: encrypted payload", "ciphertext_length", len(ct))

	token := base64.URLEncoding.EncodeToString(ct)
	slog.Info("GenerateSessionToken: session token generated", "userID", userID, "token_length", len(token))

	return token, nil
}

func VerifySessionToken(token string) (uuid.UUID, error) {
	// 1. Load key
	key := []byte(os.Getenv("SESSION_ENC_KEY"))
	if len(key) != 32 {
		slog.Error("VerifySessionToken: invalid SESSION_ENC_KEY length", "length", len(key))
		return uuid.Nil, errors.New("SESSION_ENC_KEY must be 32 bytes")
	}
	slog.Info("VerifySessionToken: loaded encryption key", "length", len(key))

	ct, err := base64.URLEncoding.DecodeString(token)
	if err != nil {
		slog.Error("VerifySessionToken: failed to decode token", "error", err)
		return uuid.Nil, fmt.Errorf("VerifySessionToken: decode: %w", err)
	}
	slog.Info("VerifySessionToken: decoded token", "ciphertext_length", len(ct))

	block, err := aes.NewCipher(key)
	if err != nil {
		slog.Error("VerifySessionToken: aes.NewCipher failed", "error", err)
		return uuid.Nil, fmt.Errorf("VerifySessionToken: NewCipher: %w", err)
	}
	slog.Info("VerifySessionToken: created AES cipher block")

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		slog.Error("VerifySessionToken: cipher.NewGCM failed", "error", err)
		return uuid.Nil, fmt.Errorf("VerifySessionToken: NewGCM: %w", err)
	}
	slog.Info("VerifySessionToken: initialized GCM", "nonceSize", gcm.NonceSize())

	if len(ct) < gcm.NonceSize() {
		slog.Error("VerifySessionToken: ciphertext too short",
			"length", len(ct), "nonceSize", gcm.NonceSize(),
		)
		return uuid.Nil, errors.New("VerifySessionToken: ciphertext too short")
	}
	nonce, ciphertext := ct[:gcm.NonceSize()], ct[gcm.NonceSize():]
	slog.Info("VerifySessionToken: split nonce and ciphertext", "ciphertext_length", len(ciphertext))

	plain, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		slog.Error("VerifySessionToken: decryption failed", "error", err)
		return uuid.Nil, fmt.Errorf("VerifySessionToken: decrypt: %w", err)
	}
	slog.Info("VerifySessionToken: decrypted payload", "plain_length", len(plain))

	var p sessionPayload
	if err := json.Unmarshal(plain, &p); err != nil {
		slog.Error("VerifySessionToken: failed to unmarshal payload", "error", err)
		return uuid.Nil, fmt.Errorf("VerifySessionToken: unmarshal: %w", err)
	}
	slog.Info("VerifySessionToken: unmarshaled payload", "userID", p.UserID, "exp", p.Exp)

	now := time.Now().Unix()
	if now > p.Exp {
		slog.Error("VerifySessionToken: token expired", "exp", p.Exp, "now", now)
		return uuid.Nil, errors.New("VerifySessionToken: token expired")
	}
	slog.Info("VerifySessionToken: session token verified", "userID", p.UserID)

	return p.UserID, nil
}
