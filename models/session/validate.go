package session

import (
	"errors"
	"fmt"
	"strings"

	"github.com/google/uuid"
)

func parsePayload(plain []byte) (uuid.UUID, error) {
	parts := strings.SplitN(string(plain), ":", 2)
	if len(parts) != 2 {
		return uuid.Nil, errors.New("parsePayload: invalid payload format")
	}
	uid, err := uuid.Parse(parts[0])
	if err != nil {
		return uuid.Nil, fmt.Errorf("parsePayload: invalid userID: %w", err)
	}
	return uid, nil
}

func ValidateSessionToken(token string) (uuid.UUID, uuid.UUID, error) {
	if token == "" {
		return uuid.Nil, uuid.Nil, fmt.Errorf("ValidateSessionToken: no token passed")
	}

	plain, err := decryptString(token)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	userID, err := parsePayload(plain)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	sessionID, err := dbGetUserSessionByID(userID, token)
	if err != nil {
		return uuid.Nil, uuid.Nil, err
	}

	return userID, sessionID, nil
}
