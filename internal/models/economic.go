package models

import "time"

// IndicadorEconomico guarda el valor del BCV (Dolar/Euro) y la UCAUU.
type IndicadorEconomico struct {
	ID        uint      `gorm:"primaryKey;autoIncrement" json:"id"`
	Tipo      string    `gorm:"type:varchar(10);not null;uniqueIndex:idx_tipo_fecha" json:"tipo"`
	Valor     float64   `gorm:"type:numeric(15,4);not null" json:"valor"`
	Fecha     time.Time `gorm:"type:date;not null;uniqueIndex:idx_tipo_fecha" json:"fecha"`
	CreatedAt time.Time `json:"created_at"`
}

func (IndicadorEconomico) TableName() string {
	return "indicador_economicos"
}
