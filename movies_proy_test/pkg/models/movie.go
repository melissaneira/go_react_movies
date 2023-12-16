package models

import (
	"database/sql"
	"fmt"
	"time"
)

// Movie representa un registro de película en la base de datos.
type Movie struct {
	ID          int       `json:"id"`
	Title       string    `json:"title"`
	ReleaseDate time.Time `json:"release_date"`
	Runtime     int       `json:"runtime"`
	MPAArating  string    `json:"mpaa_rating"`
	Description string    `json:"description"`
	Image       string    `json:"image"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
}

// GetAllMovies devuelve todas las películas de la base de datos.
func GetAllMovies(db *sql.DB) ([]Movie, error) {
	const query = `SELECT id, title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at FROM movies`
	rows, err := db.Query(query)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var movies []Movie
	for rows.Next() {
		var m Movie
		if err := rows.Scan(&m.ID, &m.Title, &m.ReleaseDate, &m.Runtime, &m.MPAArating, &m.Description, &m.Image, &m.CreatedAt, &m.UpdatedAt); err != nil {
			return nil, err
		}
		movies = append(movies, m)
	}
	return movies, nil
}

// GetMovieByID devuelve una película por su ID.
func GetMovieByID(db *sql.DB, id int) (Movie, error) {
	const query = `SELECT id, title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at FROM movies WHERE id = $1`
	var m Movie
	row := db.QueryRow(query, id)
	if err := row.Scan(&m.ID, &m.Title, &m.ReleaseDate, &m.Runtime, &m.MPAArating, &m.Description, &m.Image, &m.CreatedAt, &m.UpdatedAt); err != nil {
		if err == sql.ErrNoRows {
			return Movie{}, fmt.Errorf("no movie found with id: %d", id)
		}
		return Movie{}, err
	}
	return m, nil
}

// CreateMovie crea una nueva película en la base de datos.
func CreateMovie(db *sql.DB, m Movie) error {
	const query = `INSERT INTO movies (title, release_date, runtime, mpaa_rating, description, image, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8)`
	_, err := db.Exec(query, m.Title, m.ReleaseDate, m.Runtime, m.MPAArating, m.Description, m.Image, m.CreatedAt, m.UpdatedAt)
	return err
}

// UpdateMovie actualiza una película existente en la base de datos.
func UpdateMovie(db *sql.DB, m Movie) error {
	const query = `UPDATE movies SET title = $1, release_date = $2, runtime = $3, mpaa_rating = $4, description = $5, image = $6, updated_at = $7 WHERE id = $8`
	_, err := db.Exec(query, m.Title, m.ReleaseDate, m.Runtime, m.MPAArating, m.Description, m.Image, time.Now(), m.ID)
	return err
}

// DeleteMovie elimina una película de la base de datos.
func DeleteMovie(db *sql.DB, id int) error {
	const query = `DELETE FROM movies WHERE id = $1`
	_, err := db.Exec(query, id)
	return err
}
