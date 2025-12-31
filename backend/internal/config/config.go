package config

import (
	"fmt"
	"os"
	"strconv"
)

// Config holds all configuration for the application
type Config struct {
	Server   ServerConfig
	Database DatabaseConfig
	Redis    RedisConfig
	JWT      JWTConfig
	FCM      FCMConfig
}

// ServerConfig holds server-related configuration
type ServerConfig struct {
	Port           string
	Environment    string
	AllowedOrigins string
}

// DatabaseConfig holds database-related configuration
type DatabaseConfig struct {
	Host                   string
	Port                   string
	User                   string
	Password               string
	Name                   string
	SSLMode                string
	Timezone               string
	MaxIdleConns           int
	MaxOpenConns           int
	ConnMaxLifetimeMinutes int
	LogLevel               string
}

// RedisConfig holds Redis-related configuration
type RedisConfig struct {
	Host         string
	Port         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
}

// JWTConfig holds JWT-related configuration
type JWTConfig struct {
	SecretKey            string
	AccessTokenDuration  int // in minutes
	RefreshTokenDuration int // in hours
	Issuer               string
}

// FCMConfig holds Firebase Cloud Messaging configuration
type FCMConfig struct {
	CredentialsFile string
	ProjectID       string
}

// Load loads configuration from environment variables
func Load() (*Config, error) {
	cfg := &Config{
		Server: ServerConfig{
			Port:           getEnv("SERVER_PORT", "8080"),
			Environment:    getEnv("ENVIRONMENT", "development"),
			AllowedOrigins: getEnv("ALLOWED_ORIGINS", "*"),
		},
		Database: DatabaseConfig{
			Host:                   getEnv("DB_HOST", "localhost"),
			Port:                   getEnv("DB_PORT", "5432"),
			User:                   getEnv("DB_USER", "postgres"),
			Password:               getEnv("DB_PASSWORD", "postgres"),
			Name:                   getEnv("DB_NAME", "school_management"),
			SSLMode:                getEnv("DB_SSL_MODE", "disable"),
			Timezone:               getEnv("DB_TIMEZONE", "Asia/Jakarta"),
			MaxIdleConns:           getEnvAsInt("DB_MAX_IDLE_CONNS", 10),
			MaxOpenConns:           getEnvAsInt("DB_MAX_OPEN_CONNS", 100),
			ConnMaxLifetimeMinutes: getEnvAsInt("DB_CONN_MAX_LIFETIME_MINUTES", 60),
			LogLevel:               getEnv("DB_LOG_LEVEL", "info"),
		},
		Redis: RedisConfig{
			Host:         getEnv("REDIS_HOST", "localhost"),
			Port:         getEnv("REDIS_PORT", "6379"),
			Password:     getEnv("REDIS_PASSWORD", ""),
			DB:           getEnvAsInt("REDIS_DB", 0),
			PoolSize:     getEnvAsInt("REDIS_POOL_SIZE", 10),
			MinIdleConns: getEnvAsInt("REDIS_MIN_IDLE_CONNS", 5),
		},
		JWT: JWTConfig{
			SecretKey:            getEnv("JWT_SECRET_KEY", "your-secret-key-change-in-production"),
			AccessTokenDuration:  getEnvAsInt("JWT_ACCESS_TOKEN_DURATION", 15),   // 15 minutes
			RefreshTokenDuration: getEnvAsInt("JWT_REFRESH_TOKEN_DURATION", 168), // 7 days (168 hours)
			Issuer:               getEnv("JWT_ISSUER", "school-management-api"),
		},
		FCM: FCMConfig{
			CredentialsFile: getEnv("FCM_CREDENTIALS_FILE", ""),
			ProjectID:       getEnv("FCM_PROJECT_ID", ""),
		},
	}

	// Validate required configuration
	if err := cfg.Validate(); err != nil {
		return nil, err
	}

	return cfg, nil
}

// Validate validates the configuration
func (c *Config) Validate() error {
	// Database validation
	if c.Database.Host == "" {
		return fmt.Errorf("DB_HOST is required")
	}
	if c.Database.User == "" {
		return fmt.Errorf("DB_USER is required")
	}
	if c.Database.Name == "" {
		return fmt.Errorf("DB_NAME is required")
	}

	// JWT validation for production
	if c.Server.Environment == "production" {
		if c.JWT.SecretKey == "your-secret-key-change-in-production" {
			return fmt.Errorf("JWT_SECRET_KEY must be changed in production")
		}
	}

	return nil
}

// IsDevelopment returns true if running in development mode
func (c *Config) IsDevelopment() bool {
	return c.Server.Environment == "development"
}

// IsProduction returns true if running in production mode
func (c *Config) IsProduction() bool {
	return c.Server.Environment == "production"
}

// Helper functions

func getEnv(key, defaultValue string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return defaultValue
}

func getEnvAsInt(key string, defaultValue int) int {
	if value, exists := os.LookupEnv(key); exists {
		if intValue, err := strconv.Atoi(value); err == nil {
			return intValue
		}
	}
	return defaultValue
}

func getEnvAsBool(key string, defaultValue bool) bool {
	if value, exists := os.LookupEnv(key); exists {
		if boolValue, err := strconv.ParseBool(value); err == nil {
			return boolValue
		}
	}
	return defaultValue
}
