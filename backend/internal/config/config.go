package config

import (
	"os"
	"strconv"
	"time"
)

const (
	DefaultRateLimitRequests = 1000
	DefaultRateLimitWindow   = 3600
	defaultMaxOpenConns      = 25
	defaultMaxIdleConns      = 5
	defaultConnMaxLifetime   = 300
)

type DatabaseConfig struct {
	Host            string
	Port            string
	User            string
	Password        string
	Name            string
	SSLMode         string
	MaxOpenConns    int
	MaxIdleConns    int
	ConnMaxLifetime time.Duration
}

type Config struct {
	// Slice types first (24 bytes on 64-bit)
	AllowedOrigins []string

	// Nested struct
	Database DatabaseConfig

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
		Database: DatabaseConfig{
			Host:            getEnv("DB_HOST", "localhost"),
			Port:            getEnv("DB_PORT", "5432"),
			User:            getEnv("DB_USER", "crewee"),
			Password:        getEnv("DB_PASSWORD", "password"),
			Name:            getEnv("DB_NAME", "crewee_dev"),
			SSLMode:         getEnv("DB_SSL_MODE", "disable"),
			MaxOpenConns:    getEnvInt("DB_MAX_OPEN_CONNS", defaultMaxOpenConns),
			MaxIdleConns:    getEnvInt("DB_MAX_IDLE_CONNS", defaultMaxIdleConns),
			ConnMaxLifetime: time.Duration(getEnvInt("DB_CONN_MAX_LIFETIME", defaultConnMaxLifetime)) * time.Second,
		},
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
