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
		Message: "Estados obtenidos con éxito",
		Data:    response,
	})
}

// GetMunicipios godoc
// @Summary      Obtener municipios por estado
// @Description  Retorna la lista de municipios pertenecientes a un estado específico
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
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "ID de estado inválido"})
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
		Message: "Municipios obtenidos con éxito",
		Data:    response,
	})
}

// GetParroquias godoc
// @Summary      Obtener parroquias por municipio
// @Description  Retorna la lista de parroquias pertenecientes a un municipio específico
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
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "ID de municipio inválido"})
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
		Message: "Parroquias obtenidas con éxito",
		Data:    response,
	})
}

// GetTribunalesEstadales godoc
// @Summary      Obtener tribunales estadales por estado
// @Description  Retorna la lista de tribunales de categoría Superior y Primera Instancia pertenecientes a un estado
// @Tags         Territorio
// @Produce      json
// @Param        estado_id path int true "ID del Estado"
// @Success      200  {object}  dto.GenericResponse{data=[]dto.TribunalResponse}
// @Failure      400  {object}  dto.GenericResponse
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/territorio/estados/{estado_id}/tribunales [get]
func (h *TerritoryHandler) GetTribunalesEstadales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	estadoID, err := strconv.Atoi(chi.URLParam(r, "estado_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "ID de estado inválido"})
		return
	}

	tribunales, err := repositories.GetTribunalesByEstado(h.DB, estadoID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo tribunales estadales"})
		return
	}

	var response []dto.TribunalResponse
	for _, t := range tribunales {
		response = append(response, dto.TribunalResponse{ID: t.ID, Nombre: t.Nombre, Categoria: t.Categoria})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Tribunales estadales obtenidos con éxito",
		Data:    response,
	})
}

// GetTribunalesMunicipales godoc
// @Summary      Obtener tribunales municipales por municipio
// @Description  Retorna la lista de tribunales de categoría Municipio asociados a un municipio específico, incluyendo los compartidos con otros municipios
// @Tags         Territorio
// @Produce      json
// @Param        municipio_id path int true "ID del Municipio"
// @Success      200  {object}  dto.GenericResponse{data=[]dto.TribunalResponse}
// @Failure      400  {object}  dto.GenericResponse
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/territorio/municipios/{municipio_id}/tribunales [get]
func (h *TerritoryHandler) GetTribunalesMunicipales(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	municipioID, err := strconv.Atoi(chi.URLParam(r, "municipio_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "ID de municipio inválido"})
		return
	}

	tribunales, err := repositories.GetTribunalesByMunicipio(h.DB, municipioID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo tribunales municipales"})
		return
	}

	var response []dto.TribunalResponse
	for _, t := range tribunales {
		response = append(response, dto.TribunalResponse{ID: t.ID, Nombre: t.Nombre, Categoria: t.Categoria})
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Tribunales municipales obtenidos con éxito",
		Data:    response,
	})
}
