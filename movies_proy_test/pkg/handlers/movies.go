package handlers

import (
	"database/sql"
	"encoding/json"
	"log"
	"movies-proy/pkg/models"
	"net/http"
	"strconv"
	"strings"
)

// GetAllMovies maneja solicitudes GET para obtener todas las películas.
func GetAllMovies(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		movies, err := models.GetAllMovies(db)
		if err != nil {
			http.Error(w, "Server error", http.StatusInternalServerError)
			log.Printf("Error getting movies: %v", err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movies)
	}
}

// GetMovieByID maneja solicitudes GET para obtener una película por su ID.
func GetMovieByID(db *sql.DB) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		pathParts := strings.Split(r.URL.Path, "/")
		if len(pathParts) < 3 {
			http.Error(w, "Invalid request, movie ID is missing", http.StatusBadRequest)
			return
		}
		id, err := strconv.Atoi(pathParts[2])
		if err != nil {
			http.Error(w, "Invalid movie ID format", http.StatusBadRequest)
			return
		}

		movie, err := models.GetMovieByID(db, id)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "Movie not found", http.StatusNotFound)
			} else {
				http.Error(w, "Server error", http.StatusInternalServerError)
			}
			log.Printf("Error getting movie with ID %d: %v", id, err)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(movie)
	}
}

// Aquí puedes añadir más manejadores como CreateMovie, UpdateMovie, DeleteMovie, etc.
