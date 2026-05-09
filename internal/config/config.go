package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBUrl       string
	CorsOrigins []string
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	corsString := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001")

	var origins []string
	if corsString != "" {
		originsRaw := strings.Split(corsString, ",")
		for _, origin := range originsRaw {
			origins = append(origins, strings.TrimSpace(origin))
		}
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBUrl:       getEnv("DATABASE_URL", ""),
		CorsOrigins: origins,
	}
}

func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
