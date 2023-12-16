package handlers

import (
	"encoding/json"
	"movies-proy/pkg/models"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
)

func TestGetAllMoviesHandler(t *testing.T) {
	// Crear un mock de la base de datos
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("An error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	// Configurar las expectativas del mock
	rows := sqlmock.NewRows([]string{"id", "title", "release_date", "runtime", "mpaa_rating", "description", "image", "created_at", "updated_at"}).
		AddRow(1, "Inception", time.Now(), 148, "PG-13", "A thief...", "image.jpg", time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM movies$").WillReturnRows(rows)

	// Crear una solicitud HTTP y un ResponseRecorder
	req, err := http.NewRequest("GET", "/api/movies", nil)
	if err != nil {
		t.Fatal(err)
	}
	rr := httptest.NewRecorder()

	// Crear e invocar el handler
	handler := GetAllMovies(db)
	handler.ServeHTTP(rr, req)

	// Comprobar el status code y la respuesta
	assert.Equal(t, http.StatusOK, rr.Code, "Expected status code 200")
	//models.Movie
	// Comprobar el cuerpo de la respuesta
	var movies []models.Movie
	err = json.Unmarshal(rr.Body.Bytes(), &movies)
	assert.NoError(t, err, "Error al deserializar la respuesta")

	// Comprobar que la respuesta contiene la información correcta
	assert.Len(t, movies, 1, "Expected 1 movie in the response")
	assert.Equal(t, "Inception", movies[0].Title, "Expected title to match")
	assert.Equal(t, 148, movies[0].Runtime, "Expected runtime to match")
	assert.Equal(t, "PG-13", movies[0].MPAArating, "Expected MPAA rating to match")
}
