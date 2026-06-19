package database

import (
	"fmt"
	"log"
	"time"

	"api-global/internal/config"

	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// InitDB inicializa la conexión a PostgreSQL con sistema de reintentos
func InitDB(cfg *config.Config) (*gorm.DB, error) {
	dsn := cfg.GetDSN()

	var db *gorm.DB
	var err error

	// Sistema de reintentos: Intentará conectar 5 veces, esperando 2 segundos entre intentos.
	for i := 1; i <= 5; i++ {
		db, err = gorm.Open(postgres.Open(dsn), &gorm.Config{})
		if err == nil {
			log.Println("✅ Conexión a PostgreSQL establecida exitosamente")
			return db, nil // Conexión exitosa, salimos del bucle
		}

		log.Printf("⏳ Esperando a la base de datos (Intento %d/5)...\n", i)
		time.Sleep(2 * time.Second)
	}

	// Si después de 5 intentos falló, devolvemos el error fatal
	return nil, fmt.Errorf("error conectando a la base de datos tras varios intentos: %v", err)
}
