package database

import (
	"database/sql"
	"fmt"
	"os"
	"time"

	_ "github.com/lib/pq"
)

// InitDatabase creates and returns a database connection
func InitDatabase() (*sql.DB, error) {
	// Database connection parameters
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "urlshortener")

	// Create connection string with optimized parameters
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	// Open database connection
	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, fmt.Errorf("error opening database: %w", err)
	}

	// Configure connection pool for high concurrency and scaling
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(10 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	// Test the connection
	if err := db.Ping(); err != nil {
		return nil, fmt.Errorf("error connecting to database: %w", err)
	}

	// Create the table if it doesn't exist
	if err := createTable(db); err != nil {
		return nil, fmt.Errorf("error creating table: %w", err)
	}

	return db, nil
}

// createTable creates the urls table if it doesn't exist
func createTable(db *sql.DB) error {
	query := `
	CREATE TABLE IF NOT EXISTS urls (
		id BIGSERIAL PRIMARY KEY,
		short_code VARCHAR(10) UNIQUE NOT NULL,
		long_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	
	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	`

	_, err := db.Exec(query)
	return err
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
