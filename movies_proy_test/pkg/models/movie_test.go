package models

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGetAllMovies(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	rows := sqlmock.NewRows([]string{"id", "title", "release_date", "runtime", "mpaa_rating", "description", "image", "created_at", "updated_at"}).
		AddRow(1, "Inception", time.Now(), 148, "PG-13", "A thief...", "image.jpg", time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM movies$").WillReturnRows(rows)

	movies, err := GetAllMovies(db)
	if err != nil {
		t.Errorf("error was not expected while fetching data: %s", err)
	}

	if len(movies) != 1 {
		t.Errorf("expected a single movie, got %d", len(movies))
	}

	assert.NoError(t, err)

	assert.Equal(t, 1, len(movies))

	movie := movies[0]

	assert.Equal(t, 1, movie.ID)
	assert.Equal(t, "Inception", movie.Title)
	assert.WithinDuration(t, time.Now(), movie.ReleaseDate, time.Second)
	assert.Equal(t, 148, movie.Runtime)

}

func TestGetMovieByID(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	row := sqlmock.NewRows([]string{"id", "title", "release_date", "runtime", "mpaa_rating", "description", "image", "created_at", "updated_at"}).
		AddRow(1, "Inception", time.Now(), 148, "PG-13", "A thief...", "image.jpg", time.Now(), time.Now())

	mock.ExpectQuery("^SELECT (.+) FROM movies WHERE id = \\$1$").WithArgs(1).WillReturnRows(row)

	movie, err := GetMovieByID(db, 1)
	if err != nil {
		t.Errorf("error was not expected while fetching data: %s", err)
	}

	if movie.ID != 1 {
		t.Errorf("expected movie ID 1, got %d", movie.ID)
	}

	assert.NoError(t, err, "Fetching movie should not produce an error")
	assert.NotNil(t, movie, "Movie should not be nil")
	assert.Equal(t, 1, movie.ID, "Expected movie ID to be 1")
}

func TestCreateMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("INSERT INTO movies").
		WithArgs("Inception", sqlmock.AnyArg(), 148, "PG-13", "A thief...", "image.jpg", sqlmock.AnyArg(), sqlmock.AnyArg()).
		WillReturnResult(sqlmock.NewResult(1, 1))

	movie := Movie{
		Title:       "Inception",
		ReleaseDate: time.Now(),
		Runtime:     148,
		MPAArating:  "PG-13",
		Description: "A thief...",
		Image:       "image.jpg",
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}

	err = CreateMovie(db, movie)
	if err != nil {
		t.Errorf("error was not expected while creating data: %s", err)
	}

	assert.NoError(t, err, "Creating movie should not produce an error")
}

func TestUpdateMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("UPDATE movies").
		WithArgs("Inception", sqlmock.AnyArg(), 148, "PG-13", "Updated Description", "image.jpg", sqlmock.AnyArg(), 1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	movie := Movie{
		ID:          1,
		Title:       "Inception",
		ReleaseDate: time.Now(),
		Runtime:     148,
		MPAArating:  "PG-13",
		Description: "Updated Description",
		Image:       "image.jpg",
		UpdatedAt:   time.Now(),
	}

	err = UpdateMovie(db, movie)
	if err != nil {
		t.Errorf("error was not expected while updating data: %s", err)
	}
	assert.NoError(t, err, "Updating movie should not produce an error")
}

func TestDeleteMovie(t *testing.T) {
	db, mock, err := sqlmock.New()
	if err != nil {
		t.Fatalf("an error '%s' was not expected when opening a stub database connection", err)
	}
	defer db.Close()

	mock.ExpectExec("DELETE FROM movies").
		WithArgs(1).
		WillReturnResult(sqlmock.NewResult(1, 1))

	err = DeleteMovie(db, 1)
	if err != nil {
		t.Errorf("error was not expected while deleting data: %s", err)
	}

	assert.NoError(t, err, "Deleting movie should not produce an error")

}
