package session

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
	"math/big"
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
	if len(key) != 32 {
		return "", fmt.Errorf("encryptString: session key is %d bytes (want 32)", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("aes.NewCipher: %w", err)
	}

	aead, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("cipher.NewGCM: %w", err)
	}

	nonce := make([]byte, aead.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("nonce read: %w", err)
	}

	ciphertext := aead.Seal(nonce, nonce, []byte(plain), nil)
	return base64.URLEncoding.EncodeToString(ciphertext), nil
}

func decryptString(enc string) ([]byte, error) {
	data, err := base64.URLEncoding.DecodeString(enc)
	if err != nil {
		return nil, fmt.Errorf("base64 decode: %w", err)
	}

	key, err := getSessionKey()
	if err != nil {
		return nil, err
	}
	if len(key) != 32 {
		return nil, fmt.Errorf("decryptString: session key is %d bytes (want 32)", len(key))
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("aes.NewCipher: %w", err)
	}
	aead, err := cipher.NewGCM(block)
	if err != nil {
		return nil, fmt.Errorf("cipher.NewGCM: %w", err)
	}

	nonceSize := aead.NonceSize()
	if len(data) < nonceSize {
		return nil, errors.New("decryptString: ciphertext too short")
	}
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]

	plain, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, fmt.Errorf("decryptString: decrypt failed: %w", err)
	}
	return plain, nil
}
