package models

import (
	"acore/database/db"
	"fmt"
	"time"
)

type User struct {
	ID           string
	Username     string
	Email        string
	PasswordHash string
	CreatedAt    time.Time
	UpdatedAt    time.Time
}

func GetByEmail(email string) (*User, error) {
	var u User
	if err := db.CallFunc(&u, "get_user_by_email", email); err != nil {
		return nil, fmt.Errorf("user.GetByEmail: %w", err)
	}
	return &u, nil
}

func GetByID(id string) (*User, error) {
	var u User
	if err := db.CallFunc(&u, "get_user_by_id", id); err != nil {
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return &u, nil
}
