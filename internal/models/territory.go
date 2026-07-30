package models

// Estado representa la entidad principal de división territorial
type Estado struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre string `gorm:"type:varchar(100);not null;unique" json:"nombre"`
	// Relación: Un Estado "tiene muchos" (Has Many) Municipios
	Municipios []Municipio `gorm:"foreignKey:EstadoID" json:"municipios,omitempty"`
	// Relación: Un Estado "tiene muchas" Ciudades
	Ciudades []Ciudad `gorm:"foreignKey:EstadoID" json:"ciudades,omitempty"`
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

// Tribunal representa un tribunal judicial del sistema de justicia venezolano.
// Los tribunales de categoría "Superior" y "Primera Instancia" son estadales (usan EstadoID).
// Los de categoría "Municipio" son municipales y pueden pertenecer a varios municipios (many-to-many).
type Tribunal struct {
	ID         uint        `gorm:"primaryKey;autoIncrement" json:"id"`
	Nombre     string      `gorm:"type:varchar(500);not null" json:"nombre"`
	Categoria  string      `gorm:"type:varchar(50);not null" json:"categoria"` // "Superior", "Primera Instancia", "Municipio"
	EstadoID   *uint       `gorm:"index" json:"estado_id,omitempty"`
	Municipios []Municipio `gorm:"many2many:tribunal_municipios;" json:"municipios,omitempty"`
}

// Ciudad pertenece a un Estado
type Ciudad struct {
	ID       uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	EstadoID uint   `gorm:"not null" json:"estado_id"`
	Nombre   string `gorm:"type:varchar(100);not null" json:"nombre"`
}

// CodigoArea representa un código de área telefónico de Venezuela
type CodigoArea struct {
	ID     uint   `gorm:"primaryKey;autoIncrement" json:"id"`
	Codigo string `gorm:"type:varchar(10);not null;unique" json:"codigo"`
}
