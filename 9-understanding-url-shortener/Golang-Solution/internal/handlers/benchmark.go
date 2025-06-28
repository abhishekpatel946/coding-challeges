package handlers

import (
	"bytes"
	"database/sql"
	"encoding/json"
	"fmt"
	"log"
	"math"
	"net/http"
	"net/http/httptest"
	"os"
	"runtime"
	"sort"
	"strings"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

// Global variables for benchmarks
var serverURL string

// logSystemInfo logs detailed system information
func logSystemInfo() {
	fmt.Println("=== SYSTEM INFORMATION ===")

	// CPU Information
	fmt.Printf("CPU Architecture: %s\n", runtime.GOARCH)
	fmt.Printf("Operating System: %s\n", runtime.GOOS)
	fmt.Printf("Number of CPUs: %d\n", runtime.NumCPU())

	// Memory Information
	var m runtime.MemStats
	runtime.ReadMemStats(&m)
	fmt.Printf("Total Memory: %d MB\n", m.TotalAlloc/1024/1024)
	fmt.Printf("System Memory: %d MB\n", m.Sys/1024/1024)
	fmt.Printf("Heap Memory: %d MB\n", m.HeapAlloc/1024/1024)

	// Goroutine Information
	fmt.Printf("Number of Goroutines: %d\n", runtime.NumGoroutine())

	// Process Information
	pid := os.Getpid()
	fmt.Printf("Process ID: %d\n", pid)

	// Environment Information
	fmt.Printf("Go Version: %s\n", runtime.Version())
	fmt.Printf("GOMAXPROCS: %d\n", runtime.GOMAXPROCS(0))

	fmt.Println("==========================")
	fmt.Println()
}

// BenchmarkResults holds the results of our benchmark tests
type BenchmarkResults struct {
	ConcurrentUsers     int
	RequestsPerSecond   float64
	AverageResponseTime time.Duration
	SuccessRate         float64
	ErrorRate           float64
	TotalRequests       int
	TotalErrors         int
	TotalTime           time.Duration
}

// setupBenchmarkDatabase creates a test database connection for benchmarks
func setupBenchmarkDatabase() (*sql.DB, error) {
	// Use the same database connection parameters as the main application
	host := getEnv("DB_HOST", "localhost")
	port := getEnv("DB_PORT", "5432")
	user := getEnv("DB_USER", "postgres")
	password := getEnv("DB_PASSWORD", "password")
	dbname := getEnv("DB_NAME", "urlshortener")

	// Create connection string
	connStr := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	db, err := sql.Open("postgres", connStr)
	if err != nil {
		return nil, err
	}

	// Configure connection pool for benchmarks - REDUCED for better stability
	db.SetMaxOpenConns(20) // Reduced from 50 to 20
	db.SetMaxIdleConns(5)  // Reduced from 10 to 5
	db.SetConnMaxLifetime(2 * time.Minute)
	db.SetConnMaxIdleTime(30 * time.Second)

	if err := db.Ping(); err != nil {
		return nil, err
	}

	// Create the table if it doesn't exist (same as main application)
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

// setupBenchmarkRouter creates a test router with handlers
func setupBenchmarkRouter(db *sql.DB) *gin.Engine {
	gin.SetMode(gin.TestMode)
	router := gin.New()

	// Set the global db variable for handlers to use
	GlobalDB = db

	router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "healthy"})
	})

	router.POST("/shorten", ShortenURL)
	router.GET("/:shortCode", RedirectURL)

	return router
}

// cleanupBenchmarkData cleans up test data
func cleanupBenchmarkData(db *sql.DB) error {
	_, err := db.Exec("DELETE FROM urls")
	return err
}

// runConcurrencyTest tests how many concurrent users the system can handle
func runConcurrencyTest(db *sql.DB, concurrentUsers int, requestsPerUser int) BenchmarkResults {
	router := setupBenchmarkRouter(db)
	defer cleanupBenchmarkData(db)

	var wg sync.WaitGroup
	results := make(chan bool, concurrentUsers*requestsPerUser)
	startTime := time.Now()

	// Start concurrent users
	for i := 0; i < concurrentUsers; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			for j := 0; j < requestsPerUser; j++ {
				longURL := fmt.Sprintf("https://example%d-user%d.com", userID, j)
				payload := map[string]string{"long_url": longURL}
				jsonData, _ := json.Marshal(payload)

				req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				results <- w.Code == http.StatusCreated
			}
		}(i)
	}

	wg.Wait()
	close(results)
	totalTime := time.Since(startTime)

	// Collect results
	successCount := 0
	totalRequests := 0
	for success := range results {
		totalRequests++
		if success {
			successCount++
		}
	}

	requestsPerSecond := float64(totalRequests) / totalTime.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100
	errorRate := 100 - successRate

	return BenchmarkResults{
		ConcurrentUsers:     concurrentUsers,
		RequestsPerSecond:   requestsPerSecond,
		AverageResponseTime: totalTime / time.Duration(totalRequests),
		SuccessRate:         successRate,
		ErrorRate:           errorRate,
		TotalRequests:       totalRequests,
		TotalErrors:         totalRequests - successCount,
		TotalTime:           totalTime,
	}
}

// runLoadTest simulates sustained load over time
func runLoadTest(db *sql.DB, duration time.Duration, requestsPerSecond int) BenchmarkResults {
	router := setupBenchmarkRouter(db)
	defer cleanupBenchmarkData(db)

	results := make(chan bool, requestsPerSecond*int(duration.Seconds()))
	startTime := time.Now()
	ticker := time.NewTicker(time.Second / time.Duration(requestsPerSecond))
	defer ticker.Stop()

	var wg sync.WaitGroup
	requestID := 0

	// Start load generation
	go func() {
		for range ticker.C {
			if time.Since(startTime) >= duration {
				break
			}
			wg.Add(1)
			go func(id int) {
				defer wg.Done()
				longURL := fmt.Sprintf("https://loadtest%d.example.com", id)
				payload := map[string]string{"long_url": longURL}
				jsonData, _ := json.Marshal(payload)

				req := httptest.NewRequest("POST", "/shorten", bytes.NewBuffer(jsonData))
				req.Header.Set("Content-Type", "application/json")
				w := httptest.NewRecorder()

				router.ServeHTTP(w, req)
				results <- w.Code == http.StatusCreated
			}(requestID)
			requestID++
		}
		wg.Wait()
		close(results)
	}()

	// Wait for completion
	time.Sleep(duration + time.Second)

	// Collect results
	successCount := 0
	totalRequests := 0
	for success := range results {
		totalRequests++
		if success {
			successCount++
		}
	}

	totalTime := time.Since(startTime)
	actualRPS := float64(totalRequests) / totalTime.Seconds()
	successRate := float64(successCount) / float64(totalRequests) * 100
	errorRate := 100 - successRate

	return BenchmarkResults{
		ConcurrentUsers:     requestsPerSecond, // Approximate concurrent users
		RequestsPerSecond:   actualRPS,
		AverageResponseTime: totalTime / time.Duration(totalRequests),
		SuccessRate:         successRate,
		ErrorRate:           errorRate,
		TotalRequests:       totalRequests,
		TotalErrors:         totalRequests - successCount,
		TotalTime:           totalTime,
	}
}

// printBenchmarkReport prints a formatted benchmark report
func printBenchmarkReport(results []BenchmarkResults, testType string) {
	fmt.Printf("\n=== %s BENCHMARK REPORT ===\n", testType)
	fmt.Printf("%-15s %-15s %-15s %-15s %-15s %-15s\n",
		"Concurrent Users", "RPS", "Avg Response", "Success Rate", "Error Rate", "Total Requests")
	fmt.Println(string(make([]byte, 90, 90)))

	for _, result := range results {
		fmt.Printf("%-15d %-15.2f %-15s %-15.2f%% %-15.2f%% %-15d\n",
			result.ConcurrentUsers,
			result.RequestsPerSecond,
			result.AverageResponseTime.String(),
			result.SuccessRate,
			result.ErrorRate,
			result.TotalRequests)
	}
}

// RunBenchmarks runs comprehensive benchmark tests
func RunBenchmarks() {
	fmt.Println("🚀 Starting URL Shortener Benchmark Tests...")
	fmt.Println()

	// Log system information
	logSystemInfo()

	// Initialize database for benchmarks
	db, err := setupBenchmarkDatabase()
	if err != nil {
		log.Fatalf("Failed to setup benchmark database: %v", err)
	}
	defer db.Close()

	// Create a test server for benchmarks
	router := setupBenchmarkRouter(db)
	testServer := httptest.NewServer(router)
	defer testServer.Close()

	// Store the server URL for tests
	serverURL = testServer.URL

	fmt.Println("1. Running Concurrency Tests...")
	concurrencyResults := runConcurrencyTests()
	printBenchmarkReport(concurrencyResults, "CONCURRENCY BENCHMARK REPORT")

	fmt.Println("\n2. Running Latency Tests...")
	latencyResults := runLatencyTests()
	printBenchmarkReport(latencyResults, "LATENCY BENCHMARK REPORT")

	fmt.Println("\n3. Running Capacity Analysis...")
	runCapacityAnalysis()

	fmt.Println("\n4. Running Sustained Load Test...")
	runSustainedLoadTest()

	fmt.Println("\n5. Running Scaling Analysis...")
	runScalingAnalysis()

	fmt.Println("\n6. Running Adaptive Load Test...")
	runAdaptiveLoadTest()

	fmt.Println("\n🎉 Benchmark tests completed!")
}

// runSustainedLoadTest tests sustained load over a longer period
func runSustainedLoadTest() {
	fmt.Println("=== SUSTAINED LOAD TEST ===")
	fmt.Println("Testing sustained load (30 seconds)...")

	start := time.Now()
	totalRequests := 0
	successfulRequests := 0
	errors := 0

	// Run sustained load for 30 seconds
	duration := 30 * time.Second
	endTime := time.Now().Add(duration)

	// Use multiple goroutines for sustained load - REDUCED for better stability
	numWorkers := 5 // Reduced from 10 to 5
	var wg sync.WaitGroup
	var mu sync.Mutex

	for i := 0; i < numWorkers; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for time.Now().Before(endTime) {
				// Generate test URL
				testURL := fmt.Sprintf("https://example.com/test-%d-%d", time.Now().UnixNano(), i)

				// Create request
				jsonData := fmt.Sprintf(`{"long_url": "%s"}`, testURL)
				req, err := http.NewRequest("POST", serverURL+"/shorten", strings.NewReader(jsonData))
				if err != nil {
					mu.Lock()
					errors++
					mu.Unlock()
					continue
				}
				req.Header.Set("Content-Type", "application/json")

				// Send request
				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)

				mu.Lock()
				totalRequests++
				if err == nil && resp.StatusCode == http.StatusCreated {
					successfulRequests++
				} else {
					errors++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}

				// Small delay to prevent overwhelming
				time.Sleep(10 * time.Millisecond)
			}
		}()
	}

	wg.Wait()
	elapsed := time.Since(start)

	// Calculate metrics
	rps := float64(totalRequests) / elapsed.Seconds()
	successRate := float64(successfulRequests) / float64(totalRequests) * 100
	errorRate := float64(errors) / float64(totalRequests) * 100

	fmt.Printf("Sustained Load Results:\n")
	fmt.Printf("Duration: %v\n", elapsed)
	fmt.Printf("Total Requests: %d\n", totalRequests)
	fmt.Printf("Successful Requests: %d\n", successfulRequests)
	fmt.Printf("Errors: %d\n", errors)
	fmt.Printf("Requests/Second: %.2f\n", rps)
	fmt.Printf("Success Rate: %.2f%%\n", successRate)
	fmt.Printf("Error Rate: %.2f%%\n", errorRate)

	// Estimate daily capacity
	dailyCapacity := int64(rps * 86400) // 24 hours in seconds
	fmt.Printf("Estimated Daily Capacity: %d URLs\n", dailyCapacity)

	if successRate >= 99.0 {
		fmt.Println("✅ EXCELLENT - System handles sustained load well")
	} else if successRate >= 95.0 {
		fmt.Println("✅ GOOD - System handles sustained load adequately")
	} else {
		fmt.Println("❌ POOR - System struggles with sustained load")
	}
}

// runScalingAnalysis provides scaling recommendations
func runScalingAnalysis() {
	fmt.Println("=== SCALING ANALYSIS ===")

	// Test different concurrency levels to find optimal scaling
	concurrencyLevels := []int{5, 10, 15}
	var results []struct {
		concurrency int
		rps         float64
		successRate float64
	}

	for _, concurrency := range concurrencyLevels {
		fmt.Printf("Testing %d concurrent users...\n", concurrency)

		start := time.Now()
		totalRequests := 0
		successfulRequests := 0
		errors := 0

		// Run test for each concurrency level
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Make 10 requests per goroutine
				for j := 0; j < 10; j++ {
					testURL := fmt.Sprintf("https://example.com/scaling-test-%d-%d", time.Now().UnixNano(), j)
					jsonData := fmt.Sprintf(`{"long_url": "%s"}`, testURL)

					req, err := http.NewRequest("POST", serverURL+"/shorten", strings.NewReader(jsonData))
					if err != nil {
						mu.Lock()
						errors++
						mu.Unlock()
						continue
					}
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{Timeout: 5 * time.Second}
					resp, err := client.Do(req)

					mu.Lock()
					totalRequests++
					if err == nil && resp.StatusCode == http.StatusCreated {
						successfulRequests++
					} else {
						errors++
					}
					mu.Unlock()

					if resp != nil {
						resp.Body.Close()
					}
				}
			}()
		}

		wg.Wait()
		elapsed := time.Since(start)

		rps := float64(totalRequests) / elapsed.Seconds()
		successRate := float64(successfulRequests) / float64(totalRequests) * 100

		results = append(results, struct {
			concurrency int
			rps         float64
			successRate float64
		}{concurrency, rps, successRate})
	}

	// Find optimal concurrency level
	var optimalConcurrency int
	var maxRPS float64

	for _, result := range results {
		if result.successRate >= 95.0 && result.rps > maxRPS {
			maxRPS = result.rps
			optimalConcurrency = result.concurrency
		}
	}

	fmt.Printf("\nScaling Analysis Results:\n")
	fmt.Printf("Optimal Concurrency Level: %d users\n", optimalConcurrency)
	fmt.Printf("Maximum RPS at Optimal Level: %.2f\n", maxRPS)

	// Calculate scaling requirements
	requiredRPS := 100000000.0 / 86400 // 100M URLs per day
	instancesNeeded := int(math.Ceil(requiredRPS / maxRPS))

	fmt.Printf("Required RPS for 100M URLs/day: %.2f\n", requiredRPS)
	fmt.Printf("Instances needed: %d\n", instancesNeeded)

	if instancesNeeded <= 10 {
		fmt.Println("✅ EXCELLENT - System can scale efficiently")
	} else if instancesNeeded <= 50 {
		fmt.Println("✅ GOOD - System can scale with moderate resources")
	} else {
		fmt.Println("❌ CHALLENGING - System needs significant scaling")
	}
}

// runConcurrencyTests tests different concurrency levels
func runConcurrencyTests() []BenchmarkResults {
	concurrencyLevels := []int{5, 10, 25, 50, 100, 200} // Increased to test higher loads
	var results []BenchmarkResults

	for _, concurrency := range concurrencyLevels {
		fmt.Printf("Testing %d concurrent users...\n", concurrency)

		start := time.Now()
		totalRequests := 0
		successfulRequests := 0
		errors := 0
		var totalResponseTime time.Duration

		// Run test for each concurrency level
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i := 0; i < concurrency; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				// Make 5 requests per goroutine
				for j := 0; j < 5; j++ {
					testURL := fmt.Sprintf("https://example.com/concurrency-test-%d-%d", time.Now().UnixNano(), j)
					jsonData := fmt.Sprintf(`{"long_url": "%s"}`, testURL)

					reqStart := time.Now()
					req, err := http.NewRequest("POST", serverURL+"/shorten", strings.NewReader(jsonData))
					if err != nil {
						mu.Lock()
						errors++
						mu.Unlock()
						continue
					}
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{Timeout: 5 * time.Second}
					resp, err := client.Do(req)
					reqDuration := time.Since(reqStart)

					mu.Lock()
					totalRequests++
					totalResponseTime += reqDuration
					if err == nil && resp.StatusCode == http.StatusCreated {
						successfulRequests++
					} else {
						errors++
					}
					mu.Unlock()

					if resp != nil {
						resp.Body.Close()
					}
				}
			}()
		}

		wg.Wait()
		elapsed := time.Since(start)

		// Calculate metrics
		rps := float64(totalRequests) / elapsed.Seconds()
		avgResponseTime := totalResponseTime / time.Duration(totalRequests)
		successRate := float64(successfulRequests) / float64(totalRequests) * 100
		errorRate := float64(errors) / float64(totalRequests) * 100

		fmt.Printf("%-15d %-15.2f %-15s %-15.2f%% %-15.2f%% %-15d\n",
			concurrency, rps, avgResponseTime, successRate, errorRate, totalRequests)

		results = append(results, BenchmarkResults{
			ConcurrentUsers:     concurrency,
			RequestsPerSecond:   rps,
			AverageResponseTime: avgResponseTime,
			SuccessRate:         successRate,
			ErrorRate:           errorRate,
			TotalRequests:       totalRequests,
			TotalErrors:         errors,
			TotalTime:           elapsed,
		})

		// Stop testing if success rate drops below 90%
		if successRate < 90.0 {
			fmt.Printf("⚠️  Stopping concurrency tests - success rate dropped to %.2f%%\n", successRate)
			break
		}
	}

	return results
}

// runCapacityAnalysis tests different load levels
func runCapacityAnalysis() {
	fmt.Println("=== SYSTEM CAPACITY ANALYSIS ===")

	loadTests := []struct {
		name        string
		users       int
		requestsPer int
	}{
		{"Light Load", 5, 20},
		{"Medium Load", 10, 20},
		{"Heavy Load", 15, 20},
		{"Stress Test", 15, 40},
	}

	for _, test := range loadTests {
		fmt.Printf("\n--- %s (%d users, %d requests each) ---\n", test.name, test.users, test.requestsPer)

		start := time.Now()
		totalRequests := 0
		successfulRequests := 0
		errors := 0
		var totalResponseTime time.Duration

		// Run test
		var wg sync.WaitGroup
		var mu sync.Mutex

		for i := 0; i < test.users; i++ {
			wg.Add(1)
			go func() {
				defer wg.Done()

				for j := 0; j < test.requestsPer; j++ {
					testURL := fmt.Sprintf("https://example.com/capacity-test-%d-%d", time.Now().UnixNano(), j)
					jsonData := fmt.Sprintf(`{"long_url": "%s"}`, testURL)

					reqStart := time.Now()
					req, err := http.NewRequest("POST", serverURL+"/shorten", strings.NewReader(jsonData))
					if err != nil {
						mu.Lock()
						errors++
						mu.Unlock()
						continue
					}
					req.Header.Set("Content-Type", "application/json")

					client := &http.Client{Timeout: 5 * time.Second}
					resp, err := client.Do(req)
					reqDuration := time.Since(reqStart)

					mu.Lock()
					totalRequests++
					totalResponseTime += reqDuration
					if err == nil && resp.StatusCode == http.StatusCreated {
						successfulRequests++
					} else {
						errors++
					}
					mu.Unlock()

					if resp != nil {
						resp.Body.Close()
					}
				}
			}()
		}

		wg.Wait()
		elapsed := time.Since(start)

		// Calculate metrics
		rps := float64(totalRequests) / elapsed.Seconds()
		avgResponseTime := totalResponseTime / time.Duration(totalRequests)
		successRate := float64(successfulRequests) / float64(totalRequests) * 100
		errorRate := float64(errors) / float64(totalRequests) * 100

		fmt.Printf("Requests/Second: %.2f\n", rps)
		fmt.Printf("Average Response Time: %s\n", avgResponseTime)
		fmt.Printf("Success Rate: %.2f%%\n", successRate)
		fmt.Printf("Error Rate: %.2f%%\n", errorRate)

		if successRate >= 99.0 {
			fmt.Println("✅ EXCELLENT - System handles this load perfectly")
		} else if successRate >= 95.0 {
			fmt.Println("✅ GOOD - System handles this load well")
		} else if successRate >= 90.0 {
			fmt.Println("⚠️  ACCEPTABLE - System handles this load adequately")
		} else {
			fmt.Println("❌ POOR - System struggles with this load")
		}
	}
}

// runLatencyTests tests latency at different request rates
func runLatencyTests() []BenchmarkResults {
	requestRates := []int{5, 10, 25, 50, 100, 200} // Increased to test higher throughput
	var results []BenchmarkResults

	for _, rate := range requestRates {
		fmt.Printf("Testing latency at %d requests/second...\n", rate)

		start := time.Now()
		totalRequests := 0
		successfulRequests := 0
		errors := 0
		var totalResponseTime time.Duration

		// Calculate delay between requests
		delay := time.Second / time.Duration(rate)

		// Run test for 6 seconds
		endTime := time.Now().Add(6 * time.Second)
		var wg sync.WaitGroup
		var mu sync.Mutex

		for time.Now().Before(endTime) {
			wg.Add(1)
			go func() {
				defer wg.Done()

				testURL := fmt.Sprintf("https://example.com/latency-test-%d", time.Now().UnixNano())
				jsonData := fmt.Sprintf(`{"long_url": "%s"}`, testURL)

				reqStart := time.Now()
				req, err := http.NewRequest("POST", serverURL+"/shorten", strings.NewReader(jsonData))
				if err != nil {
					mu.Lock()
					errors++
					mu.Unlock()
					return
				}
				req.Header.Set("Content-Type", "application/json")

				client := &http.Client{Timeout: 5 * time.Second}
				resp, err := client.Do(req)
				reqDuration := time.Since(reqStart)

				mu.Lock()
				totalRequests++
				totalResponseTime += reqDuration
				if err == nil && resp.StatusCode == http.StatusCreated {
					successfulRequests++
				} else {
					errors++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}()

			time.Sleep(delay)
		}

		wg.Wait()
		elapsed := time.Since(start)

		// Calculate metrics
		rps := float64(totalRequests) / elapsed.Seconds()
		avgResponseTime := totalResponseTime / time.Duration(totalRequests)
		successRate := float64(successfulRequests) / float64(totalRequests) * 100
		errorRate := float64(errors) / float64(totalRequests) * 100

		results = append(results, BenchmarkResults{
			ConcurrentUsers:     rate,
			RequestsPerSecond:   rps,
			AverageResponseTime: avgResponseTime,
			SuccessRate:         successRate,
			ErrorRate:           errorRate,
			TotalRequests:       totalRequests,
			TotalErrors:         errors,
			TotalTime:           elapsed,
		})

		// Stop testing if success rate drops below 90%
		if successRate < 90.0 {
			fmt.Printf("⚠️  Stopping latency tests - success rate dropped to %.2f%%\n", successRate)
			break
		}
	}

	return results
}

// AdaptiveLoadTest performs dynamic load testing to find the system's breaking point
func runAdaptiveLoadTest() {
	fmt.Println("=== ADAPTIVE LOAD TEST (FINDING BREAKING POINT) ===")
	fmt.Println("Dynamically increasing load until failure detected...")
	fmt.Println()

	// Test parameters - can be adjusted for different scenarios
	initialConcurrency := 5
	maxConcurrency := 1000
	concurrencyStep := 10
	requestsPerUser := 10
	testDuration := 5 * time.Second
	errorThreshold := 5.0               // 5% error rate threshold
	timeoutThreshold := 2 * time.Second // 2 second response time threshold

	// Check if stress test mode is enabled
	stressMode := os.Getenv("STRESS_TEST") == "true"
	if stressMode {
		fmt.Println("🔥 STRESS TEST MODE ENABLED - Testing extreme loads!")
		maxConcurrency = 5000
		concurrencyStep = 50
		requestsPerUser = 20
		testDuration = 3 * time.Second
		errorThreshold = 10.0 // Higher tolerance for stress test
		timeoutThreshold = 5 * time.Second
	}

	var results []AdaptiveTestResult
	var breakingPoint *AdaptiveTestResult

	fmt.Printf("Starting adaptive test from %d to %d concurrent users...\n", initialConcurrency, maxConcurrency)
	fmt.Printf("Error threshold: %.1f%%, Timeout threshold: %v\n", errorThreshold, timeoutThreshold)
	fmt.Printf("Test duration per level: %v\n", testDuration)
	fmt.Println()

	for concurrency := initialConcurrency; concurrency <= maxConcurrency; concurrency += concurrencyStep {
		fmt.Printf("Testing %d concurrent users... ", concurrency)

		// Run test for this concurrency level
		result := runSingleAdaptiveTest(concurrency, requestsPerUser, testDuration)
		results = append(results, result)

		// Check if we've hit the breaking point
		if result.ErrorRate > errorThreshold || result.AvgResponseTime > timeoutThreshold {
			breakingPoint = &result
			fmt.Printf("❌ BREAKING POINT DETECTED!\n")
			fmt.Printf("   Error Rate: %.2f%% (threshold: %.1f%%)\n", result.ErrorRate, errorThreshold)
			fmt.Printf("   Avg Response Time: %v (threshold: %v)\n", result.AvgResponseTime, timeoutThreshold)
			break
		}

		fmt.Printf("✅ PASSED (%.2f%% errors, %v avg response)\n", result.ErrorRate, result.AvgResponseTime)

		// Adaptive step sizing for efficiency
		if concurrency > 100 && concurrencyStep == 10 {
			concurrencyStep = 25
			fmt.Printf("   Increasing step size to %d for faster testing...\n", concurrencyStep)
		}
		if concurrency > 500 && concurrencyStep == 25 {
			concurrencyStep = 50
			fmt.Printf("   Increasing step size to %d for faster testing...\n", concurrencyStep)
		}
		if stressMode && concurrency > 2000 && concurrencyStep == 50 {
			concurrencyStep = 100
			fmt.Printf("   Increasing step size to %d for stress testing...\n", concurrencyStep)
		}
	}

	// Print comprehensive analysis
	printAdaptiveTestAnalysis(results, breakingPoint)
}

// AdaptiveTestResult holds results for a single adaptive test
type AdaptiveTestResult struct {
	Concurrency        int
	TotalRequests      int
	SuccessfulRequests int
	Errors             int
	RPS                float64
	AvgResponseTime    time.Duration
	ErrorRate          float64
	SuccessRate        float64
	MaxResponseTime    time.Duration
	MinResponseTime    time.Duration
	P95ResponseTime    time.Duration
	P99ResponseTime    time.Duration
}

// runSingleAdaptiveTest runs a single test for the adaptive load test
func runSingleAdaptiveTest(concurrency, requestsPerUser int, duration time.Duration) AdaptiveTestResult {
	start := time.Now()
	totalRequests := 0
	successfulRequests := 0
	errors := 0
	var responseTimes []time.Duration
	var mu sync.Mutex

	// Calculate how many requests to make per goroutine
	requestsPerGoroutine := requestsPerUser
	if concurrency*requestsPerGoroutine > 1000 { // Cap total requests
		requestsPerGoroutine = 1000 / concurrency
	}

	var wg sync.WaitGroup
	endTime := time.Now().Add(duration)

	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < requestsPerGoroutine && time.Now().Before(endTime); j++ {
				testURL := fmt.Sprintf("https://example.com/adaptive-test-%d-%d-%d", concurrency, i, j)

				requestStart := time.Now()
				resp, err := http.Post(serverURL+"/shorten", "application/json",
					strings.NewReader(fmt.Sprintf(`{"long_url":"%s"}`, testURL)))
				responseTime := time.Since(requestStart)

				mu.Lock()
				totalRequests++
				responseTimes = append(responseTimes, responseTime)
				if err != nil || resp.StatusCode != http.StatusCreated {
					errors++
				} else {
					successfulRequests++
				}
				mu.Unlock()

				if resp != nil {
					resp.Body.Close()
				}
			}
		}()
	}

	wg.Wait()
	testDuration := time.Since(start)

	// Calculate statistics
	var avgResponseTime time.Duration
	var maxResponseTime time.Duration
	var minResponseTime time.Duration = time.Hour

	if len(responseTimes) > 0 {
		var totalDuration time.Duration
		for _, rt := range responseTimes {
			totalDuration += rt
			if rt > maxResponseTime {
				maxResponseTime = rt
			}
			if rt < minResponseTime {
				minResponseTime = rt
			}
		}
		avgResponseTime = totalDuration / time.Duration(len(responseTimes))
	}

	// Calculate percentiles
	p95ResponseTime := calculatePercentile(responseTimes, 95)
	p99ResponseTime := calculatePercentile(responseTimes, 99)

	rps := float64(totalRequests) / testDuration.Seconds()
	errorRate := float64(errors) / float64(totalRequests) * 100
	successRate := float64(successfulRequests) / float64(totalRequests) * 100

	return AdaptiveTestResult{
		Concurrency:        concurrency,
		TotalRequests:      totalRequests,
		SuccessfulRequests: successfulRequests,
		Errors:             errors,
		RPS:                rps,
		AvgResponseTime:    avgResponseTime,
		ErrorRate:          errorRate,
		SuccessRate:        successRate,
		MaxResponseTime:    maxResponseTime,
		MinResponseTime:    minResponseTime,
		P95ResponseTime:    p95ResponseTime,
		P99ResponseTime:    p99ResponseTime,
	}
}

// calculatePercentile calculates the nth percentile from a slice of durations
func calculatePercentile(times []time.Duration, percentile int) time.Duration {
	if len(times) == 0 {
		return 0
	}

	// Sort the times
	sorted := make([]time.Duration, len(times))
	copy(sorted, times)
	sort.Slice(sorted, func(i, j int) bool {
		return sorted[i] < sorted[j]
	})

	// Calculate index
	index := (percentile * len(sorted)) / 100
	if index >= len(sorted) {
		index = len(sorted) - 1
	}

	return sorted[index]
}

// printAdaptiveTestAnalysis prints comprehensive analysis of adaptive test results
func printAdaptiveTestAnalysis(results []AdaptiveTestResult, breakingPoint *AdaptiveTestResult) {
	fmt.Println("\n" + strings.Repeat("=", 80))
	fmt.Println("ADAPTIVE LOAD TEST ANALYSIS")
	fmt.Println(strings.Repeat("=", 80))

	if breakingPoint != nil {
		fmt.Printf("🚨 BREAKING POINT: %d concurrent users\n", breakingPoint.Concurrency)
		fmt.Printf("   Error Rate: %.2f%%\n", breakingPoint.ErrorRate)
		fmt.Printf("   Avg Response Time: %v\n", breakingPoint.AvgResponseTime)
		fmt.Printf("   RPS: %.2f\n", breakingPoint.RPS)
		fmt.Printf("   Total Requests: %d\n", breakingPoint.TotalRequests)
	} else {
		fmt.Println("✅ NO BREAKING POINT DETECTED - System handled all tested loads!")
	}

	fmt.Println("\n📊 DETAILED RESULTS:")
	fmt.Printf("%-8s %-8s %-12s %-15s %-12s %-12s %-12s\n",
		"Users", "RPS", "Avg Resp", "Error Rate", "P95 Resp", "P99 Resp", "Requests")
	fmt.Println(strings.Repeat("-", 80))

	for _, result := range results {
		fmt.Printf("%-8d %-8.1f %-12v %-15.2f%% %-12v %-12v %-12d\n",
			result.Concurrency,
			result.RPS,
			result.AvgResponseTime,
			result.ErrorRate,
			result.P95ResponseTime,
			result.P99ResponseTime,
			result.TotalRequests)
	}

	// Find optimal performance point
	var optimalResult AdaptiveTestResult
	maxRPS := 0.0
	for _, result := range results {
		if result.RPS > maxRPS && result.ErrorRate < 1.0 {
			maxRPS = result.RPS
			optimalResult = result
		}
	}

	fmt.Printf("\n🎯 OPTIMAL PERFORMANCE: %d concurrent users\n", optimalResult.Concurrency)
	fmt.Printf("   Peak RPS: %.2f\n", optimalResult.RPS)
	fmt.Printf("   Error Rate: %.2f%%\n", optimalResult.ErrorRate)
	fmt.Printf("   Avg Response Time: %v\n", optimalResult.AvgResponseTime)

	// Calculate scaling recommendations
	if breakingPoint != nil {
		safeConcurrency := int(float64(breakingPoint.Concurrency) * 0.8) // 80% of breaking point
		fmt.Printf("\n📈 SCALING RECOMMENDATIONS:\n")
		fmt.Printf("   Safe Production Load: %d concurrent users\n", safeConcurrency)
		fmt.Printf("   Maximum Load: %d concurrent users\n", breakingPoint.Concurrency)
		fmt.Printf("   Recommended Auto-scaling Threshold: %d concurrent users\n", int(float64(breakingPoint.Concurrency)*0.7))
	}

	// Estimate daily capacity
	if len(results) > 0 {
		bestResult := results[len(results)-1]
		if breakingPoint != nil {
			bestResult = *breakingPoint
		}
		dailyCapacity := int64(bestResult.RPS * 86400 * 0.8) // 80% of max capacity
		fmt.Printf("\n📊 CAPACITY ESTIMATES:\n")
		fmt.Printf("   Estimated Daily Capacity: %d URLs\n", dailyCapacity)
		fmt.Printf("   Peak Hourly Capacity: %d URLs\n", int64(bestResult.RPS*3600*0.8))
	}

	fmt.Println(strings.Repeat("=", 80))
}

// getEnv gets an environment variable or returns a default value
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}
