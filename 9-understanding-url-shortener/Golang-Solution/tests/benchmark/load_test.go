package benchmark

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"time"
)

// LoadTestConfig holds configuration for load testing
type LoadTestConfig struct {
	BaseURL         string
	ConcurrentUsers int
	RequestsPerUser int
	TestDuration    time.Duration
	URLs            []string
}

// LoadTestResult holds the results of load testing
type LoadTestResult struct {
	TotalRequests       int
	SuccessfulRequests  int
	FailedRequests      int
	AverageResponseTime time.Duration
	MinResponseTime     time.Duration
	MaxResponseTime     time.Duration
	RequestsPerSecond   float64
	Errors              []string
}

// LoadTest performs load testing on the URL shortener
func LoadTest(config LoadTestConfig) LoadTestResult {
	var wg sync.WaitGroup
	var mu sync.Mutex

	result := LoadTestResult{}
	startTime := time.Now()

	// Channel to collect response times
	responseTimes := make(chan time.Duration, config.ConcurrentUsers*config.RequestsPerUser)
	errors := make(chan string, config.ConcurrentUsers*config.RequestsPerUser)

	// Start concurrent users
	for i := 0; i < config.ConcurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()

			for j := 0; j < config.RequestsPerUser; j++ {
				// Randomly select a URL to shorten
				urlIndex := (userID + j) % len(config.URLs)
				testURL := config.URLs[urlIndex]

				// Test URL shortening
				reqStart := time.Now()
				success := testShortenURL(config.BaseURL, testURL)
				responseTime := time.Since(reqStart)

				mu.Lock()
				result.TotalRequests++
				if success {
					result.SuccessfulRequests++
				} else {
					result.FailedRequests++
				}
				mu.Unlock()

				responseTimes <- responseTime

				// Small delay to prevent overwhelming the server
				time.Sleep(10 * time.Millisecond)
			}
		}(i)
	}

	// Wait for all goroutines to complete
	wg.Wait()
	close(responseTimes)
	close(errors)

	// Calculate statistics
	var totalResponseTime time.Duration
	var minResponseTime time.Duration
	var maxResponseTime time.Duration
	first := true

	for responseTime := range responseTimes {
		if first {
			minResponseTime = responseTime
			maxResponseTime = responseTime
			first = false
		} else {
			if responseTime < minResponseTime {
				minResponseTime = responseTime
			}
			if responseTime > maxResponseTime {
				maxResponseTime = responseTime
			}
		}
		totalResponseTime += responseTime
	}

	// Collect errors
	for err := range errors {
		result.Errors = append(result.Errors, err)
	}

	// Calculate final statistics
	actualDuration := time.Since(startTime)
	result.AverageResponseTime = totalResponseTime / time.Duration(result.TotalRequests)
	result.MinResponseTime = minResponseTime
	result.MaxResponseTime = maxResponseTime
	result.RequestsPerSecond = float64(result.TotalRequests) / actualDuration.Seconds()

	return result
}

// testShortenURL tests the URL shortening endpoint
func testShortenURL(baseURL, longURL string) bool {
	payload := map[string]string{
		"long_url": longURL,
	}

	jsonData, err := json.Marshal(payload)
	if err != nil {
		return false
	}

	resp, err := http.Post(baseURL+"/shorten", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		return false
	}
	defer resp.Body.Close()

	return resp.StatusCode == http.StatusCreated
}

// StressTest performs stress testing with gradually increasing load
func StressTest(baseURL string) {
	fmt.Println("🚀 Starting Stress Test...")

	// Test URLs
	testURLs := []string{
		"https://www.google.com",
		"https://www.github.com",
		"https://www.stackoverflow.com",
		"https://www.reddit.com",
		"https://www.youtube.com",
		"https://www.amazon.com",
		"https://www.netflix.com",
		"https://www.twitter.com",
		"https://www.linkedin.com",
		"https://www.medium.com",
	}

	// Gradually increase load
	loadLevels := []struct {
		concurrentUsers int
		requestsPerUser int
		duration        time.Duration
	}{
		{10, 10, 30 * time.Second},   // Light load
		{50, 20, 30 * time.Second},   // Medium load
		{100, 30, 30 * time.Second},  // High load
		{200, 50, 30 * time.Second},  // Very high load
		{500, 100, 30 * time.Second}, // Extreme load
	}

	for i, level := range loadLevels {
		fmt.Printf("\n📊 Load Level %d: %d users, %d requests each\n", i+1, level.concurrentUsers, level.requestsPerUser)

		config := LoadTestConfig{
			BaseURL:         baseURL,
			ConcurrentUsers: level.concurrentUsers,
			RequestsPerUser: level.requestsPerUser,
			TestDuration:    level.duration,
			URLs:            testURLs,
		}

		result := LoadTest(config)
		printLoadTestResult(result)

		// Check if system is still responsive
		if result.RequestsPerSecond < 10 || result.FailedRequests > result.TotalRequests/2 {
			fmt.Println("⚠️  System showing signs of stress, stopping test")
			break
		}

		// Wait between load levels
		time.Sleep(5 * time.Second)
	}
}

// printLoadTestResult prints the results of a load test
func printLoadTestResult(result LoadTestResult) {
	fmt.Printf("📈 Load Test Results:\n")
	fmt.Printf("   Total Requests: %d\n", result.TotalRequests)
	fmt.Printf("   Successful: %d\n", result.SuccessfulRequests)
	fmt.Printf("   Failed: %d\n", result.FailedRequests)
	fmt.Printf("   Success Rate: %.2f%%\n", float64(result.SuccessfulRequests)/float64(result.TotalRequests)*100)
	fmt.Printf("   Requests/Second: %.2f\n", result.RequestsPerSecond)
	fmt.Printf("   Avg Response Time: %v\n", result.AverageResponseTime)
	fmt.Printf("   Min Response Time: %v\n", result.MinResponseTime)
	fmt.Printf("   Max Response Time: %v\n", result.MaxResponseTime)

	if len(result.Errors) > 0 {
		fmt.Printf("   Errors: %d\n", len(result.Errors))
		for i, err := range result.Errors {
			if i < 5 { // Show only first 5 errors
				fmt.Printf("     - %s\n", err)
			}
		}
		if len(result.Errors) > 5 {
			fmt.Printf("     ... and %d more errors\n", len(result.Errors)-5)
		}
	}
}
