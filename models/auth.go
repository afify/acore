package models

import (
	"acore/database"
	"crypto/rand"
	"database/sql"
	"encoding/base64"
	"errors"
	"log"

	"golang.org/x/crypto/bcrypt"
)

func ValidateUser(username, password string) bool {
	var user User

	err := database.DB.QueryRow("SELECT id, username, password_hash FROM users WHERE username = $1", username).Scan(&user.ID, &user.Username, &user.PasswordHash)
	if err == sql.ErrNoRows {
		// User not found
		return false
	} else if err != nil {
		// Handle other errors
		log.Println("Error fetching user:", err)
		return false
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password))
	if err != nil {
		return false
	}

	return true
}

func GenerateToken(username string) (string, error) {
	tokenLength := 32
	tokenBytes := make([]byte, tokenLength)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return "", errors.New("unable to generate secure token")
	}
	token := base64.URLEncoding.EncodeToString(tokenBytes)
	return token, nil
}
