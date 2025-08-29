package repository

import (
	"bookstore-api/model"
	"database/sql"
	"log"
	"testing"

	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestDB(t *testing.T) *sql.DB{
	connStr := "postgres://postgres:postgres@localhost:5433/bookstore_db_test?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil{
		log.Fatalf("Failed to connect to test database: %v", err)
	}

	createTableQuery := `CREATE TABLE IF NOT EXISTS books (
				id SERIAL PRIMARY KEY,
				title VARCHAR(255) NOT NULL,
				author VARCHAR(255) NOT NULL,
				description TEXT);`

	_, err = db.Exec(createTableQuery)
	if err != nil {
		log.Fatalf("Failed to create books table: %v", err)
	}

	db.Exec("DELETE FROM books")

	return db
}

func TestCreateBook(t *testing.T){
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBookRepository(db)

	// test data
	book := model.Book{
		Title: "Test Book",
		Author: "Test Author",
		Description: "This is a test book",
	}

	bookID, err := repo.CreateBook(book)
	if err != nil{
		t.Fatalf("CreateBook() failed: %v", err)
	}

	if bookID == 0{
		t.Fatalf("Expected book ID to be non-zero, but got 0")
	}

	var count int
	db.QueryRow("SELECT COUNT(*) FROM books WHERE id = $1", bookID).Scan(&count)
	if count != 1 {
		t.Fatalf("Expected to find 1 book with ID %d in the database, but found %d", bookID, count)
	}
}