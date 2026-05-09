package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
	httpSwagger "github.com/swaggo/http-swagger"

	"api-global/docs"
	"api-global/internal/config"
	"api-global/internal/database"
	"api-global/internal/handlers"

	_ "api-global/docs"
)

// @title           API Global Universitas
// @version         1.0
// @description     API centralizada para indicadores economicos y territoriales.
// @contact.name    Luiger
// @BasePath        /
func main() {
	cfg := config.LoadConfig()

	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("No se pudo iniciar la base de datos: %v", err)
	}

	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Use(cors.Handler(cors.Options{
		AllowedOrigins:   cfg.CorsOrigins,
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300,
	}))

	ecoHandler := handlers.NewEconomicHandler(db)
	terrHandler := handlers.NewTerritoryHandler(db)

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
			r.Get("/municipios/{municipio_id}/ciudades", terrHandler.GetCiudades)
		})
	})

	r.Get("/api/docs", func(w http.ResponseWriter, r *http.Request) {
		http.Redirect(w, r, "/api/docs/index.html", http.StatusMovedPermanently)
	})

	r.Get("/api/docs/*", httpSwagger.Handler(
		httpSwagger.URL("/api/docs/doc.json"),
	))

	docs.SwaggerInfo.Host = ""

	log.Printf("Servidor corriendo en el puerto %s", cfg.Port)
	log.Printf("Documentacion Swagger en: http://localhost:%s/api/docs", cfg.Port)

	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
	if err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
