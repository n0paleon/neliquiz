package config

import (
	"fmt"
	"github.com/joho/godotenv"
	"log"
	"os"
	"path/filepath"
	"strconv"
	"sync"
)

var (
	once   sync.Once
	config *Config
)

// Config holds all app config values.
type Config struct {
	DBUser            string
	DBPassword        string
	DBHost            string
	DBPort            string
	DBName            string
	DBSSLMode         string
	DBMaxIdleConn     int
	DBMaxOpenConn     int
	DBMaxConnLifetime int
	HTTPHost          string
	HTTPPort          string
	HTTPPrefork       bool
}

// New initializes the config.
func New() *Config {
	once.Do(func() {
		wd, _ := os.Getwd()
		fmt.Println("WD:", wd)

		envPath := os.Getenv("ENV_PATH")
		if envPath == "" {
			rootPath, err := findRootPathWithEnv(wd)
			if err != nil {
				log.Printf("Failed to find .env: %v", err)
			}
			envPath = filepath.Join(rootPath, ".env")
		}

		if err := godotenv.Load(envPath); err != nil {
			log.Printf("Warning: .env file not found: %v", err)
		}

		config = &Config{
			DBUser:            GetEnv("DB_USER", "postgres"),
			DBPassword:        GetEnv("DB_PASSWORD", "postgres"),
			DBHost:            GetEnv("DB_HOST", "localhost"),
			DBPort:            GetEnv("DB_PORT", "5432"),
			DBName:            GetEnv("DB_NAME", "postgres"),
			DBSSLMode:         GetEnv("DB_SSLMODE", "disable"),
			DBMaxIdleConn:     GetEnvAsInt("DB_MAX_IDLE_CONN", 2),
			DBMaxOpenConn:     GetEnvAsInt("DB_MAX_OPEN_CONN", 5),
			DBMaxConnLifetime: GetEnvAsInt("DB_MAX_CONN_LIFETIME", 60),
			HTTPHost:          GetEnv("HTTP_HOST", "0.0.0.0"),
			HTTPPort:          GetEnv("HTTP_PORT", "3000"),
			HTTPPrefork:       GetEnvAsBool("HTTP_PREFORK", false),
		}
	})

	return config
}

func findRootPathWithEnv(startDir string) (string, error) {
	dir := startDir
	for {
		envPath := filepath.Join(dir, ".env")
		if _, err := os.Stat(envPath); err == nil {
			return dir, nil
		}

		parent := filepath.Dir(dir)
		if parent == dir {
			// Reached filesystem root
			return "", fmt.Errorf(".env not found from %s up", startDir)
		}
		dir = parent
	}
}

// GetConfig returns the singleton config.
func GetConfig() *Config {
	return New()
}

// GetEnv reads an env var as string.
func GetEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}

// GetEnvAsInt reads an env var as int.
func GetEnvAsInt(key string, fallback int) int {
	valStr := GetEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.Atoi(valStr)
	if err != nil {
		return fallback
	}
	return val
}

// GetEnvAsBool reads an env var as bool.
func GetEnvAsBool(key string, fallback bool) bool {
	valStr := GetEnv(key, "")
	if valStr == "" {
		return fallback
	}
	val, err := strconv.ParseBool(valStr)
	if err != nil {
		return fallback
	}
	return val
}
