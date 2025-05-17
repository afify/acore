package models

import (
	"acore/database/db"
	"fmt"

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
	EmailUsername string `form:"email-username"`
	Password      string `form:"password" binding:"required"`
}

type ForgetReq struct {
	Email string `form:"email"      binding:"required"`
}

type userCred struct {
	ID           string `db:"id"`
	PasswordHash string `db:"password_hash"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(form SignUpReq) (string, error) {
	var newID string

	err := db.CallFunc(&newID,
		"create_user",
		form.UserName,
		form.Email,
		form.Password,
	)

	if err != nil {
		return "", fmt.Errorf("User.Create: %w", err)
	}
	return newID, nil
}
func Authenticate(form SignInReq) (string, error) {
	var uc userCred

	if err := db.CallFunc(&uc,
		"get_user_password_hash",
		form.EmailUsername,
	); err != nil {
		return "", fmt.Errorf("auth.Authenticate (fetch hash): %w", err)
	}

	if !CheckPasswordHash(form.Password, uc.PasswordHash) {
		return "", fmt.Errorf("invalid credentials")
	}

	return uc.ID, nil
}
