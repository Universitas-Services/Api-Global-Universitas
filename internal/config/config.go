package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port  string
	DBUrl string // URL de conexión a la base de datos (Supabase)
}

func LoadConfig() *Config {
	// Intentamos cargar .env en local
	_ = godotenv.Load()

	return &Config{
		Port:  getEnv("PORT", "8080"),
		DBUrl: getEnv("DATABASE_URL", ""),
	}
}

// Helper para obtener variable o un valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
