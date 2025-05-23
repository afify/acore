package models

import (
	"acore/database/db"
	"fmt"
	"log/slog"

	user "acore/models/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpReq struct {
	UserName string `form:"username" validate:"required,alphanum",max=50`
	Email    string `form:"email"    validate:"required,email",max=50`
	Password string `form:"password" validate:"required,min=8",max=50`
}

type SignInReq struct {
	EmailUsername string `form:"email-username" validate:"required"`
	Password      string `form:"password"       validate:"required,min=8"`
}

type userCred struct {
	ID           uuid.UUID `db:"id"`
	PasswordHash string    `db:"password_hash"`
}

func HashPassword(password string) (string, error) {
	slog.Info("Hashing Password")
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	slog.Info("CheckPasswordHash")
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(form SignUpReq) (*user.User, error) {
	u, err := db.CallFuncSingle[user.User](db.CallFuncParams{
		FuncName: "create_user",
		FuncArgs: []interface{}{form.UserName, form.Email, form.Password},
	})
	if err != nil {
		slog.Error("Create User failed", "error", err)
		return nil, fmt.Errorf("User.Create: %w", err)
	}

	slog.Info("Create User", "user", u)
	return u, nil
}

func Authenticate(form SignInReq) (uuid.UUID, error) {
	uc, err := db.CallFuncSingle[userCred](db.CallFuncParams{
		FuncName: "get_user_password_hash",
		FuncArgs: []interface{}{form.EmailUsername},
	})

	if err != nil {
		slog.Error("Authenticate failed", "error", err)
		return uuid.Nil, fmt.Errorf("auth.Authenticate (fetch hash): %w", err)
	}

	if !CheckPasswordHash(form.Password, uc.PasswordHash) {
		slog.Error("Authenticate failed", "error", "invalid credentials")
		return uuid.Nil, fmt.Errorf("invalid credentials")
	}

	slog.Info("Authenticated", "User", uc.ID)
	return uc.ID, nil
}
