package main

import (
	"log/slog"
	"net/http"

	"acore/database"
	"acore/middleware"
	"acore/routes"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found")
	}
	slog.Info("[0] Logger initialized.")
}

func main() {
	database.InitRedis()
	database.InitDB()
	defer database.CloseDB()

	mux := routes.SetupRoutes()
	handler := middleware.LoggingMiddleware(mux)

	slog.Info("server starting", "port", 8080)
	http.ListenAndServe(":8080", handler)
}
