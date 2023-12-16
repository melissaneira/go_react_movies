package main

import (
	"database/sql"
	"log"
	"net/http"

	"movies-proy/internal/config"
	"movies-proy/internal/db"
	"movies-proy/pkg/handlers"
)

// AdaptHandler adapta un handler que requiere un *sql.DB a un http.HandlerFunc
func AdaptHandler(db *sql.DB, h func(db *sql.DB) http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		h(db)(w, r)
	}
}

func main() {
	// Cargar la configuración
	cfg, err := config.LoadConfig()
	if err != nil {
		log.Fatalf("Error loading configuration: %v", err)
	}

	// Inicializar la conexión a la base de datos
	database, err := db.OpenDB(cfg.DatabaseURL)
	if err != nil {
		log.Fatalf("Error opening database: %v", err)
	}
	defer database.Close()

	// Configurar el router
	router := http.NewServeMux()

	// Asignar manejadores a rutas específicas, utilizando el adaptador
	router.HandleFunc("/api/movies", AdaptHandler(database, handlers.GetAllMovies))
	router.HandleFunc("/api/movies/", AdaptHandler(database, handlers.GetMovieByID)) // Nota: la ruta incluye una barra al final

	// Iniciar el servidor
	log.Printf("Starting server on %s", cfg.ServerAddress)
	if err := http.ListenAndServe(cfg.ServerAddress, router); err != nil {
		log.Fatalf("Error starting server: %v", err)
	}
}
