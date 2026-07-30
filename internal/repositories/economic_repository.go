package repositories

import (
	"api-global/internal/models"
	"api-global/internal/timeutil"
	"time"

	"gorm.io/gorm"
)

// SaveUCAUU guarda o actualiza el valor de la UCAUU para el día actual (Caracas)
func SaveUCAUU(db *gorm.DB, valor float64) error {
	hoy := timeutil.HoyCaracas()

	indicador := models.IndicadorEconomico{
		Tipo:  "UCAUU",
		Valor: valor,
		Fecha: hoy,
	}

	result := db.Where(models.IndicadorEconomico{Tipo: "UCAUU", Fecha: hoy}).
		Assign(models.IndicadorEconomico{Valor: valor}).
		FirstOrCreate(&indicador)

	return result.Error
}

// GetLatestUCAUU obtiene el último valor registrado de la UCAUU
func GetLatestUCAUU(db *gorm.DB) (*models.IndicadorEconomico, error) {
	var indicador models.IndicadorEconomico
	result := db.Where("tipo = ?", "UCAUU").Order("fecha desc").First(&indicador)

	if result.Error != nil {
		return nil, result.Error
	}
	return &indicador, nil
}

// GetIndicadorDeHoy busca si ya tenemos registrado un indicador para el día civil en Caracas
func GetIndicadorDeHoy(db *gorm.DB, tipo string) (*models.IndicadorEconomico, error) {
	var indicador models.IndicadorEconomico
	hoy := timeutil.HoyCaracas()

	err := db.Where("tipo = ? AND fecha = ?", tipo, hoy).First(&indicador).Error
	if err != nil {
		return nil, err
	}
	return &indicador, nil
}

// SaveIndicador guarda o actualiza un valor económico con la fecha de hoy (Caracas).
// Si ya existe tipo+fecha, actualiza el valor (upsert).
func SaveIndicador(db *gorm.DB, tipo string, valor float64) error {
	hoy := timeutil.HoyCaracas()

	indicador := models.IndicadorEconomico{
		Tipo:  tipo,
		Valor: valor,
		Fecha: hoy,
	}

	return db.Where(models.IndicadorEconomico{Tipo: tipo, Fecha: hoy}).
		Assign(models.IndicadorEconomico{Valor: valor}).
		FirstOrCreate(&indicador).Error
}

// GetBCVByFecha obtiene USD y EUR del BCV para una fecha específica (solo lectura de BD)
func GetBCVByFecha(db *gorm.DB, fecha time.Time) (usd float64, eur float64, err error) {
	dia := fecha.Truncate(24 * time.Hour)

	var usdInd, eurInd models.IndicadorEconomico
	if err = db.Where("tipo = ? AND fecha = ?", "USD_BCV", dia).First(&usdInd).Error; err != nil {
		return 0, 0, err
	}
	if err = db.Where("tipo = ? AND fecha = ?", "EUR_BCV", dia).First(&eurInd).Error; err != nil {
		return 0, 0, err
	}
	return usdInd.Valor, eurInd.Valor, nil
}
