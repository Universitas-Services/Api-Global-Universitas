package handlers

import (
	"api-global/internal/dto"
	"api-global/internal/repositories"
	"encoding/json"
	"net/http"

	"github.com/go-playground/validator/v10"
	"gorm.io/gorm"
)

// Usaremos una única instancia del validador
var validate *validator.Validate

func init() {
	validate = validator.New()
}

// Estructura para inyectar la base de datos a los endpoints
type EconomicHandler struct {
	DB *gorm.DB
}

func NewEconomicHandler(db *gorm.DB) *EconomicHandler {
	return &EconomicHandler{DB: db}
}

// CreateUCAUU godoc
// @Summary      Registrar o Actualizar UCAUU
// @Description  Permite registrar el valor de la UCAUU del día actual. Si ya existe, lo actualiza.
// @Tags         Economia
// @Accept       json
// @Produce      json
// @Param        request body dto.CreateUCAUURequest true "Valor de la UCAUU"
// @Success      201  {object}  dto.GenericResponse
// @Failure      400  {object}  dto.GenericResponse
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/economia/ucauu [post]
func (h *EconomicHandler) CreateUCAUU(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req dto.CreateUCAUURequest

	// 1. Decodificar JSON
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "JSON inválido"})
		return
	}

	// 2. Validar DTO
	if err := validate.Struct(req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "El valor debe ser mayor a 0"})
		return
	}

	// 3. Guardar en BD
	if err := repositories.SaveUCAUU(h.DB, req.Valor); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error guardando en base de datos"})
		return
	}

	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(dto.GenericResponse{Message: "UCAUU registrada correctamente"})
}

// GetUCAUU godoc
// @Summary      Obtener última UCAUU
// @Description  Retorna el último valor registrado de la unidad UCAUU
// @Tags         Economia
// @Produce      json
// @Success      200  {object}  dto.UCAUUResponse
// @Failure      404  {object}  dto.GenericResponse
// @Router       /api/v1/economia/ucauu [get]
func (h *EconomicHandler) GetUCAUU(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	indicador, err := repositories.GetLatestUCAUU(h.DB)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "No hay registros de UCAUU aún"})
		return
	}

	response := dto.UCAUUResponse{
		Valor: indicador.Valor,
		Fecha: indicador.Fecha,
	}
	json.NewEncoder(w).Encode(response)
}
