package main

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"api-global/internal/config"
	"api-global/internal/database"
	"api-global/internal/handlers"
	"api-global/internal/jobs"
	"api-global/internal/models"

	_ "api-global/docs"
)

// @title           API Global Universitas
// @version         1.0
// @description     API centralizada para indicadores económicos y territoriales.
// @contact.name    Jose
// @BasePath        /
func main() {
	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("No se pudo iniciar la base de datos: %v", err)
	}

	// 1. Ejecutar Migraciones
	db.AutoMigrate(&models.Estado{}, &models.Ciudad{}, &models.Municipio{}, &models.Parroquia{}, &models.IndicadorEconomico{}, &models.Tribunal{}, &models.CodigoArea{})

	// 2. Ejecutar Seeders (Poblar BD)
	database.SeedTerritories(db)
	database.SeedCiudades(db)
	database.SeedTribunales(db)
	database.SeedBCVHistorico(db)
	database.SeedCodigosArea(db)

	// 3. Job BCV: scrape 2x/día (hora Caracas) con upsert
	jobs.StartBCVDailyJob(db, cfg.BCVDailyHours)

	// Inicializar Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	// ==========================================
	// CONFIGURACIÓN DE CORS
	// ==========================================
	corsOrigins := os.Getenv("CORS_ALLOWED_ORIGINS")
	var allowedOrigins []string
	if corsOrigins != "" {
		allowedOrigins = strings.Split(corsOrigins, ",")
	} else {
		allowedOrigins = []string{"http://localhost:3000", "http://localhost:3001"}
	}

	r.Use(cors.Handler(cors.Options{
		// Aquí defines los orígenes permitidos. Luego podrás añadir los dominios de producción.
		AllowedOrigins: allowedOrigins,
		// Métodos permitidos (GET, POST, etc.)
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		// Cabeceras que el frontend tiene permitido enviar
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token", "X-Cron-Secret"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Tiempo en segundos que el navegador cachea esta regla
	}))

	// Instanciar Handlers
	ecoHandler := handlers.NewEconomicHandler(db, cfg.BCVCronSecret)
	terrHandler := handlers.NewTerritoryHandler(db)

	// ==========================================
	// RUTAS DE LA API
	// ==========================================
	r.Route("/api/v1", func(r chi.Router) {

		r.Route("/economia", func(r chi.Router) {
			r.Get("/ucauu", ecoHandler.GetUCAUU)
			r.Post("/ucauu", ecoHandler.CreateUCAUU)
			r.Get("/bcv", ecoHandler.GetBCV)
			r.Get("/bcv/historico", ecoHandler.GetBCVHistorico)
			r.Post("/bcv/capturar", ecoHandler.CaptureBCV)
		})

		r.Route("/territorio", func(r chi.Router) {
			r.Get("/estados", terrHandler.GetEstados)
			r.Get("/estados/{estado_id}/ciudades", terrHandler.GetCiudades)
			r.Get("/estados/{estado_id}/municipios", terrHandler.GetMunicipios)
			r.Get("/estados/{estado_id}/tribunales", terrHandler.GetTribunalesEstadales)
			r.Get("/municipios/{municipio_id}/parroquias", terrHandler.GetParroquias)
			r.Get("/municipios/{municipio_id}/tribunales", terrHandler.GetTribunalesMunicipales)
			r.Get("/codigos-area", terrHandler.GetCodigosArea)
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
