# URL Shortener - Go Learning Project

A scalable URL shortener built in Go with PostgreSQL, Redis, and comprehensive testing. This project demonstrates modern Go development practices, system design concepts, and performance optimization techniques.

## 🏗️ Project Structure

```
url-shortener/
├── cmd/
│   └── server/           # Application entry point
├── internal/             # Private application code
│   ├── handlers/         # HTTP request handlers
│   ├── models/           # Data models and structures
│   ├── database/         # Database operations
│   ├── cache/            # Caching layer (Redis)
│   └── utils/            # Utility functions
├── pkg/                  # Public packages
│   ├── config/           # Configuration utilities
│   └── logger/           # Logging utilities
├── tests/                # Test files
│   ├── unit/             # Unit tests
│   ├── integration/      # Integration tests
│   └── benchmark/        # Performance tests
├── scripts/              # Build and deployment scripts
├── deployments/          # Deployment configurations
├── docs/                 # Documentation
├── go.mod
├── go.sum
└── README.md
```

## 🚀 Quick Start

### 1. Start the Database

```bash
# Start PostgreSQL and Redis using Docker Compose
cd deployments
docker-compose up -d db redis

# Check if services are running
docker-compose ps
```

### 2. Run the Application

```bash
# Install dependencies
go mod tidy

# Run the application
go run cmd/server/main.go
```

The server will start on `http://localhost:8080`

## 🧪 Testing the API

### 1. Health Check

```bash
curl http://localhost:8080/health
```

### 2. Create a Short URL

```bash
curl -X POST http://localhost:8080/shorten \
  -H "Content-Type: application/json" \
  -d '{"long_url": "https://www.google.com"}'
```

### 3. Redirect to Original URL

```bash
# Replace {shortCode} with the actual short code from step 2
curl -L http://localhost:8080/{shortCode}
```

## 🔧 Environment Variables

The application uses these environment variables (with defaults):

- `DB_HOST`: localhost
- `DB_PORT`: 5432
- `DB_USER`: postgres
- `DB_PASSWORD`: password
- `DB_NAME`: urlshortener
- `REDIS_HOST`: localhost
- `REDIS_PORT`: 6379
- `PORT`: 8080

## 🛠️ Development

### Stop the Services

```bash
cd deployments
docker-compose down
```

### View Service Logs

```bash
# View all logs
docker-compose logs

# View specific service logs
docker-compose logs app
docker-compose logs postgres
docker-compose logs redis
```

### Connect to Database (Optional)

```bash
docker exec -it urlshortener-postgres psql -U postgres -d urlshortener
```

## 📚 Learning Concepts

This project demonstrates:

1. **Go Best Practices**: Clean architecture, proper package organization, error handling
2. **System Design**: URL shortening algorithm, database schema design, caching strategies
3. **Performance Optimization**: Connection pooling, prepared statements, caching layers
4. **Testing**: Unit tests, integration tests, load testing, benchmarking
5. **DevOps**: Docker containerization, Docker Compose orchestration
6. **Monitoring**: Health checks, performance metrics, logging

## 🔄 Key Features

- **High Performance**: Optimized for high concurrency with connection pooling
- **Caching**: Redis caching for frequently accessed URLs
- **Scalable**: Designed for horizontal scaling
- **Comprehensive Testing**: Unit, integration, and performance tests
- **Production Ready**: Health checks, proper error handling, logging
- **Benchmarking**: Built-in performance testing and analysis

## 📊 Performance Metrics

Based on benchmark results:

- **Peak Performance**: 6,892 RPS
- **Optimal Load**: 255 concurrent users
- **Breaking Point**: 405 concurrent users
- **Daily Capacity**: 73.4M URLs
- **Response Time**: 24.7ms (optimal)
- **Reliability**: 100% up to 355 users

## 🔄 Next Steps (Future Improvements)

- [ ] Add authentication and rate limiting
- [ ] Implement URL expiration
- [ ] Add click tracking and analytics
- [ ] Add API documentation (Swagger/OpenAPI)
- [ ] Implement URL validation and sanitization
- [ ] Add monitoring and alerting
- [ ] Implement horizontal scaling with load balancer
- [ ] Add database sharding for higher capacity

## 🔄 Running Tests

To run tests, use the following commands:

```bash
# Run all tests
go test -v ./tests

# Run specific test types
go test -v ./tests/integration
go test -v ./tests/benchmark

# Run with coverage
go test -v -cover ./tests/...

# Run benchmarks
go test -bench=. ./tests/benchmark
```

## 📋 API Documentation

### Endpoints

#### POST /shorten

Creates a short URL from a long URL.

**Request:**

```json
{
  "long_url": "https://www.example.com/very/long/url"
}
```

**Response:**

```json
{
  "short_url": "http://localhost:8080/abc123",
  "long_url": "https://www.example.com/very/long/url",
  "short_code": "abc123",
  "cached": false,
  "cache_type": "none"
}
```

#### GET /{shortCode}

Redirects to the original URL.

**Response:** HTTP 302 redirect to the original URL

#### GET /health

Health check endpoint.

**Response:**

```json
{
  "status": "healthy"
}
```

---

*For detailed documentation, see [docs/README.md](docs/README.md)*  
*For benchmark reports, see [BENCHMARK_REPORT.md](BENCHMARK_REPORT.md) and [BENCHMARK_SUMMARY.md](BENCHMARK_SUMMARY.md)*
