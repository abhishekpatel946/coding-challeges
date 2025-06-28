package utils

import (
	"net/url"
	"strings"
	"sync"
	"sync/atomic"
	"time"
)

// Global counter for short code generation (atomic operations)
var ShortCodeCounter int64

// URL cache for frequently accessed URLs (in-memory fallback)
var URLCache = struct {
	sync.RWMutex
	data map[string]string
}{
	data: make(map[string]string),
}

// GetFromCache retrieves a value from the URL cache
func GetFromCache(key string) (string, bool) {
	URLCache.RLock()
	defer URLCache.RUnlock()
	value, exists := URLCache.data[key]
	return value, exists
}

// SetInCache stores a value in the URL cache
func SetInCache(key, value string) {
	URLCache.Lock()
	defer URLCache.Unlock()
	URLCache.data[key] = value
}

// GetFromCacheByValue retrieves a key from the URL cache by searching for a value
func GetFromCacheByValue(value string) (string, bool) {
	URLCache.RLock()
	defer URLCache.RUnlock()
	for key, val := range URLCache.data {
		if val == value {
			return key, true
		}
	}
	return "", false
}

// isValidURL performs robust URL validation
func IsValidURL(u string) bool {
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
func IsUniqueConstraintViolation(err error) bool {
	if err == nil {
		return false
	}
	return strings.Contains(err.Error(), "duplicate key value violates unique constraint")
}

// generateUniqueShortCode creates a unique short code with better concurrency handling
func GenerateUniqueShortCode() string {
	// Use atomic increment for thread-safe counter
	counter := atomic.AddInt64(&ShortCodeCounter, 1)

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
