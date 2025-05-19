package models

import (
	"acore/database/db"
	"fmt"
	"log/slog"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type SignUpReq struct {
	UserName string `form:"username"   binding:"required"`
	Email    string `form:"email"      binding:"required"`
	Password string `form:"password"   binding:"required"`
}

type ChangePassReq struct {
	Password        string `form:"password"         binding:"required"`
	ConfirmPassword string `form:"confirm-password" binding:"required"`
}

type SignInReq struct {
	EmailUsername string `form:"email-username" binding:"required"`
	Password      string `form:"password"       binding:"required"`
}

type ForgetReq struct {
	Email string `form:"email"      binding:"required"`
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

func CreateUser(form SignUpReq) (uuid.UUID, error) {
	userID, err := db.CallFuncSingle[uuid.UUID]("create_user", form)
	if err != nil {
		slog.Error("Create User failed", "error", err)
		return uuid.Nil, fmt.Errorf("User.Create: %w", err)
	}

	slog.Info("Create User", "user", userID)
	return *userID, nil
}

func Authenticate(form SignInReq) (uuid.UUID, error) {
	uc, err := db.CallFuncSingle[userCred]("get_user_password_hash", form.EmailUsername)
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
