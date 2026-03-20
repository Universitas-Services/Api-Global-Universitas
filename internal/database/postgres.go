package database

import (
	"api-global/internal/config"
	"fmt"
	"log"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inicializa la conexión a PostgreSQL usando GORM
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	// Construimos la cadena de conexión (DSN)
	dsn := fmt.Sprintf("host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		cfg.DBHost, cfg.DBUser, cfg.DBPassword, cfg.DBName, cfg.DBPort, cfg.DBSSLMode)

	// Abrimos la conexión
	db, err := gorm.Open(postgres.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, fmt.Errorf("error conectando a la base de datos: %v", err)
	}

	log.Println("✅ Conexión a PostgreSQL establecida exitosamente")
	return db, nil
}
