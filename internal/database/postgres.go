package database

import (
	"fmt"
	"log"
	"strings"
	"time"

	"api-global/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inicializa la conexión a Supabase usando GORM
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("❌ Error fatal: La variable DATABASE_URL está vacía")
	}

	var db *gorm.DB
	var err error

	// Sistema de reintentos para manejar el arranque de Supabase
	for i := 1; i <= 5; i++ {
		// Forzamos prepare_threshold=0 en el DSN para evitar errores con PgBouncer/Supabase
		dsn := cfg.DBUrl
		if !strings.Contains(dsn, "prepare_threshold=0") {
			if strings.Contains(dsn, "?") {
				dsn += "&prepare_threshold=0"
			} else {
				dsn += "?prepare_threshold=0"
			}
		}

		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{
			PrepareStmt:            false, // Desactivado para GORM
			SkipDefaultTransaction: true,  // Recomendado para PgBouncer
		})
		if err == nil {
			log.Println("✅ Conexión a Supabase establecida exitosamente")
			return db, nil
		}

		log.Printf("⏳ Esperando a Supabase (Intento %d/5)...\n", i)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("error conectando a Supabase tras varios intentos: %v", err)
}
