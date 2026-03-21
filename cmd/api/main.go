package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	httpSwagger "github.com/swaggo/http-swagger"

	"api-global/internal/config"
	"api-global/internal/database"
	"api-global/internal/handlers"
	"api-global/internal/models"

	_ "api-global/docs"
)

// @title           API Global Universitas
// @version         1.0
// @description     API centralizada para indicadores económicos y territoriales.
// @contact.name    Jose
// @host            localhost:8080
// @BasePath        /
func main() {
	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("No se pudo iniciar la base de datos: %v", err)
	}

	// 1. Ejecutar Migraciones
	db.AutoMigrate(&models.Estado{}, &models.Municipio{}, &models.Parroquia{}, &models.IndicadorEconomico{})

	// 2. Ejecutar Seeder (Poblar BD)
	database.SeedTerritories(db)

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// Instanciar Handlers
	ecoHandler := handlers.NewEconomicHandler(db)
	terrHandler := handlers.NewTerritoryHandler(db)

	// ==========================================
	// RUTAS DE LA API
	// ==========================================
	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/economia", func(r chi.Router) {
			r.Get("/ucauu", ecoHandler.GetUCAUU)
			r.Post("/ucauu", ecoHandler.CreateUCAUU)
			r.Get("/bcv", ecoHandler.GetBCV)
		})

		r.Route("/territorio", func(r chi.Router) {
			r.Get("/estados", terrHandler.GetEstados)
			r.Get("/estados/{estado_id}/municipios", terrHandler.GetMunicipios)
			r.Get("/municipios/{municipio_id}/parroquias", terrHandler.GetParroquias)
		})
	})

	// ==========================================
	// CONFIGURACIÓN SWAGGER
	// ==========================================
	// Redirección limpia (sin index.html en la URL)
	r.Get("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/docs/index.html", http.StatusMovedPermanently)
	})

	r.Get("/api/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/api/docs/doc.json"),
	))

	// ==========================================
	// INICIO DEL SERVIDOR
	// ==========================================
	log.Printf("🚀 Servidor corriendo en el puerto %s", cfg.Port)
	log.Printf("📚 Documentación Swagger en: http://localhost:%s/api/docs", cfg.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
	if err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
