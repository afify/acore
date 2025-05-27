package main

import (
	"log/slog"
	"net/http"
	"os"
	"time"

	"acore/database/pg"
	"acore/database/redis"
	"acore/logger"
	"acore/render"
	"acore/routes"

	"github.com/joho/godotenv"
)

func init() {
	if err := godotenv.Load(".env"); err != nil {
		slog.Error("Loading .env:", "Error", err)
		os.Exit(1)
	}
	logger.Init()
}

func main() {
	redis.InitRedis()
	pg.InitDB()
	defer pg.CloseDB()

	mux := routes.SetupRoutes()
	render.InitTemplates()

	port := os.Getenv("APP_CONTAINER_PORT")
	if port == "" {
		slog.Error("environment variable APP_CONTAINER_PORT not set")
		os.Exit(1)
	}
	addr := ":" + port

	srv := &http.Server{
		Addr:         addr,
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	slog.Info("server starting",
		slog.String("addr", addr),
		slog.String("read_timeout", srv.ReadTimeout.String()),
		slog.String("write_timeout", srv.WriteTimeout.String()),
		slog.String("idle_timeout", srv.IdleTimeout.String()),
	)

	if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		slog.Error("server failed", "error", err)
		os.Exit(1)
	}
}
