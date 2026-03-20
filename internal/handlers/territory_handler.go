package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"

	"api-global/internal/dto"
	"api-global/internal/repositories"
)

type TerritoryHandler struct {
	DB *gorm.DB
}

func NewTerritoryHandler(db *gorm.DB) *TerritoryHandler {
	return &TerritoryHandler{DB: db}
}

// GetEstados godoc
// @Summary      Obtener todos los estados
// @Description  Retorna la lista de todos los estados de Venezuela
// @Tags         Territorio
// @Produce      json
// @Success      200  {array}   dto.EstadoResponse
// @Router       /api/v1/territorio/estados [get]
func (h *TerritoryHandler) GetEstados(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	estados, err := repositories.GetEstados(h.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo estados"})
		return
	}
	json.NewEncoder(w).Encode(estados)
}

// GetMunicipios godoc
// @Summary      Obtener municipios por estado
// @Description  Retorna la lista de municipios pertenecientes a un estado específico
// @Tags         Territorio
// @Produce      json
// @Param        estado_id path int true "ID del Estado"
// @Success      200  {array}   dto.MunicipioResponse
// @Router       /api/v1/territorio/estados/{estado_id}/municipios [get]
func (h *TerritoryHandler) GetMunicipios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	estadoID, _ := strconv.Atoi(chi.URLParam(r, "estado_id"))

	municipios, err := repositories.GetMunicipiosByEstado(h.DB, estadoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(municipios)
}

// GetParroquias godoc
// @Summary      Obtener parroquias por municipio
// @Description  Retorna la lista de parroquias pertenecientes a un municipio específico
// @Tags         Territorio
// @Produce      json
// @Param        municipio_id path int true "ID del Municipio"
// @Success      200  {array}   dto.ParroquiaResponse
// @Router       /api/v1/territorio/municipios/{municipio_id}/parroquias [get]
func (h *TerritoryHandler) GetParroquias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	municipioID, _ := strconv.Atoi(chi.URLParam(r, "municipio_id"))

	parroquias, err := repositories.GetParroquiasByMunicipio(h.DB, municipioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(parroquias)
}
