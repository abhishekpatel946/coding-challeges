-- Database initialization script for URL shortener
-- This script runs when the PostgreSQL container starts for the first time

-- Create the urls table with optimized structure
CREATE TABLE IF NOT EXISTS urls (
    id BIGSERIAL PRIMARY KEY,
    short_code VARCHAR(10) UNIQUE NOT NULL,
    long_url TEXT NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create optimized indexes for fast lookups
CREATE INDEX IF NOT EXISTS idx_short_code ON urls(short_code);
CREATE INDEX IF NOT EXISTS idx_created_at ON urls(created_at);
CREATE INDEX IF NOT EXISTS idx_long_url_hash ON urls USING hash(long_url);

-- Add some helpful comments
COMMENT ON TABLE urls IS 'Stores mappings between short codes and long URLs';
COMMENT ON COLUMN urls.id IS 'Auto-incrementing primary key';
COMMENT ON COLUMN urls.short_code IS 'Unique short code (e.g., "abc123")';
COMMENT ON COLUMN urls.long_url IS 'The original long URL';
COMMENT ON COLUMN urls.created_at IS 'Timestamp when the short URL was created';

-- Optimize PostgreSQL settings for high concurrency
ALTER SYSTEM SET max_connections = 200;
ALTER SYSTEM SET shared_buffers = '256MB';
ALTER SYSTEM SET effective_cache_size = '1GB';
ALTER SYSTEM SET maintenance_work_mem = '64MB';
ALTER SYSTEM SET checkpoint_completion_target = 0.9;
ALTER SYSTEM SET wal_buffers = '16MB';
ALTER SYSTEM SET default_statistics_target = 100;
ALTER SYSTEM SET random_page_cost = 1.1;
ALTER SYSTEM SET effective_io_concurrency = 200;
ALTER SYSTEM SET work_mem = '4MB';
ALTER SYSTEM SET min_wal_size = '1GB';
ALTER SYSTEM SET max_wal_size = '4GB';

-- Reload configuration
SELECT pg_reload_conf(); 