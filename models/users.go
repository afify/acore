package models

import (
	"acore/database"
	"errors"
	"log"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

type User struct {
	ID           uuid.UUID `json:"id"`
	Username     string    `json:"username"`
	Email        string    `json:"email"`
	PasswordHash string    `json:"-"`
}

func GetAllUsers() ([]User, error) {
	rows, err := database.DB.Query("SELECT id, username, email FROM users")
	if err != nil {
		log.Println("Error fetching users:", err)
		return nil, err
	}
	defer rows.Close()

	users := []User{}
	for rows.Next() {
		var user User
		if err := rows.Scan(&user.ID, &user.Username, &user.Email); err != nil {
			return nil, err
		}
		users = append(users, user)
	}

	return users, nil
}

func GetUserByID(id uuid.UUID) (*User, error) {
	var user User
	err := database.DB.QueryRow("SELECT id, username, email FROM users WHERE id = $1", id).Scan(&user.ID, &user.Username, &user.Email)
	if err != nil {
		log.Println("Error fetching user by ID:", err)
		return nil, err
	}
	return &user, nil
}

func DeleteUserByID(id uuid.UUID) error {
	_, err := database.DB.Exec("DELETE FROM users WHERE id = $1", id)
	if err != nil {
		log.Println("Error deleting user by ID:", err)
		return err
	}

	return nil
}

func CreateUser(user *User) error {
	// Validate that username, email, and password are provided
	if user.Username == "" || user.Email == "" || user.PasswordHash == "" {
		return errors.New("username, email, and password are required")
	}

	// Check if the username or email already exists using the new function
	var usernameExists, emailExists bool
	err := database.DB.QueryRow(
		"SELECT username_exists, email_exists FROM check_user_exists($1, $2)",
		user.Username, user.Email).Scan(&usernameExists, &emailExists)
	if err != nil {
		log.Println("Error checking if username or email exists:", err)
		return err
	}

	if usernameExists {
		return errors.New("username already in use")
	}
	if emailExists {
		return errors.New("email already in use")
	}

	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.PasswordHash), bcrypt.DefaultCost)
	if err != nil {
		log.Println("Error hashing password:", err)
		return err
	}
	user.PasswordHash = string(hashedPassword)

	// Insert the user into the database using a stored procedure
	err = database.DB.QueryRow(
		"CALL create_new_user($1, $2, $3) RETURNING id", // Use your existing stored procedure
		user.Username, user.Email, user.PasswordHash).Scan(&user.ID)

	if err != nil {
		log.Println("Error creating user:", err)
		return err
	}

	return nil
}
