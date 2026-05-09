package repositories

import (
	"api-global/internal/models"

	"gorm.io/gorm"
)

func GetEstados(db *gorm.DB) ([]models.Estado, error) {
	var estados []models.Estado
	err := db.Order("nombre asc").Find(&estados).Error
	return estados, err
}

func GetMunicipiosByEstado(db *gorm.DB, estadoID int) ([]models.Municipio, error) {
	var municipios []models.Municipio
	err := db.Where("estado_id = ?", estadoID).Order("nombre asc").Find(&municipios).Error
	return municipios, err
}

func GetParroquiasByMunicipio(db *gorm.DB, municipioID int) ([]models.Parroquia, error) {
	var parroquias []models.Parroquia
	err := db.Where("municipio_id = ?", municipioID).Order("nombre asc").Find(&parroquias).Error
	return parroquias, err
}

func GetCiudadesByMunicipio(db *gorm.DB, municipioID int) ([]models.Ciudad, error) {
	var ciudades []models.Ciudad
	err := db.Where("municipio_id = ?", municipioID).Order("id asc").Find(&ciudades).Error
	return ciudades, err
}
