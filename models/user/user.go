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
	PasswordHash string    `db:"-"`
	CreatedAt    time.Time `db:"created_at"`
	UpdatedAt    time.Time `db:"updated_at"`
}

func GetByEmail(email string) (*User, error) {
	u, err := db.CallFuncSingle[User](db.CallFuncParams{
		FuncName: "get_user_by_email",
		FuncArgs: []interface{}{email},
	})

	if err != nil {
		slog.Error("GetByEmail", "error", err)
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return u, nil
}

func GetByID(id string) (*User, error) {
	u, err := db.CallFuncSingle[User](db.CallFuncParams{
		FuncName: "get_user_by_id",
		FuncArgs: []interface{}{id},
	})
	if err != nil {
		slog.Error("GetByID", "error", err)
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return u, nil
}
