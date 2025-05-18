package pg

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/jackc/pgx/v5"
)

var DB *sql.DB

func InitDB() {
	var err error

	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	dbName := os.Getenv("PG_NAME")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_CONT_PORT")
	sslMode := os.Getenv("PG_SSL_MODE")

	if user == "" || password == "" || dbName == "" || host == "" ||
		sslMode == "" || port == "" {
		log.Fatal("All database environment variables are required")
	}

	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		user, password, host, port, dbName, sslMode,
	)

	DB, err := pgx.Connect(context.Background(), connStr)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	deadline := time.Now().Add(30 * time.Second)
	for {
		if err := DB.Ping(context.Background()); err == nil {
			break
		} else if time.Now().After(deadline) {
			log.Fatalf("Timed out waiting for Postgres: %v", err)
		} else {
			log.Printf("Waiting for Postgres… %v", err)
			time.Sleep(2 * time.Second)
		}
	}
	log.Println("[2] Connected to PostgreSQL")

}

func CloseDB() {
	if DB != nil {
		err := DB.Close()
		if err != nil {
			log.Printf("Failed to close database connection: %v", err)
		}
	}
}
