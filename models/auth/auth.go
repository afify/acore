package auth

import (
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

type SignupReq struct {
	UserName string `form:"username" validate:"required,alphanum,max=50"`
	Email    string `form:"email"    validate:"required,email,max=50"`
	Password string `form:"password" validate:"required,min=8,max=50"`
}

type LoginReq struct {
	EmailUsername string `form:"email-username" validate:"required"`
	Password      string `form:"password"       validate:"required,min=8"`
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func CheckPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func CreateUser(form SignupReq) (*user.User, error) {
	u, err := dbCreateUser(form)
	if err != nil {
		return nil, err
	}
	return u, nil
}

func Authenticate(form LoginReq) (uuid.UUID, error) {
	u, err := dbGetUserPassword(form)
	if err != nil {
		return uuid.Nil, err
	}

	if !CheckPasswordHash(form.Password, u.PasswordHash) {
		return uuid.Nil, err
	}

	return u.ID, nil
}
