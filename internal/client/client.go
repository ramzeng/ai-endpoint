package client

import "time"

type Client struct {
	Id         uint64
	Name       string
	Secret     string
	RateLimits []RateLimit `gorm:"serializer:json"`
	CreatedAt  time.Time   `gorm:"autoCreateTime"`
	LastUsedAt time.Time
}

type RateLimit struct {
	Model                string
	MaxRequestsPerMinute uint64 `json:"max_requests_per_minute"`
}
