package database

import (
	"api-global/internal/models"

	"gorm.io/gorm"
)

// MigrateAll sincroniza el esquema esperado por la aplicacion.
func MigrateAll(db *gorm.DB) error {
	return db.AutoMigrate(
		&models.Estado{},
		&models.Municipio{},
		&models.Parroquia{},
		&models.Ciudad{},
		&models.IndicadorEconomico{},
	)
}
