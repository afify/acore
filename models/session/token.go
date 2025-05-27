package session

import (
	"fmt"

	"github.com/google/uuid"
)

func GenerateSessionToken(userID uuid.UUID) (string, error) {
	randPart, err := randomString(128)
	if err != nil {
		return "", err
	}
	payload := fmt.Sprintf("%s:%s", userID, randPart)
	return encryptString(payload)
}
