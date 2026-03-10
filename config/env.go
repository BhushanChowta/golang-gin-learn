package config

import (
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

type AppConfig struct {
	ServerPort    string
	BaseURL       string
	RedisAddr     string
	RedisPassword string
	RedisDB       int
}

func LoadEnv() *AppConfig {
	if err := godotenv.Load(); err != nil {
		log.Printf(".env not found, using system environment/defaults")
	}

	redisDB := getEnvAsInt("REDIS_DB", 0)

	return &AppConfig{
		ServerPort:    getEnv("APP_PORT", "8080"),
		BaseURL:       getEnv("BASE_URL", "http://localhost:8080/get"),
		RedisAddr:     getEnv("REDIS_ADDR", "127.0.0.1:6379"),
		RedisPassword: getEnv("REDIS_PASSWORD", ""),
		RedisDB:       redisDB,
	}
}

func getEnv(key, fallback string) string {
	val, ok := os.LookupEnv(key)
	if !ok || val == "" {
		return fallback
	}
	return val
}

func getEnvAsInt(key string, fallback int) int {
	val := os.Getenv(key)
	if val == "" {
		return fallback
	}

	parsed, err := strconv.Atoi(val)
	if err != nil {
		return fallback
	}

	return parsed
}
