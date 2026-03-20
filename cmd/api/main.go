package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
)

func main() {
	// 1. Inicializar el enrutador chi
	r := chi.NewRouter()

	// 2. Añadir Middlewares globales integrados en chi
	r.Use(middleware.Logger)    // Registra cada petición HTTP en la consola
	r.Use(middleware.Recoverer) // Evita que la API se caiga si hay un "panic" (error crítico)

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡Gogeta!!"))
	})

	// 3. Definir una ruta de prueba (Healthcheck)
	r.Get("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("¡La API Global está viva!"))
	})

	// 4. Iniciar el servidor
	puerto := ":8080"
	fmt.Printf("Servidor corriendo en el puerto %s\n", puerto)

	// ListenAndServe es de la librería estándar de Go, pero le pasamos 'r' (nuestro router chi)
	err := http.ListenAndServe(puerto, r)
	if err != nil {
		log.Fatalf("Error al iniciar el servidor: %v", err)
	}
}
