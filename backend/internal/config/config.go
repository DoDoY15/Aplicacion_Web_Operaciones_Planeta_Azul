package config

import (
	"os"
	"time"
)

type Config struct {
	Port             string
	Env              string
	JWTPrivateKeyPath string
	JWTPublicKeyPath  string
	AccessExpiry     time.Duration
	RefreshExpiry    time.Duration
	DB               DBConfig
	FabricaDB        DBConfig
	OCRServiceURL    string
}

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
}

func Load() *Config {
	accessExpiry, _ := time.ParseDuration(getEnv("JWT_ACCESS_EXPIRY", "8h"))
	refreshExpiry, _ := time.ParseDuration(getEnv("JWT_REFRESH_EXPIRY", "168h"))

	return &Config{
		Port:             getEnv("PORT", "8080"),
		Env:              getEnv("ENV", "development"),
		JWTPrivateKeyPath: getEnv("JWT_PRIVATE_KEY_PATH", "./scripts/private.pem"),
		JWTPublicKeyPath:  getEnv("JWT_PUBLIC_KEY_PATH", "./scripts/public.pem"),
		AccessExpiry:     accessExpiry,
		RefreshExpiry:    refreshExpiry,
		DB: DBConfig{
			Host:     getEnv("DB_HOST", "localhost"),
			Port:     getEnv("DB_PORT", "5432"),
			User:     getEnv("DB_USER", "planeta_azul"),
			Password: getEnv("DB_PASSWORD", ""),
			Name:     getEnv("DB_NAME", "planeta_azul_app"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		FabricaDB: DBConfig{
			Host:     getEnv("FABRICA_DB_HOST", "localhost"),
			Port:     getEnv("FABRICA_DB_PORT", "5432"),
			User:     getEnv("FABRICA_DB_USER", ""),
			Password: getEnv("FABRICA_DB_PASSWORD", ""),
			Name:     getEnv("FABRICA_DB_NAME", "fabrica"),
			SSLMode:  getEnv("DB_SSLMODE", "disable"),
		},
		OCRServiceURL: getEnv("OCR_SERVICE_URL", "http://localhost:9000"),
	}
}

func getEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
