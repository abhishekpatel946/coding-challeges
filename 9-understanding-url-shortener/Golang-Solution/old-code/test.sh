#!/bin/bash

# URL Shortener Test Script
# Make sure the application is running on localhost:8080

echo "🧪 Testing URL Shortener API"
echo "=============================="

# Test 1: Health Check
echo "1. Testing Health Check..."
curl -s http://localhost:8080/health | jq .
echo ""

# Test 2: Create Short URL
echo "2. Creating Short URL..."
RESPONSE=$(curl -s -X POST http://localhost:8080/shorten \
    -H "Content-Type: application/json" \
    -d '{"long_url": "https://www.google.com"}')

echo $RESPONSE | jq .

# Extract short code from response
SHORT_CODE=$(echo $RESPONSE | jq -r '.short_code')
echo "Short code: $SHORT_CODE"
echo ""

# Test 3: Redirect (just check if it works, don't follow redirect)
echo "3. Testing Redirect..."
curl -s -I http://localhost:8080/$SHORT_CODE | head -5
echo ""

# Test 4: Test with another URL
echo "4. Creating another Short URL..."
RESPONSE2=$(curl -s -X POST http://localhost:8080/shorten \
    -H "Content-Type: application/json" \
    -d '{"long_url": "https://github.com"}')

echo $RESPONSE2 | jq .
echo ""

# Test 5: Test invalid URL
echo "5. Testing Invalid URL..."
curl -s -X POST http://localhost:8080/shorten \
    -H "Content-Type: application/json" \
    -d '{"long_url": "invalid-url"}' | jq .
echo ""

echo "✅ Testing completed!"
