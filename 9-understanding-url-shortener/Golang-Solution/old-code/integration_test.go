package main

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// TestSetup holds test configuration
type TestSetup struct {
	DB     *sql.DB
	Router *gin.Engine
}

// setupTestDatabase creates a test database connection
func setupTestDatabase() (*sql.DB, error) {
	// Use test database configuration
	connStr := "host=localhost port=5432 user=postgres password=password dbname=urlshortener sslmode=disable"

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Test connection
	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create test table
	createTableQuery := `
	CREATE TABLE IF NOT EXISTS urls (
		id BIGSERIAL PRIMARY KEY,
		short_code VARCHAR(10) UNIQUE NOT NULL,
		long_url TEXT NOT NULL,
		created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
	);
	CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
	`

	_, err = db.Exec(createTableQuery)
	return db, err
}

// setupTestRouter creates a test router with handlers
func setupTestRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Set the global db variable for handlers to use
	globalDB = db

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.POST("/shorten", shortenURL)
	router.GET("/:shortCode", redirectURL)

	return router
}

// cleanupTestData cleans up test data
func cleanupTestData(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM urls")
	return err
}

// TestURLShorteningFlow tests the complete URL shortening flow
func TestURLShorteningFlow(t *testing.T) {
	// Setup
	db, err := setupTestDatabase()
	if err != nil {
		t.Skipf("Database not available: %v", err)
	}
	defer db.Close()

	router := setupTestRouter(db)
	defer cleanupTestData(db)

	// Test data
	testURLs := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
	}

	for _, longURL := range testURLs {
		t.Run(fmt.Sprintf("Test_%s", longURL), func(t *testing.T) {
			// Step 1: Create short URL
			payload := map[string]string{"long_url": longURL}
			jsonData, _ := json.Marshal(payload)

			req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Check response
			if w.Code != http.StatusCreated {
				t.Errorf("Expected status 201, got %d", w.Code)
			}

			var response map[string]interface{}
			if err := json.Unmarshal(w.Body.Bytes(), &response); err != nil {
				t.Fatalf("Failed to parse response: %v", err)
			}

			shortCode, ok := response["short_code"].(string)
			if !ok {
				t.Fatal("Short code not found in response")
			}

			// Step 2: Test redirect
			req = httptest.NewRequest("GET", "/"+shortCode, nil)
			w = httptest.NewRecorder()

			router.ServeHTTP(w, req)

			// Check redirect
			if w.Code != http.StatusFound {
				t.Errorf("Expected status 302, got %d", w.Code)
			}

			location := w.Header().Get("Location")
			if location != longURL {
				t.Errorf("Expected redirect to %s, got %s", longURL, location)
			}
		})
	}
}

// TestConcurrentURLShortening tests concurrent URL shortening
func TestConcurrentURLShortening(t *testing.T) {
	db, err := setupTestDatabase()
	if err != nil {
		t.Skipf("Database not available: %v", err)
	}
	defer db.Close()

	router := setupTestRouter(db)
	defer cleanupTestData(db)

	// Test concurrent requests
	numRequests := 100
	results := make(chan bool, numRequests)

	for i := 0; i < numRequests; i++ {
		go func(id int) {
			longURL := fmt.Sprintf("https://example%d.com", id)
			payload := map[string]string{"long_url": longURL}
			jsonData, _ := json.Marshal(payload)

			req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			results <- w.Code == http.StatusCreated
		}(i)
	}

	// Collect results
	successCount := 0
	for i := 0; i < numRequests; i++ {
		if <-results {
			successCount++
		}
	}

	successRate := float64(successCount) / float64(numRequests)
	if successRate < 0.95 { // 95% success rate threshold
		t.Errorf("Success rate too low: %.2f%%", successRate*100)
	}
}

// TestInvalidURLs tests handling of invalid URLs
func TestInvalidURLs(t *testing.T) {
	db, err := setupTestDatabase()
	if err != nil {
		t.Skipf("Database not available: %v", err)
	}
	defer db.Close()

	router := setupTestRouter(db)
	defer cleanupTestData(db)

	invalidURLs := []string{
		"not-a-url",
		"ftp://example.com",
		"",
		"http://",
		"https://",
	}

	for _, invalidURL := range invalidURLs {
		t.Run(fmt.Sprintf("Invalid_%s", invalidURL), func(t *testing.T) {
			payload := map[string]string{"long_url": invalidURL}
			jsonData, _ := json.Marshal(payload)

			req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			w := httptest.NewRecorder()

			router.ServeHTTP(w, req)

			if w.Code != http.StatusBadRequest {
				t.Errorf("Expected status 400 for invalid URL, got %d", w.Code)
			}
		})
	}
}

// TestDuplicateURLs tests handling of duplicate URLs
func TestDuplicateURLs(t *testing.T) {
	db, err := setupTestDatabase()
	if err != nil {
		t.Skipf("Database not available: %v", err)
	}
	defer db.Close()

	router := setupTestRouter(db)
	defer cleanupTestData(db)

	longURL := "https://www.google.com"

	// Create first short URL
	payload := map[string]string{"long_url": longURL}
	jsonData, _ := json.Marshal(payload)

	req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusCreated {
		t.Fatalf("Failed to create first short URL: %d", w.Code)
	}

	// Create second short URL for same long URL
	req = httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	w = httptest.NewRecorder()

	router.ServeHTTP(w, req)

	// Should still succeed (different short codes for same URL)
	if w.Code != http.StatusCreated {
		t.Errorf("Expected status 201 for duplicate URL, got %d", w.Code)
	}
}

// TestNonExistentShortCode tests handling of non-existent short codes
func TestNonExistentShortCode(t *testing.T) {
	db, err := setupTestDatabase()
	if err != nil {
		t.Skipf("Database not available: %v", err)
	}
	defer db.Close()

	router := setupTestRouter(db)
	defer cleanupTestData(db)

	// Try to access non-existent short code
	req := httptest.NewRequest("GET", "/nonexistent", nil)
	w := httptest.NewRecorder()

	router.ServeHTTP(w, req)

	if w.Code != http.StatusNotFound {
		t.Errorf("Expected status 404 for non-existent short code, got %d", w.Code)
	}
}

// BenchmarkURLShortening benchmarks URL shortening performance
func BenchmarkURLShortening(b *testing.B) {
	db, err := setupTestDatabase()
	if err != nil {
		b.Skipf("Database not available: %v", err)
	}
	defer db.Close()

	router := setupTestRouter(db)
	defer cleanupTestData(db)

	longURL := "https://www.google.com"
	payload := map[string]string{"long_url": longURL}
	jsonData, _ := json.Marshal(payload)

	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		w := httptest.NewRecorder()

		router.ServeHTTP(w, req)

		if w.Code != http.StatusCreated {
			b.Fatalf("Failed to create short URL: %d", w.Code)
		}
	}
}
