package service

import (
	"context"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"strings"
	redis "github.com/redis/go-redis/v9"
)

type URLService struct {
	redisClient *redis.Client
	baseURL     string
}

func NewURLService(redisClient *redis.Client, baseURL string) *URLService {
	return &URLService{
		redisClient: redisClient,
		baseURL:     baseURL,
	}
}

// generateShortCode creates a short code from a URL using MD5 hash
func generateShortCode(longURL string) string {
	hash := md5.Sum([]byte(longURL))
	hashStr := hex.EncodeToString(hash[:])
	// Take first 8 characters for a reasonably short code
	return hashStr[:8]
}

func (s *URLService) ShortenURL(longURL string) (string, error) {
	// Generate short code
	shortCode := generateShortCode(longURL)

	// Store the mapping in Redis (no expiration by setting duration to 0)
	err := s.redisClient.Set(context.Background(), shortCode, longURL, 0).Err()
	if err != nil {
		return "", fmt.Errorf("failed to store URL in Redis: %w", err)
	}

	// Return a fully qualified short URL built from config.
	return fmt.Sprintf("%s/%s", strings.TrimRight(s.baseURL, "/"), shortCode), nil
}

func (s *URLService) GetLongURL(ctx context.Context, shortURL string) (string, error) {
	// Fetch long URL by short code.
	val, err := s.redisClient.Get(ctx, shortURL).Result()
	if err != nil {
		return "", fmt.Errorf("URL not found: %w", err)
	}
	return val, nil
}