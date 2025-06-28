package main

import (
	"database/sql"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// URLMapping represents a URL entry in our database
type URLMapping struct {
	ID        int64  `json:"id"`
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
}

// Global database connection
var db *sql.DB

func main() {
	// Parse command line flags
	var benchmarkMode bool
	var stressTest bool
	flag.BoolVar(&benchmarkMode, "benchmark", false, "Run benchmark tests")
	flag.BoolVar(&stressTest, "stress", false, "Enable stress test mode for benchmarks")
	flag.Parse()

	// Set stress test environment variable if flag is provided
	if stressTest {
		os.Setenv("STRESS_TEST", "true")
		fmt.Println("🔥 Stress test mode enabled!")
	}

	if benchmarkMode {
		RunBenchmarks()
		return
	}

	// Initialize database connection
	var err error
	db, err = initDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Successfully connected to database")

	// Initialize Redis connection
	if err := initRedis(); err != nil {
		log.Printf("Warning: Redis initialization failed: %v", err)
	}

	// Initialize batch processor for high-throughput URL shortening
	initBatchProcessor()
	log.Println("✅ Batch processor initialized")

	// Set Gin to release mode for better performance
	gin.SetMode(gin.ReleaseMode)

	// Set up Gin router with optimized settings
	router := gin.New() // Use gin.New() instead of gin.Default() for better performance

	// Add only essential middleware
	router.Use(gin.Recovery())

	// Add custom logging middleware for performance monitoring
	router.Use(func(c *gin.Context) {
		start := time.Now()
		c.Next()
		duration := time.Since(start)
		if duration > 100*time.Millisecond {
			log.Printf("Slow request: %s %s took %v", c.Request.Method, c.Request.URL.Path, duration)
		}
	})

	// Health check endpoint
	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	// URL shortening endpoint
	router.POST("/shorten", shortenURL)

	// URL redirect endpoint
	router.GET("/:shortCode", redirectURL)

	// Start server
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("🚀 Server starting on port %s", port)
	if err := router.Run(":" + port); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}

// initDatabase creates and returns a database connection
func initDatabase() (*sql.DB, error) {
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
