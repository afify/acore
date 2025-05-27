package session

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"

	"golang.org/x/crypto/chacha20poly1305"
)

func randomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyz" +
		"ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	result := make([]byte, length)
	max := big.NewInt(int64(len(charset)))
	for i := range result {
		n, err := rand.Int(rand.Reader, max)
		if err != nil {
			return "", fmt.Errorf("RandomString: %w", err)
		}
		result[i] = charset[n.Int64()]
	}
	return string(result), nil
}

func encryptString(plain string) (string, error) {
	key, err := getSessionKey()
	if err != nil {
		return "", err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return "", fmt.Errorf("chacha20poly1305.NewX: %w", err)
	}

	nonce := make([]byte, chacha20poly1305.NonceSizeX)
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce read: %w", err)
	}

	ct := aead.Seal(nonce, nonce, []byte(plain), nil)
	return base64.URLEncoding.EncodeToString(ct), nil
}

func decryptString(ct []byte) ([]byte, error) {
	key, err := getSessionKey()
	if err != nil {
		return nil, err
	}

	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, fmt.Errorf("decryptString: NewX: %w", err)
	}

	nonceSize := chacha20poly1305.NonceSizeX
	if len(ct) < nonceSize {
		return nil, errors.New("decryptString: ciphertext too short")
	}
	nonce, ciphertext := ct[:nonceSize], ct[nonceSize:]

	plain, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryptString: decrypt failed: %w", err)
	}
	return plain, nil
}
