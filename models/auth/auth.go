package auth

import (
	"fmt"
	"log/slog"

	"acore/models/user"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type AuthProvider int16

const (
	AuthProviderEmail    AuthProvider = 1
	AuthProviderGoogle   AuthProvider = 2
	AuthProviderApple    AuthProvider = 3
	AuthProviderGithub   AuthProvider = 4
	OauthStateCookieName string       = "oauthstate"
)

func (ap AuthProvider) String() string {
	switch ap {
	case AuthProviderEmail:
		return "email"
	case AuthProviderApple:
		return "apple"
	case AuthProviderGoogle:
		return "google"
	default:
		return "unknown"
	}
}

type SignUpReq struct {
	UserName string `form:"username" validate:"required,alphanum,max=50"`
	Email    string `form:"email"    validate:"required,email,max=50"`
	Password string `form:"password" validate:"required,min=8,max=50"`
}

type SignInReq struct {
	EmailUsername string `form:"email-username" validate:"required"`
	Password      string `form:"password"       validate:"required,min=8"`
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
	slog.Info("Creating User", "form", form)
	u, err := dbCreateUser(form)
	if err != nil {
		slog.Error("Create User failed", "error", err)
		return nil, fmt.Errorf("User.Create: %w", err)
	}

	slog.Info("Create User", "user", u)
	return u, nil
}

func Authenticate(form SignInReq) (uuid.UUID, error) {
	u, err := dbGetUserPassword(form)
	if err != nil {
		slog.Error("Authenticate failed", "error", err)
		return uuid.Nil, fmt.Errorf("auth.Authenticate (fetch hash): %w", err)
	}

	if !CheckPasswordHash(form.Password, u.PasswordHash) {
		slog.Error("Authenticate failed", "error", "invalid credentials")
		return uuid.Nil, fmt.Errorf("invalid credentials")
	}

	slog.Info("Authenticated", "User", u.ID)
	return u.ID, nil
}
