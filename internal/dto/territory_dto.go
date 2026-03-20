package dto

type EstadoResponse struct {
	ID     uint   `json:"id" example:"1"`
	Nombre string `json:"nombre" example:"Lara"`
}

type MunicipioResponse struct {
	ID       uint   `json:"id" example:"10"`
	EstadoID uint   `json:"estado_id" example:"1"`
	Nombre   string `json:"nombre" example:"Iribarren"`
}

type ParroquiaResponse struct {
	ID          uint   `json:"id" example:"50"`
	MunicipioID uint   `json:"municipio_id" example:"10"`
	Nombre      string `json:"nombre" example:"Catedral"`
}
