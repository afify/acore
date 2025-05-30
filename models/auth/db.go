package auth

import (
	"acore/database/db"
	"acore/models/user"

	"github.com/google/uuid"
)

type userCred struct {
	ID           uuid.UUID `db:"id"`
	PasswordHash string    `db:"password_hash"`
}

type providerLink struct {
	UserID uuid.UUID `db:"user_id"`
}

func dbCreateUser(form SignupReq) (*user.User, error) {
	u, err := db.CallFuncSingle[user.User](db.CallFuncParams{
		FuncName: "create_user",
		FuncArgs: []any{form.UserName, form.Email, form.Password},
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func dbGetUserPassword(form LoginReq) (*userCred, error) {
	u, err := db.CallFuncSingle[userCred](db.CallFuncParams{
		FuncName: "get_user_password_hash",
		FuncArgs: []any{form.EmailUsername},
	})
	if err != nil {
		return nil, err
	}
	return u, nil
}

func dbGetUserByProvider(provider AuthProvider, sub string) (uuid.UUID, error) {
	pl, err := db.CallFuncSingle[providerLink](db.CallFuncParams{
		FuncName: "get_user_by_provider",
		FuncArgs: []any{int(provider), sub},
	})
	if err != nil {
		return uuid.Nil, err
	}
	return pl.UserID, nil
}

func dbCreateUserProvider(userID uuid.UUID, provider AuthProvider, sub string) error {
	_, err := db.CallFuncSingle[providerLink](db.CallFuncParams{
		FuncName: "create_user_provider",
		FuncArgs: []any{userID, int(provider), sub},
	})
	if err != nil {
		return err
	}
	return nil
}
