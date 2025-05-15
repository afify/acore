package controllers

import (
	"acore/models"
	"encoding/json"
	"net/http"
)

// Example data; in a real application, you would retrieve this from a database
var products = []models.Product{
	{ID: 1, Name: "jacket", Description: "Description for Product 1", Price: 29.99, Image: "jacket"},
	{ID: 2, Name: "Belt", Description: "Description for Product 2", Price: 49.99, Image: "belt"},
	{ID: 3, Name: "Sneakers", Description: "Description for Product 3", Price: 99.99, Image: "sneakers"},
}

// HandleGetProducts handles the GET request for the /products endpoint
func HandleGetProducts(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(products); err != nil {
		http.Error(w, "Failed to encode response", http.StatusInternalServerError)
	}
}
