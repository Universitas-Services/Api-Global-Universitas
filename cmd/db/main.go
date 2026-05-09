package main

import (
	"log"
	"os"

	"api-global/internal/config"
	"api-global/internal/database"
)

func main() {
	if len(os.Args) < 2 {
		log.Fatalf("Uso: go run ./cmd/db <migrate|seed-territories|bootstrap-territories>")
	}

	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("No se pudo iniciar la base de datos: %v", err)
	}

	switch os.Args[1] {
	case "migrate":
		if err := database.MigrateAll(db); err != nil {
			log.Fatalf("Error ejecutando migraciones: %v", err)
		}
		log.Println("Migraciones completadas correctamente.")
	case "seed-territories":
		database.SeedTerritories(db)
	case "bootstrap-territories":
		if err := database.MigrateAll(db); err != nil {
			log.Fatalf("Error ejecutando migraciones: %v", err)
		}
		database.SeedTerritories(db)
	default:
		log.Fatalf("Comando no soportado: %s", os.Args[1])
	}
}
