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

func GetTribunalesByEstado(db *gorm.DB, estadoID int) ([]models.Tribunal, error) {
	var tribunales []models.Tribunal
	err := db.Where("estado_id = ?", estadoID).Order("categoria asc, nombre asc").Find(&tribunales).Error
	return tribunales, err
}

func GetTribunalesByMunicipio(db *gorm.DB, municipioID int) ([]models.Tribunal, error) {
	var tribunales []models.Tribunal
	err := db.Joins("JOIN tribunal_municipios ON tribunal_municipios.tribunal_id = tribunals.id").
		Where("tribunal_municipios.municipio_id = ?", municipioID).
		Order("tribunals.nombre asc").Find(&tribunales).Error
	return tribunales, err
}

func GetCiudadesByEstado(db *gorm.DB, estadoID int) ([]models.Ciudad, error) {
	var ciudades []models.Ciudad
	err := db.Where("estado_id = ?", estadoID).Order("nombre asc").Find(&ciudades).Error
	return ciudades, err
}

func GetCodigosArea(db *gorm.DB) ([]models.CodigoArea, error) {
	var codigos []models.CodigoArea
	err := db.Order("codigo asc").Find(&codigos).Error
	return codigos, err
}
