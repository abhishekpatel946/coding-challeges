package models

// URLMapping represents a URL entry in our database
type URLMapping struct {
	ID        int64  `json:"id"`
	ShortCode string `json:"short_code"`
	LongURL   string `json:"long_url"`
}

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
