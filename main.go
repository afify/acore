package main

import (
	"log/slog"
	"net/http"

	"acore/database/pg"
	"acore/database/redis"
	"acore/middleware"
	"acore/render"
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
	redis.InitRedis()
	pg.InitDB()
	defer pg.CloseDB()

	mux := routes.SetupRoutes()
	render.InitTemplates()
	handler := middleware.LoggingMiddleware(mux)

	slog.Info("server starting", "port", 8080)
	http.ListenAndServe(":8080", handler)
}
