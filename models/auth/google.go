package auth

import (
	"fmt"
	"log/slog"

	"github.com/google/uuid"
)

// GetUserByProvider wraps dbGetUserByProvider, logging and annotating errors.
func GetUserByProvider(provider AuthProvider, sub string) (uuid.UUID, error) {
	userID, err := dbGetUserByProvider(provider, sub)
	if err != nil {
		slog.Warn("GetUserByProvider failed", "Warn", err)
		return uuid.Nil, err
	}
	return userID, nil
}

// LinkProvider wraps dbCreateUserProvider, logging and annotating errors.
func LinkProvider(userID uuid.UUID, provider AuthProvider, sub string) error {
	if err := dbCreateUserProvider(userID, provider, sub); err != nil {
		slog.Error("LinkProvider failed", "error", err)
		return fmt.Errorf("auth.LinkProvider: %w", err)
	}
	return nil
}
