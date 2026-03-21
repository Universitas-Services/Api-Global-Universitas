package models

// Estado representa la entidad principal de división territorial
type Estado struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"type:varchar(100);not null;unique" json:"nombre"`
	// Relación: Un Estado "tiene muchos" (Has Many) Municipios
	Municipios []Municipio `gorm:"foreignKey:EstadoID" json:"municipios,omitempty"`
}

// Municipio pertenece a un Estado y tiene muchas Parroquias
type Municipio struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	EstadoID uint   `gorm:"not null" json:"estado_id"`
	Nombre   string `gorm:"type:varchar(100);not null" json:"nombre"`
	// Relación: Un Municipio "tiene muchas" (Has Many) Parroquias
	Parroquias []Parroquia `gorm:"foreignKey:MunicipioID" json:"parroquias,omitempty"`
}

// Parroquia es el nivel más granular y pertenece a un Municipio
type Parroquia struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	MunicipioID uint   `gorm:"not null" json:"municipio_id"`
	Nombre      string `gorm:"type:varchar(100);not null" json:"nombre"`
}
