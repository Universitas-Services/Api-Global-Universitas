package repositories

import (
	"api-global/internal/models"
	"time"

	"gorm.io/gorm"
)

// SaveUCAUU guarda o actualiza el valor de la UCAUU para el día actual
func SaveUCAUU(db *gorm.DB, valor float64) error {
	hoy := time.Now().Truncate(24 * time.Hour) // Solo la fecha, sin horas

	indicador := models.IndicadorEconomico{
		Tipo:  "UCAUU",
		Valor: valor,
		Fecha: hoy,
	}

	// FirstOrCreate: Si ya existe un registro hoy para UCAUU, lo actualiza, si no, lo crea.
	// Esto cumple tu requerimiento de "actualizar ese dato" si se envía varias veces el mismo día.
	result := db.Where(models.IndicadorEconomico{Tipo: "UCAUU", Fecha: hoy}).
		Assign(models.IndicadorEconomico{Valor: valor}).
		FirstOrCreate(&indicador)

	return result.Error
}

// GetLatestUCAUU obtiene el último valor registrado de la UCAUU
func GetLatestUCAUU(db *gorm.DB) (*models.IndicadorEconomico, error) {
	var indicador models.IndicadorEconomico
	// Ordenamos por fecha descendente y tomamos el primero
	result := db.Where("tipo = ?", "UCAUU").Order("fecha desc").First(&indicador)

	if result.Error != nil {
		return nil, result.Error
	}
	return &indicador, nil
}

// GetIndicadorDeHoy busca si ya tenemos registrado un indicador específico para la fecha actual
func GetIndicadorDeHoy(db *gorm.DB, tipo string) (*models.IndicadorEconomico, error) {
	var indicador models.IndicadorEconomico
	hoy := time.Now().Truncate(24 * time.Hour)

	err := db.Where("tipo = ? AND fecha = ?", tipo, hoy).First(&indicador).Error
	if err != nil {
		return nil, err // Retornará error si no lo encuentra
	}
	return &indicador, nil
}

// SaveIndicador guarda un nuevo valor económico con la fecha de hoy
func SaveIndicador(db *gorm.DB, tipo string, valor float64) error {
	hoy := time.Now().Truncate(24 * time.Hour)

	indicador := models.IndicadorEconomico{
		Tipo:  tipo,
		Valor: valor,
		Fecha: hoy,
	}

	// Usamos FirstOrCreate por seguridad, por si hay peticiones concurrentes
	return db.Where(models.IndicadorEconomico{Tipo: tipo, Fecha: hoy}).
		Assign(models.IndicadorEconomico{Valor: valor}).
		FirstOrCreate(&indicador).Error
}
