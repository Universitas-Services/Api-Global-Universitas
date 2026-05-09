package models

// Estado representa la entidad principal de division territorial.
type Estado struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre     string      `gorm:"type:varchar(100);not null;unique" json:"nombre"`
	Municipios []Municipio `gorm:"foreignKey:EstadoID" json:"municipios,omitempty"`
}

func (Estado) TableName() string {
	return "estados"
}

// Municipio pertenece a un Estado y tiene muchas parroquias y ciudades.
type Municipio struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	EstadoID   uint        `gorm:"not null;uniqueIndex:idx_municipio_estado_nombre" json:"estado_id"`
	Nombre     string      `gorm:"type:varchar(100);not null;uniqueIndex:idx_municipio_estado_nombre" json:"nombre"`
	Parroquias []Parroquia `gorm:"foreignKey:MunicipioID" json:"parroquias,omitempty"`
	Ciudades   []Ciudad    `gorm:"foreignKey:MunicipioID" json:"ciudades,omitempty"`
}

func (Municipio) TableName() string {
	return "municipios"
}

// Parroquia es el nivel mas granular y pertenece a un Municipio.
type Parroquia struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	MunicipioID uint   `gorm:"not null;uniqueIndex:idx_parroquia_municipio_nombre" json:"municipio_id"`
	Nombre      string `gorm:"type:varchar(100);not null;uniqueIndex:idx_parroquia_municipio_nombre" json:"nombre"`
}

func (Parroquia) TableName() string {
	return "parroquia"
}

// Ciudad representa un centro poblado principal asociado a un Municipio.
type Ciudad struct {
	ID          uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	MunicipioID uint   `gorm:"not null;uniqueIndex:idx_ciudad_municipio_nombre" json:"municipio_id"`
	Nombre      string `gorm:"type:varchar(100);not null;uniqueIndex:idx_ciudad_municipio_nombre" json:"nombre"`
}

func (Ciudad) TableName() string {
	return "ciudades"
}
