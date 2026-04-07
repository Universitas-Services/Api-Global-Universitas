package database

import (
	"fmt"
	"log"
	"time"

	"api-global/internal/config"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/stdlib"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inicializa la conexión a Supabase usando GORM + pgx con Simple Protocol
// Esto es OBLIGATORIO para compatibilidad con PgBouncer (Supabase/Render)
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	if cfg.DBUrl == "" {
		return nil, fmt.Errorf("❌ Error fatal: La variable DATABASE_URL está vacía")
	}

	var db *gorm.DB
	var err error

	// Sistema de reintentos para manejar el arranque de Supabase
	for i := 1; i <= 5; i++ {
		// Parseamos la URL de conexión con pgx
		pgxConfig, parseErr := pgx.ParseConfig(cfg.DBUrl)
		if parseErr != nil {
			return nil, fmt.Errorf("❌ Error parseando DATABASE_URL: %v", parseErr)
		}

		// SOLUCIÓN DEFINITIVA: Forzar Simple Protocol en pgx
		// Esto evita que pgx genere "prepared statements" en el servidor,
		// lo cual es incompatible con PgBouncer en Transaction Mode (Supabase)
		pgxConfig.DefaultQueryExecMode = pgx.QueryExecModeSimpleProtocol

		// Convertimos la config de pgx a un *sql.DB estándar
		sqlDB := stdlib.OpenDB(*pgxConfig)

		// Le pasamos el *sql.DB ya configurado a GORM
		db, err = gorm.Open(postgres.New(postgres.Config{
			Conn: sqlDB,
		}), &gorm.Config{
			PrepareStmt:            false, // Doble seguridad: desactivado en GORM
			SkipDefaultTransaction: true,  // Optimización para PgBouncer
		})

		if err == nil {
			log.Println("✅ Conexión a Supabase establecida exitosamente")
			return db, nil
		}

		log.Printf("⏳ Esperando a Supabase (Intento %d/5)...\n", i)
		time.Sleep(2 * time.Second)
	}

	return nil, fmt.Errorf("❌ Error conectando a Supabase tras varios intentos: %v", err)
}
