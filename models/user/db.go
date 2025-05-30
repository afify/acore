package user

import (
	"acore/database/db"
	"fmt"
)

func dbGetUserByEmail(email string) (*User, error) {
	u, err := db.CallFuncSingle[User](db.CallFuncParams{
		FuncName: "get_user_by_email",
		FuncArgs: []any{email},
	})

	if err != nil {
		return nil, fmt.Errorf("user.GetByEmail: %w", err)
	}
	return u, nil
}

func dbGetUserById(id string) (*User, error) {
	u, err := db.CallFuncSingle[User](db.CallFuncParams{
		FuncName: "get_user_by_id",
		FuncArgs: []any{id},
	})
	if err != nil {
		return nil, fmt.Errorf("user.GetByID: %w", err)
	}
	return u, nil
}
