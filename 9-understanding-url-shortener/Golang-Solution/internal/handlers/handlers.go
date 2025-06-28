package handlers

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"url-shortener/internal/cache"
	"url-shortener/internal/models"
	"url-shortener/internal/utils"
)

// Global variable for testing - this allows tests to inject their own database connection
var GlobalDB *sql.DB

// Global batch processor
var batchProcessor *BatchProcessor

// BatchProcessor handles batch URL shortening
type BatchProcessor struct {
	batchChan chan models.URLBatch
	stopChan  chan bool
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor() *BatchProcessor {
	bp := &BatchProcessor{
		batchChan: make(chan models.URLBatch, 1000),
		stopChan:  make(chan bool),
	}
	go bp.processBatches()
	return bp
}

// processBatches processes URL batches in the background
func (bp *BatchProcessor) processBatches() {
	for {
		select {
		case batch := <-bp.batchChan:
			bp.processBatch(batch)
		case <-bp.stopChan:
			return
		}
	}
}

// processBatch processes a single batch of URLs
func (bp *BatchProcessor) processBatch(batch models.URLBatch) {
	// Use globalDB if available (for tests), otherwise use the main db variable
	dbConn := GlobalDB
	if dbConn == nil {
		// This will be set by the main application
		return
	}

	// Prepare batch insert statement
	stmt, err := dbConn.Prepare("INSERT INTO urls (short_code, long_url) VALUES ($1, $2)")
	if err != nil {
		// Send error for all URLs in batch
		for _, longURL := range batch.LongURLs {
			batch.Results <- models.BatchResult{LongURL: longURL, Error: err}
		}
		return
	}
	defer stmt.Close()

	// Process each URL in the batch
	for _, longURL := range batch.LongURLs {
		shortCode := utils.GenerateUniqueShortCode()

		// Try to insert with retry logic
		const maxRetries = 5
		var insertErr error
		for i := 0; i < maxRetries; i++ {
			_, insertErr = stmt.Exec(shortCode, longURL)
			if insertErr == nil {
				break
			}

			// If unique constraint violation, generate new code
			if utils.IsUniqueConstraintViolation(insertErr) {
				shortCode = utils.GenerateUniqueShortCode()
				continue
			}

			break
		}

		batch.Results <- models.BatchResult{
			LongURL:   longURL,
			ShortCode: shortCode,
			Error:     insertErr,
		}
	}
}

// InitBatchProcessor initializes the global batch processor
func InitBatchProcessor() {
	batchProcessor = NewBatchProcessor()
}

// ShortenURL handles POST /shorten requests
func ShortenURL(c *gin.Context) {
	// Use globalDB if available (for tests), otherwise use the main db variable
	dbConn := GlobalDB
	if dbConn == nil {
		// This will be set by the main application
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	var request struct {
		LongURL string `json:"long_url" binding:"required"`
	}

	if err := c.ShouldBindJSON(&request); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid request. Please provide a 'long_url' field.",
		})
		return
	}

	// Improved URL validation
	if !utils.IsValidURL(request.LongURL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL format. URL must start with http:// or https:// and be a valid URL.",
		})
		return
	}

	// Check Redis cache first
	if cache.RedisClient != nil {
		ctx := context.Background()
		if shortCode, err := cache.RedisClient.Get(ctx, "url:"+request.LongURL).Result(); err == nil {
			shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)
			c.JSON(http.StatusCreated, gin.H{
				"short_url":  shortURL,
				"long_url":   request.LongURL,
				"short_code": shortCode,
				"cached":     true,
				"cache_type": "redis",
			})
			return
		}
	}

	// Check in-memory cache
	if shortCode, exists := utils.GetFromCache(request.LongURL); exists {
		shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)
		c.JSON(http.StatusCreated, gin.H{
			"short_url":  shortURL,
			"long_url":   request.LongURL,
			"short_code": shortCode,
			"cached":     true,
			"cache_type": "memory",
		})
		return
	}

	// Generate a short code with retry logic for unique constraint
	const maxRetries = 10
	var shortCode string
	var err error
	for i := 0; i < maxRetries; i++ {
		shortCode, err = generateShortCodeOptimized(dbConn, request.LongURL)
		if err == nil {
			break
		}
		// If unique constraint violation, retry
		if utils.IsUniqueConstraintViolation(err) {
			time.Sleep(time.Duration(i+1) * time.Millisecond)
			continue
		}
		log.Printf("Error generating short code: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL"})
		return
	}
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate short URL after retries"})
		return
	}

	// Cache the result in Redis
	if cache.RedisClient != nil {
		ctx := context.Background()
		cache.RedisClient.Set(ctx, "url:"+request.LongURL, shortCode, 24*time.Hour)
		cache.RedisClient.Set(ctx, "code:"+shortCode, request.LongURL, 24*time.Hour)
	}

	// Cache the result in memory
	utils.SetInCache(request.LongURL, shortCode)

	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)
	c.JSON(http.StatusCreated, gin.H{
		"short_url":  shortURL,
		"long_url":   request.LongURL,
		"short_code": shortCode,
		"cached":     false,
		"cache_type": "none",
	})
}

// RedirectURL handles GET /{shortCode} requests
func RedirectURL(c *gin.Context) {
	// Use globalDB if available (for tests), otherwise use the main db variable
	dbConn := GlobalDB
	if dbConn == nil {
		// This will be set by the main application
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database not initialized"})
		return
	}

	shortCode := c.Param("shortCode")
	longURL, err := getLongURLOptimized(dbConn, shortCode)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Short URL not found"})
			return
		}
		log.Printf("Error retrieving long URL: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve URL"})
		return
	}
	c.Redirect(http.StatusFound, longURL)
}

// generateShortCodeOptimized generates a short code with better performance
func generateShortCodeOptimized(dbConn *sql.DB, longURL string) (string, error) {
	// Use a prepared statement for better performance
	stmt, err := dbConn.Prepare("INSERT INTO urls (short_code, long_url) VALUES ($1, $2)")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	// Try to insert with retry logic - increase retries for better collision handling
	const maxRetries = 10
	for i := 0; i < maxRetries; i++ {
		// Generate short code using timestamp + random component for better uniqueness
		shortCode := utils.GenerateUniqueShortCode()

		_, err = stmt.Exec(shortCode, longURL)
		if err == nil {
			return shortCode, nil
		}

		// If unique constraint violation, try again with a new code
		if utils.IsUniqueConstraintViolation(err) {
			// Add a small delay before retrying to reduce collision probability
			time.Sleep(time.Duration(i+1) * time.Millisecond)
			continue
		}

		// For other errors, return immediately
		return "", err
	}

	return "", fmt.Errorf("failed to generate unique short code after %d retries", maxRetries)
}

// getLongURLOptimized retrieves the long URL with caching
func getLongURLOptimized(dbConn *sql.DB, shortCode string) (string, error) {
	// Check Redis cache first
	if cache.RedisClient != nil {
		ctx := context.Background()
		if longURL, err := cache.RedisClient.Get(ctx, "code:"+shortCode).Result(); err == nil {
			return longURL, nil
		}
	}

	// Check in-memory cache
	if longURL, exists := utils.GetFromCacheByValue(shortCode); exists {
		return longURL, nil
	}

	// Use prepared statement for better performance
	stmt, err := dbConn.Prepare("SELECT long_url FROM urls WHERE short_code = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var longURL string
	err = stmt.QueryRow(shortCode).Scan(&longURL)

	// Cache the result if found
	if err == nil && cache.RedisClient != nil {
		ctx := context.Background()
		cache.RedisClient.Set(ctx, "code:"+shortCode, longURL, 24*time.Hour)
	}

	return longURL, err
}
