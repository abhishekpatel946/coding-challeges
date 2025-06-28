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
	_ "github.com/lib/pq"

	"url-shortener/internal/cache"
	"url-shortener/internal/database"
	"url-shortener/internal/handlers"
)

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
		handlers.RunBenchmarks()
		return
	}

	// Initialize database connection
	var err error
	db, err = database.InitDatabase()
	if err != nil {
		log.Fatalf("Failed to connect to database: %v", err)
	}
	defer db.Close()
	log.Println("✅ Successfully connected to database")

	// Initialize Redis connection
	if err := cache.InitRedis(); err != nil {
		log.Printf("Warning: Redis initialization failed: %v", err)
	}

	// Initialize batch processor for high-throughput URL shortening
	handlers.InitBatchProcessor()
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
	router.POST("/shorten", handlers.ShortenURL)

	// URL redirect endpoint
	router.GET("/:shortCode", handlers.RedirectURL)

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
