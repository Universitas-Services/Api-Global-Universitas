package dto

import "time"

// ==========================================
// DTOs para Peticiones (Request)
// ==========================================

// CreateUCAUURequest es el JSON que esperamos recibir en el POST
type CreateUCAUURequest struct {
	// validate:"required,gt=0" asegura que el campo venga y sea mayor a cero
	Valor float64 `json:"valor" validate:"required,gt=0" example:"35.50"`
}

// ==========================================
// DTOs para Respuestas (Response)
// ==========================================

// UCAUUResponse es lo que le devolveremos al Frontend
type UCAUUResponse struct {
	Valor float64   `json:"valor" example:"35.50"`
	Fecha time.Time `json:"fecha" example:"2026-03-20T10:00:00Z"`
}

// GenericResponse estructura estándar para respuestas de éxito o error
type GenericResponse struct {
	Message string      `json:"message" example:"Operación exitosa"`
	Data    interface{} `json:"data,omitempty"`
}

// BCVResponse estructura la respuesta del endpoint de tasas BCV
type BCVResponse struct {
	USD   float64   `json:"usd" example:"36.4521"`
	EUR   float64   `json:"eur" example:"39.1234"`
	Fecha time.Time `json:"fecha" example:"2026-03-20T00:00:00Z"`
}

// BCVHistoricoResponse es la tasa BCV de una fecha concreta (calendario)
type BCVHistoricoResponse struct {
	Fecha string  `json:"fecha" example:"2024-05-15"`
	USD   float64 `json:"usd" example:"36.4521"`
	EUR   float64 `json:"eur" example:"39.1234"`
}

// BCVCaptureResponse es el resultado de una captura forzada (Cloud Scheduler)
type BCVCaptureResponse struct {
	Fecha  string  `json:"fecha" example:"2026-07-30"`
	USD    float64 `json:"usd" example:"745.6371"`
	EUR    float64 `json:"eur" example:"848.8258"`
	Status string  `json:"status" example:"updated"` // created | updated | unchanged
}

