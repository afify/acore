package pg

import (
	"context"
	"fmt"
	"os"
	"time"

	"log/slog"

	"github.com/jackc/pgx/v5"
)

var DB *pgx.Conn

func InitDB() {
	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_NAME")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_CONT_PORT")
	sslMode := os.Getenv("PG_SSL_MODE")

	if user == "" || password == "" || dbName == "" || host == "" ||
		sslMode == "" || port == "" {
		slog.Error("All database environment variables are required")
		os.Exit(1)
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode,
	)

	var err error
	DB, err = pgx.Connect(context.Background(), connStr)
	if err != nil {
		slog.Error("Failed to connect to database", slog.Any("error", err))
	}

	deadline := time.Now().Add(30 * time.Second)
	for {
		if err = DB.Ping(context.Background()); err == nil {
			break
		} else if time.Now().After(deadline) {
			slog.Error("Timed out waiting for Postgres", slog.Any("error", err))
			os.Exit(1)
		} else {
			slog.Info("Waiting for Postgres…", slog.Any("error", err))
			time.Sleep(2 * time.Second)
		}
	}

	slog.Info("[2] Connected to PostgreSQL")
}

func CloseDB() {
	if DB != nil {
		if err := DB.Close(context.Background()); err != nil {
			slog.Error("Failed to close database connection", slog.Any("error", err))
		}
	}
}
