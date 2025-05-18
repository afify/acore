package models

import (
	"acore/database/db"
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID           uuid.UUID `db:"id"`
	Username     string    `db:"username"`
	Email        string    `db:"email"`
	PasswordHash string
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func GetByEmail(email string) (*User, error) {
	var u User
	slog.Info("GetByEmail")
	if err := db.CallFunc(&u, "get_user_by_email", email); err != nil {
		slog.Error("GetByID", slog.Any("error", err))
		return nil, fmt.Errorf("user.GetByEmail: %w", err)
	}
	return &u, nil
}

func GetByID(id string) (*User, error) {
	var u User
	slog.Info("GetByID")
	if err := db.CallFunc(&u, "get_user_by_id", id); err != nil {
		slog.Error("GetByID", slog.Any("error", err))
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return &u, nil
}
