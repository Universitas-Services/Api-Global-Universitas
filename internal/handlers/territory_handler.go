package handlers

import (
	"api-global/internal/dto"
	"api-global/internal/repositories"
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"
	"gorm.io/gorm"
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
// @Success      200  {object}  dto.GenericResponse{data=[]dto.EstadoResponse}
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/territorio/estados [get]
func (h *TerritoryHandler) GetEstados(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	estados, err := repositories.GetEstados(h.DB)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo estados de la base de datos"})
		return
	}

	var response []dto.EstadoResponse
	for _, est := range estados {
		response = append(response, dto.EstadoResponse{ID: est.ID, Nombre: est.Nombre})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Estados obtenidos con exito",
		Data:    response,
	})
}

// GetMunicipios godoc
// @Summary      Obtener municipios por estado
// @Description  Retorna la lista de municipios pertenecientes a un estado especifico
// @Tags         Territorio
// @Produce      json
// @Param        estado_id path int true "ID del Estado"
// @Success      200  {object}  dto.GenericResponse{data=[]dto.MunicipioResponse}
// @Failure      400  {object}  dto.GenericResponse
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/territorio/estados/{estado_id}/municipios [get]
func (h *TerritoryHandler) GetMunicipios(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	estadoID, err := strconv.Atoi(chi.URLParam(r, "estado_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "ID de estado invalido"})
		return
	}

	municipios, err := repositories.GetMunicipiosByEstado(h.DB, estadoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo municipios"})
		return
	}

	var response []dto.MunicipioResponse
	for _, mun := range municipios {
		response = append(response, dto.MunicipioResponse{ID: mun.ID, EstadoID: mun.EstadoID, Nombre: mun.Nombre})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Municipios obtenidos con exito",
		Data:    response,
	})
}

// GetParroquias godoc
// @Summary      Obtener parroquias por municipio
// @Description  Retorna la lista de parroquias pertenecientes a un municipio especifico
// @Tags         Territorio
// @Produce      json
// @Param        municipio_id path int true "ID del Municipio"
// @Success      200  {object}  dto.GenericResponse{data=[]dto.ParroquiaResponse}
// @Failure      400  {object}  dto.GenericResponse
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/territorio/municipios/{municipio_id}/parroquias [get]
func (h *TerritoryHandler) GetParroquias(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	municipioID, err := strconv.Atoi(chi.URLParam(r, "municipio_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "ID de municipio invalido"})
		return
	}

	parroquias, err := repositories.GetParroquiasByMunicipio(h.DB, municipioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo parroquias"})
		return
	}

	var response []dto.ParroquiaResponse
	for _, parr := range parroquias {
		response = append(response, dto.ParroquiaResponse{ID: parr.ID, MunicipioID: parr.MunicipioID, Nombre: parr.Nombre})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Parroquias obtenidas con exito",
		Data:    response,
	})
}

// GetCiudades godoc
// @Summary      Obtener ciudades por municipio
// @Description  Retorna la lista de ciudades pertenecientes a un municipio especifico
// @Tags         Territorio
// @Produce      json
// @Param        municipio_id path int true "ID del Municipio"
// @Success      200  {object}  dto.GenericResponse{data=[]dto.CiudadResponse}
// @Failure      400  {object}  dto.GenericResponse
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/territorio/municipios/{municipio_id}/ciudades [get]
func (h *TerritoryHandler) GetCiudades(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	municipioID, err := strconv.Atoi(chi.URLParam(r, "municipio_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "ID de municipio invalido"})
		return
	}

	ciudades, err := repositories.GetCiudadesByMunicipio(h.DB, municipioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo ciudades"})
		return
	}

	var response []dto.CiudadResponse
	for _, ciudad := range ciudades {
		response = append(response, dto.CiudadResponse{ID: ciudad.ID, MunicipioID: ciudad.MunicipioID, Nombre: ciudad.Nombre})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Ciudades obtenidas con exito",
		Data:    response,
	})
}
