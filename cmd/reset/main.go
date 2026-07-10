package main

import (
	"log"

	"api-global/internal/config"
	"api-global/internal/database"
	"api-global/internal/models"
)

func main() {
	cfg := config.LoadConfig()
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("No se pudo iniciar la base de datos: %v", err)
	}

	log.Println("Eliminando tablas territoriales para reiniciar contadores...")

	err = db.Migrator().DropTable(
		&models.Tribunal{},
		"tribunal_municipios",
		&models.Ciudad{},
		&models.Parroquia{},
		&models.Municipio{},
		&models.Estado{},
	)
	if err != nil {
		log.Fatalf("Error eliminando tablas: %v", err)
	}

	log.Println("✅ Base de datos limpiada y contadores reiniciados exitosamente.")
}
