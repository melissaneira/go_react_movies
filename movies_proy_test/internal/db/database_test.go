package db

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

const connectionString = "host=localhost port=5433 dbname=reactgo user=postgres password=password sslmode=disable"

func TestOpenDB(t *testing.T) {

	db, err := OpenDB(connectionString)

	if err != nil {
		t.Fatalf("error opening db: %v", err)
	}

	err = db.Ping()

	if err != nil {
		t.Fatal("error pinging db: " + err.Error())
	}

}

func TestQuery(t *testing.T) {

	db, err := OpenDB(connectionString)
	if err != nil {
		t.Fatal("failed to open db")
	}
	defer db.Close()

	// Ejecuta un query
	row := db.QueryRow("SELECT title, created_at FROM movies WHERE id = $1", 2)

	var title string
	var created time.Time

	err = row.Scan(&title, &created)

	if err != nil {
		t.Fatal("failed to query table:", err)
	}

	// Hacer assertions sobre los datos retornados
	if title != "Raiders of the Lost Ark" {
		t.Errorf("invalid name: %s", title)
	}

	// Verificar fecha de creación
	if created.Year() != 2022 {
		t.Errorf("invalid creation year: %d", created.Year())
	}

}

func TestCloseDB(t *testing.T) {

	db, err := OpenDB(connectionString)

	assert.Nil(t, err, "failed to open db")

	err = db.Ping()
	assert.Nil(t, err, "couldn't ping db")

	db.Close()

	err = db.Ping()
	assert.NotNil(t, err, "expected error pinging closed db")

	assert.Error(t, err, "expected error pinging closed db")

}

func TestParameterizedQueries(t *testing.T) {
	db, err := OpenDB(connectionString)
	if err != nil {
		t.Fatal("failed to connect to db")
	}
	defer db.Close()

	// Inserta registro de prueba
	result, err := db.Exec(`INSERT INTO movies(title, runtime) VALUES($1, $2)`, "Lego", "180")

	if err != nil {
		t.Fatal("failed to insert row:", err)
	}

	// Verifica filas afectadas
	count, err := result.RowsAffected()
	assert.Equal(t, int64(1), count)

	// Query parametrizado
	row := db.QueryRow("SELECT title, runtime FROM movies WHERE title=$1", "Lego")

	var fetchedTitle string
	var fetchedRuntime int

	err = row.Scan(&fetchedTitle, &fetchedRuntime)
	assert.NoError(t, err)

	// Validar datos
	assert.Equal(t, "Lego", fetchedTitle)
	assert.Equal(t, 180, fetchedRuntime)

}

func TestTransactionCommit(t *testing.T) {
	db, err := OpenDB(connectionString)
	if err != nil {
		t.Fatalf("failed to open db: %v", err)
	}
	defer db.Close()

	// Iniciar una transacción
	tx, err := db.Begin()
	if err != nil {
		t.Fatalf("failed to begin transaction: %v", err)
	}

	// Ejecutar algunas operaciones dentro de la transacción
	_, err = tx.Exec("INSERT INTO movies(title, runtime) VALUES($1, $2)", "New Movie", 120)
	if err != nil {
		tx.Rollback()
		t.Fatalf("failed to execute insert: %v", err)
	}

	// Confirmar la transacción
	err = tx.Commit()
	if err != nil {
		t.Fatalf("failed to commit transaction: %v", err)
	}

	// Verificar que los cambios se reflejan en la base de datos
	var runtime int
	err = db.QueryRow("SELECT runtime FROM movies WHERE title = $1", "New Movie").Scan(&runtime)
	if err != nil {
		t.Errorf("failed to query inserted movie: %v", err)
	}

	if runtime != 120 {
		t.Errorf("expected runtime 120, got %d", runtime)
	}
}
