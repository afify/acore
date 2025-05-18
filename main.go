package main

import (
	"log/slog"
	"net/http"
	"os"

	"acore/database/pg"
	"acore/database/redis"
	"acore/logger"
	"acore/render"
	"acore/routes"

	"github.com/joho/godotenv"
)

func init() {
	err := godotenv.Load()
	if err != nil {
		slog.Info("No .env file found")
	}
	logger.Init()
}

func main() {
	redis.InitRedis()
	pg.InitDB()
	defer pg.CloseDB()

	mux := routes.SetupRoutes()
	render.InitTemplates()

	slog.Info("server starting",
		slog.String("port", os.Getenv("APP_CONTAINER_PORT")))
	http.ListenAndServe(":8080", mux)
}
