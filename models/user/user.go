package user

import (
	"fmt"
	"log/slog"
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID `db:"id"`
	Username  string    `db:"username"`
	Email     string    `db:"email"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

type UserGoogle struct {
	ID        uuid.UUID `db:"id"`
	GoogleSub string    `db:"google_sub"`
	Email     string    `db:"email"`
}

func GetByEmail(email string) (*User, error) {
	u, err := dbGetUserByEmail(email)
	if err != nil {
		slog.Error("GetByEmail", "error", err)
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return u, nil
}

func GetByID(id string) (*User, error) {
	u, err := dbGetUserById(id)
	if err != nil {
		slog.Error("GetByID", "error", err)
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return u, nil
}
