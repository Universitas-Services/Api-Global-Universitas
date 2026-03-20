package main

import (
	"api-global/internal/config"
	"api-global/internal/database"
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Cargar configuración
	cfg := config.LoadConfig()

	// 2. Conectar a Base de Datos
	db, err := database.InitDB(cfg)
	if err != nil {
		log.Fatalf("No se pudo iniciar la base de datos: %v", err)
	}
	// (Opcional) Guardar db en una variable global o inyectarla luego en los handlers
	_ = db

	// 3. Inicializar Router
	r := chi.NewRouter()
	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)

	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡La API Global está conectada a Docker y Postgres!"))
	})

	// 4. Iniciar Servidor
	log.Printf("🚀 Servidor corriendo en el puerto %s", cfg.Port)
	err = http.ListenAndServe(fmt.Sprintf(":%s", cfg.Port), r)
	if err != nil {
		log.Fatalf("Error en el servidor: %v", err)
	}
}
