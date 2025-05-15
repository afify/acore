package database

import (
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func InitDB() {
	var err error

	user := os.Getenv("PG_USER")
	password := os.Getenv("PG_PASSWORD")
	name := os.Getenv("PG_NAME")
	host := os.Getenv("PG_HOST")
	port := os.Getenv("PG_CONT_PORT")

	if user == "" ||
		password == "" ||
		name == "" ||
		host == "" ||
		port == "" {
		log.Fatal("All database environment variables are required")
	}

	connStr := fmt.Sprintf("user=%s password=%s dbname=%s host=%s port=%s sslmode=disable",
		user, password, name, host, port)

	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		log.Printf("Failed to connect to database: %v", err)
	}

	deadline := time.Now().Add(30 * time.Second)
	for {
		if err := DB.Ping(); err == nil {
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
