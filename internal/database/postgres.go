package database

import (
	"fmt"
	"log"
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
		// Le pasamos directamente la URL completa a GORM
		db, err = gorm.Open(postgres.Open(cfg.DBUrl), &gorm.Config{
			PrepareStmt: false, // Desactivado para compatibilidad con PgBouncer (Supabase)
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
