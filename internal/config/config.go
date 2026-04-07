package config

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

type Config struct {
	Port        string
	DBUrl       string // URL de conexión a la base de datos (Supabase)
	CorsOrigins []string
	RunSeeder    bool   // Controlar si se ejecuta el seeder territorial
	RunMigrations bool  // Controlar si se ejecuta el AutoMigrate
}

func LoadConfig() *Config {
	_ = godotenv.Load()

	// Leemos la variable de entorno
	// Si no existe, dejamos localhost:3000 y 3001 como valor por defecto por seguridad
	corsString := getEnv("CORS_ALLOWED_ORIGINS", "http://localhost:3000,http://localhost:3001")

	// Dividimos el texto por las comas para crear el arreglo
	var origins []string
	if corsString != "" {
		originsRaw := strings.Split(corsString, ",")
		// Limpiamos posibles espacios en blanco extra por si alguien escribe "url1, url2"
		for _, origin := range originsRaw {
			origins = append(origins, strings.TrimSpace(origin))
		}
	}

	return &Config{
		Port:        getEnv("PORT", "8080"),
		DBUrl:       getEnv("DATABASE_URL", ""),
		CorsOrigins: origins,
		RunSeeder:     getEnv("RUN_SEEDER", "false") == "true",
		RunMigrations: getEnv("RUN_MIGRATIONS", "false") == "true",
	}
}

// Helper para obtener variable o un valor por defecto
func getEnv(key, fallback string) string {
	if value, exists := os.LookupEnv(key); exists && value != "" {
		return value
	}
	return fallback
}
