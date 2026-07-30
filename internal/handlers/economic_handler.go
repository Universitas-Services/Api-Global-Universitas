package handlers

import (
	"api-global/internal/dto"
	"api-global/internal/repositories"
	"api-global/internal/scrapers"
	"api-global/internal/timeutil"
	"encoding/json"
	"net/http"
	"time"

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

// GetBCV godoc
// @Summary      Obtener tasas del BCV (Dólar y Euro)
// @Description  Retorna el valor actual del Dólar y Euro oficial.
// @Tags         Economia
// @Produce      json
// @Success      200  {object}  dto.GenericResponse{data=dto.BCVResponse}
// @Failure      500  {object}  dto.GenericResponse
// @Router       /api/v1/economia/bcv [get]
func (h *EconomicHandler) GetBCV(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	usdHoy, errUSD := repositories.GetIndicadorDeHoy(h.DB, "USD_BCV")
	eurHoy, errEUR := repositories.GetIndicadorDeHoy(h.DB, "EUR_BCV")

	if errUSD != nil || errEUR != nil {
		rates, err := scrapers.ScrapeBCV()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Error obteniendo datos del BCV en vivo: " + err.Error()})
			return
		}

		_ = repositories.SaveIndicador(h.DB, "USD_BCV", rates.USD)
		_ = repositories.SaveIndicador(h.DB, "EUR_BCV", rates.EUR)

		w.WriteHeader(http.StatusOK)
		json.NewEncoder(w).Encode(dto.GenericResponse{
			Message: "Tasas del BCV extraídas y actualizadas con éxito",
			Data: dto.BCVResponse{
				USD:   rates.USD,
				EUR:   rates.EUR,
				Fecha: timeutil.HoyCaracas(),
			},
		})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Tasas del BCV obtenidas con éxito (Desde Caché)",
		Data: dto.BCVResponse{
			USD:   usdHoy.Valor,
			EUR:   eurHoy.Valor,
			Fecha: usdHoy.Fecha,
		},
	})
}

// GetBCVHistorico godoc
// @Summary      Obtener tasa BCV por fecha
// @Description  Retorna USD y EUR del BCV para una fecha específica desde la base de datos (sin scrapeo). Pensado para el calendario del frontend.
// @Tags         Economia
// @Produce      json
// @Param        fecha query string true "Fecha a consultar (YYYY-MM-DD)" example:"2024-05-15"
// @Success      200  {object}  dto.GenericResponse{data=dto.BCVHistoricoResponse}
// @Failure      400  {object}  dto.GenericResponse
// @Failure      404  {object}  dto.GenericResponse
// @Router       /api/v1/economia/bcv/historico [get]
func (h *EconomicHandler) GetBCVHistorico(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")

	fechaStr := r.URL.Query().Get("fecha")
	if fechaStr == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "El parámetro fecha es obligatorio (YYYY-MM-DD)"})
		return
	}

	fecha, err := time.Parse("2006-01-02", fechaStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "Formato de fecha inválido. Use YYYY-MM-DD"})
		return
	}

	usd, eur, err := repositories.GetBCVByFecha(h.DB, fecha)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(dto.GenericResponse{Message: "No hay tasas BCV registradas para la fecha " + fechaStr})
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(dto.GenericResponse{
		Message: "Tasa BCV obtenida con éxito",
		Data: dto.BCVHistoricoResponse{
			Fecha: fechaStr,
			USD:   usd,
			EUR:   eur,
		},
	})
}
