package controllers

import (
	"acore/models"
	"encoding/json"
	"net/http"
)

// Struct for the JSON request body
type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// Struct for the JSON response
type LoginResponse struct {
	Token string `json:"token"`
}

// LoginUser handles the login logic
func LoginUser(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	// Decode the JSON request body
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	// Validate the user using the models function
	if !models.ValidateUser(req.Username, req.Password) {
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)
		return
	}

	// Generate the secure token using the models function
	token, err := models.GenerateToken(req.Username)
	if err != nil {
		http.Error(w, "Unable to generate token", http.StatusInternalServerError)
		return
	}

	// Create the response with the token
	response := LoginResponse{
		Token: token,
	}

	// Return the token as a JSON response
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(response); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
