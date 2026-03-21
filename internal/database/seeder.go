package database

import (
	"encoding/json"
	"log"
	"os"

	"api-global/internal/models"

	"gorm.io/gorm"
)

// SeedTerritories lee el JSON e inserta los datos si la tabla está vacía
func SeedTerritories(db *gorm.DB) {
	var count int64
	db.Model(&models.Estado{}).Count(&count)

	// Si ya hay registros, no hacemos nada para no duplicar
	if count > 0 {
		log.Println("⚡ Base de datos ya poblada. Omitiendo Seeder territorial.")
		return
	}

	log.Println("🌱 Iniciando Seeder: Cargando datos territoriales de Venezuela...")

	// Leemos el archivo JSON (La ruta relativa funcionará porque la configuramos en el Dockerfile)
	bytes, err := os.ReadFile("internal/database/seeds/venezuela.json")
	if err != nil {
		log.Printf("❌ Error leyendo archivo venezuela.json: %v\n", err)
		return
	}

	var estados []models.Estado
	if err := json.Unmarshal(bytes, &estados); err != nil {
		log.Printf("❌ Error decodificando el JSON: %v\n", err)
		return
	}

	// Magia de GORM: Inserta el estado, obtiene el ID, inserta los municipios con ese ID,
	// obtiene sus IDs e inserta las parroquias correspondientes. Todo en una transacción.
	if err := db.Create(&estados).Error; err != nil {
		log.Printf("❌ Error insertando datos en PostgreSQL: %v\n", err)
		return
	}

	log.Println("✅ Seeder completado: Venezuela cargada exitosamente en la base de datos.")
}
