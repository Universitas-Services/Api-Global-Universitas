package models

import "time"

// IndicadorEconomico guardará el valor del BCV (Dólar/Euro) y la UCAUU
type IndicadorEconomico struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tipo      string    `gorm:"type:varchar(10);not null;index" json:"tipo"` // Ej: "USD", "EUR", "UCAUU"
	Valor     float64   `gorm:"type:numeric(15,4);not null" json:"valor"`    // Permite 4 decimales para mayor precisión
	Fecha     time.Time `gorm:"type:date;not null;uniqueIndex:idx_tipo_fecha" json:"fecha"`
	CreatedAt time.Time `json:"created_at"`
}
