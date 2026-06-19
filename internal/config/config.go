package config

import (
	"os"

	"github.com/joho/godotenv"
)

type Config struct {
	Port       string
	DBHost     string
	DBUser     string
	DBPassword string
	DBName     string
	DBPort     string
	DBSSLMode  string
}

// LoadConfig lee las variables de entorno o carga el archivo .env si existe
func LoadConfig() *Config {
	_ = godotenv.Load()

	return &Config{
		Port:       getEnv("PORT", "8080"),
		DBHost:     getEnv("DB_HOST", "localhost"),
		DBUser:     getEnv("DB_USER", "postgres"),
		DBPassword: getEnv("DB_PASSWORD", "secret"),
		DBName:     getEnv("DB_NAME", "universitas_db"),
		DBPort:     getEnv("DB_PORT", "5432"),
		DBSSLMode:  getEnv("DB_SSLMODE", "disable"),
	}
}

// GetDSN retorna el string de conexión a la base de datos
func (c *Config) GetDSN() string {
	// Si existe DATABASE_URL (usado típicamente en Render/Supabase), la preferimos
	if dbURL := os.Getenv("DATABASE_URL"); dbURL != "" {
		return dbURL
	}
	
	// Si no, armamos el DSN con las variables individuales
	return "host=" + c.DBHost + " user=" + c.DBUser + " password=" + c.DBPassword + " dbname=" + c.DBName + " port=" + c.DBPort + " sslmode=" + c.DBSSLMode
}

// Helper para obtener variable o un valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists {
		return value
	}
	return fallback
}
