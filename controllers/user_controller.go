package controllers

import (
	"encoding/json"
	"net/http"
	"strings"

	"acore/models"

	"github.com/google/uuid"
)

// Handle multiple routes under "/users" (GET, POST)
func HandleUsers(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case "GET":
		GetUsers(w, r)
	case "POST":
		CreateUser(w, r)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// Handle route "/users/{id}" (GET by ID)
func HandleUserByID(w http.ResponseWriter, r *http.Request) {
	// Extract the UUID from the URL path
	idStr := strings.TrimPrefix(r.URL.Path, "/users/")
	id, err := uuid.Parse(idStr) // Parse the string as a UUID
	if err != nil || idStr == "" {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}

	switch r.Method {
	case "GET":
		GetUserByID(w, r, id)
	default:
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
	}
}

// GetUsers returns all users from the database
func GetUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	users, err := models.GetAllUsers()
	if err != nil {
		http.Error(w, "Failed to fetch users", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(users)
}

// GetUserByID returns a specific user by UUID
func GetUserByID(w http.ResponseWriter, r *http.Request, id uuid.UUID) {
	w.Header().Set("Content-Type", "application/json")
	user, err := models.GetUserByID(id)
	if err != nil {
		http.Error(w, "User not found", http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(user)
}

func CreateUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var newUser models.User
	if err := json.NewDecoder(r.Body).Decode(&newUser); err != nil {
		http.Error(w, "Invalid request payload", http.StatusBadRequest)
		return
	}
	if err := models.CreateUser(&newUser); err != nil {
		http.Error(w, "Failed to create user", http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(newUser)
}
