package main

import (
	"context"
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
)

// Global variable for testing - this allows tests to inject their own database connection
var globalDB *sql.DB

// Global Redis client
var redisClient *redis.Client

// URL cache for frequently accessed URLs (in-memory fallback)
var urlCache = struct {
	sync.RWMutex
	data map[string]string
}{
	data: make(map[string]string),
}

// Global counter for short code generation (atomic operations)
var shortCodeCounter int64

// URL batch for processing multiple URLs at once
type URLBatch struct {
	LongURLs []string
	Results  chan BatchResult
}

type BatchResult struct {
	LongURL   string
	ShortCode string
	Error     error
}

// Global batch processor
var batchProcessor *BatchProcessor

// BatchProcessor handles batch URL shortening
type BatchProcessor struct {
	batchChan chan URLBatch
	stopChan  chan bool
}

// NewBatchProcessor creates a new batch processor
func NewBatchProcessor() *BatchProcessor {
	bp := &BatchProcessor{
		batchChan: make(chan URLBatch, 1000),
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
func (bp *BatchProcessor) processBatch(batch URLBatch) {
	// Use globalDB if available (for tests), otherwise use the main db variable
	dbConn := globalDB
	if dbConn == nil {
		dbConn = db
	}

	// Prepare batch insert statement
	stmt, err := dbConn.Prepare("INSERT INTO urls (short_code, long_url) VALUES ($1, $2)")
	if err != nil {
		// Send error for all URLs in batch
		for _, longURL := range batch.LongURLs {
			batch.Results <- BatchResult{LongURL: longURL, Error: err}
		}
		return
	}
	defer stmt.Close()

	// Process each URL in the batch
	for _, longURL := range batch.LongURLs {
		shortCode := generateUniqueShortCode()

		// Try to insert with retry logic
		const maxRetries = 5
		var insertErr error
		for i := 0; i < maxRetries; i++ {
			_, insertErr = stmt.Exec(shortCode, longURL)
			if insertErr == nil {
				break
			}

			// If unique constraint violation, generate new code
			if isUniqueConstraintViolation(insertErr) {
				shortCode = generateUniqueShortCode()
				continue
			}

			break
		}

		batch.Results <- BatchResult{
			LongURL:   longURL,
			ShortCode: shortCode,
			Error:     insertErr,
		}
	}
}

// initBatchProcessor initializes the global batch processor
func initBatchProcessor() {
	batchProcessor = NewBatchProcessor()
}

// initRedis initializes Redis connection
func initRedis() error {
	redisHost := getEnv("REDIS_HOST", "localhost")
	redisPort := getEnv("REDIS_PORT", "6379")

	redisClient = redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%s", redisHost, redisPort),
		Password: "", // no password set
		DB:       0,  // use default DB
		PoolSize: 20, // connection pool size
	})

	// Test connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	_, err := redisClient.Ping(ctx).Result()
	if err != nil {
		log.Printf("Warning: Redis not available, using in-memory cache only: %v", err)
		redisClient = nil
		return nil
	}

	log.Println("✅ Successfully connected to Redis")
	return nil
}

// shortenURL handles POST /shorten requests
func shortenURL(c *gin.Context) {
	// Use globalDB if available (for tests), otherwise use the main db variable
	dbConn := globalDB
	if dbConn == nil {
		dbConn = db
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
	if !isValidURL(request.LongURL) {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid URL format. URL must start with http:// or https:// and be a valid URL.",
		})
		return
	}

	// Check Redis cache first
	if redisClient != nil {
		ctx := context.Background()
		if shortCode, err := redisClient.Get(ctx, "url:"+request.LongURL).Result(); err == nil {
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
	urlCache.RLock()
	if shortCode, exists := urlCache.data[request.LongURL]; exists {
		urlCache.RUnlock()
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
	urlCache.RUnlock()

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
		if isUniqueConstraintViolation(err) {
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
	if redisClient != nil {
		ctx := context.Background()
		redisClient.Set(ctx, "url:"+request.LongURL, shortCode, 24*time.Hour)
		redisClient.Set(ctx, "code:"+shortCode, request.LongURL, 24*time.Hour)
	}

	// Cache the result in memory
	urlCache.Lock()
	urlCache.data[request.LongURL] = shortCode
	urlCache.Unlock()

	shortURL := fmt.Sprintf("http://localhost:8080/%s", shortCode)
	c.JSON(http.StatusCreated, gin.H{
		"short_url":  shortURL,
		"long_url":   request.LongURL,
		"short_code": shortCode,
		"cached":     false,
		"cache_type": "none",
	})
}

// redirectURL handles GET /{shortCode} requests
func redirectURL(c *gin.Context) {
	// Use globalDB if available (for tests), otherwise use the main db variable
	dbConn := globalDB
	if dbConn == nil {
		dbConn = db
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
		shortCode := generateUniqueShortCode()

		_, err = stmt.Exec(shortCode, longURL)
		if err == nil {
			return shortCode, nil
		}

		// If unique constraint violation, try again with a new code
		if isUniqueConstraintViolation(err) {
			// Add a small delay before retrying to reduce collision probability
			time.Sleep(time.Duration(i+1) * time.Millisecond)
			continue
		}

		// For other errors, return immediately
		return "", err
	}

	return "", fmt.Errorf("failed to generate unique short code after %d retries", maxRetries)
}

// generateUniqueShortCode creates a unique short code with better concurrency handling
func generateUniqueShortCode() string {
	// Use atomic increment for thread-safe counter
	counter := atomic.AddInt64(&shortCodeCounter, 1)

	// Get current timestamp in milliseconds
	timestamp := time.Now().UnixNano() / int64(time.Millisecond)

	// Add random component (0-999)
	random := time.Now().UnixNano() % 1000

	// Combine counter, timestamp, and random for maximum uniqueness
	// Use bit shifting to create a more unique number
	combined := (counter << 32) | (timestamp & 0xFFFFFFFF) | (random << 20)

	// Generate base62 string
	base62Str := encodeToBase62(combined)

	// Ensure it's exactly 10 characters
	if len(base62Str) > 10 {
		base62Str = base62Str[:10]
	} else if len(base62Str) < 10 {
		// Pad with random characters to make it exactly 10
		const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
		for len(base62Str) < 10 {
			// Use counter for deterministic but unique padding
			padIndex := (counter + int64(len(base62Str))) % int64(len(charset))
			base62Str += string(charset[padIndex])
		}
	}

	return base62Str
}

// getLongURLOptimized retrieves the long URL with caching
func getLongURLOptimized(dbConn *sql.DB, shortCode string) (string, error) {
	// Check Redis cache first
	if redisClient != nil {
		ctx := context.Background()
		if longURL, err := redisClient.Get(ctx, "code:"+shortCode).Result(); err == nil {
			return longURL, nil
		}
	}

	// Check in-memory cache
	urlCache.RLock()
	for longURL, cachedShortCode := range urlCache.data {
		if cachedShortCode == shortCode {
			urlCache.RUnlock()
			return longURL, nil
		}
	}
	urlCache.RUnlock()

	// Use prepared statement for better performance
	stmt, err := dbConn.Prepare("SELECT long_url FROM urls WHERE short_code = $1")
	if err != nil {
		return "", err
	}
	defer stmt.Close()

	var longURL string
	err = stmt.QueryRow(shortCode).Scan(&longURL)

	// Cache the result if found
	if err == nil && redisClient != nil {
		ctx := context.Background()
		redisClient.Set(ctx, "code:"+shortCode, longURL, 24*time.Hour)
	}

	return longURL, err
}

// encodeToBase62 converts a number to base62 string
func encodeToBase62(num int64) string {
	const charset = "0123456789ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz"
	const base = 62
	if num == 0 {
		return "0"
	}
	var result strings.Builder
	for num > 0 {
		result.WriteByte(charset[num%base])
		num /= base
	}
	// Reverse the string
	bytes := []byte(result.String())
	for i, j := 0, len(bytes)-1; i < j; i, j = i+1, j-1 {
		bytes[i], bytes[j] = bytes[j], bytes[i]
	}
	return string(bytes)
}

// isValidURL performs robust URL validation
func isValidURL(u string) bool {
	if !strings.HasPrefix(u, "http://") && !strings.HasPrefix(u, "https://") {
		return false
	}
	parsed, err := url.ParseRequestURI(u)
	if err != nil {
		return false
	}
	if parsed.Host == "" || parsed.Scheme == "" {
		return false
	}
	// Disallow URLs like http:// or https:// with no host
	if parsed.Host == "" || parsed.Host == "." {
		return false
	}
	return true
}

// isUniqueConstraintViolation checks if an error is a PostgreSQL unique constraint violation
func isUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}
