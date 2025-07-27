package config

import (
	"os"
	"strconv"
)

const (
	DefaultRateLimitRequests = 1000
	DefaultRateLimitWindow   = 3600
)

type Config struct {
	// Slice types first (24 bytes on 64-bit)
	AllowedOrigins []string

	// String types (16 bytes each on 64-bit)
	Port               string
	Env                string
	DatabaseURL        string
	TestDatabaseURL    string
	RedisURL           string
	AWSRegion          string
	AWSAccessKeyID     string
	AWSSecretAccessKey string
	CognitoUserPoolID  string
	CognitoClientID    string
	CognitoRegion      string
	JWTSecret          string
	LogLevel           string
	LogFormat          string
	MaxFileSize        string
	UploadPath         string

	// Int types last (8 bytes each on 64-bit)
	RateLimitRequests int
	RateLimitWindow   int
}

func Load() (*Config, error) {
	cfg := &Config{
		Port:               getEnv("PORT", "8080"),
		Env:                getEnv("ENV", "development"),
		DatabaseURL:        getEnv("DATABASE_URL", ""),
		TestDatabaseURL:    getEnv("TEST_DATABASE_URL", ""),
		RedisURL:           getEnv("REDIS_URL", ""),
		AWSRegion:          getEnv("AWS_REGION", "ap-northeast-1"),
		AWSAccessKeyID:     getEnv("AWS_ACCESS_KEY_ID", ""),
		AWSSecretAccessKey: getEnv("AWS_SECRET_ACCESS_KEY", ""),
		CognitoUserPoolID:  getEnv("COGNITO_USER_POOL_ID", ""),
		CognitoClientID:    getEnv("COGNITO_CLIENT_ID", ""),
		CognitoRegion:      getEnv("COGNITO_REGION", "ap-northeast-1"),
		JWTSecret:          getEnv("JWT_SECRET", ""),
		LogLevel:           getEnv("LOG_LEVEL", "info"),
		LogFormat:          getEnv("LOG_FORMAT", "json"),
		MaxFileSize:        getEnv("MAX_FILE_SIZE", "10MB"),
		UploadPath:         getEnv("UPLOAD_PATH", "./tmp/uploads"),
		RateLimitRequests:  getEnvInt("RATE_LIMIT_REQUESTS", DefaultRateLimitRequests),
		RateLimitWindow:    getEnvInt("RATE_LIMIT_WINDOW", DefaultRateLimitWindow),
	}

	// Parse CORS allowed origins
	if origins := getEnv("ALLOWED_ORIGINS", ""); origins != "" {
		// TODO: Parse comma-separated origins
		cfg.AllowedOrigins = []string{origins}
	}

	return cfg, nil
}

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func getEnvInt(key string, defaultValue int) int {
	if value := os.Getenv(key); value != "" {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}
