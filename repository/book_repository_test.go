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
		t.Fatalf("Failed to connect to test database: %v", err)
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

func TestGetBooks(t *testing.T){
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBookRepository(db)

	booksToInsert := []model.Book{
		{Title: "Book One", Author: "Author One", Description: "First book"},
	}

	for _, book := range booksToInsert{
		_, err := repo.CreateBook(book)
		if err != nil{
			t.Fatalf("Failed to insert book: %v", err)
		}

		books, err := repo.GetBooks()

		if err != nil {
			t.Fatalf("GetBooks() failed: %v", err)
		}
		if len(books) != len(booksToInsert){
			t.Fatalf("Expected %d books, but got %d", len(booksToInsert), len(books))
		}
	}
}

func TestGetBookById(t *testing.T){
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBookRepository(db)

	// test data
	book := model.Book{
		Title: "Test Book",
		Author: "Test Author",
		Description: "This is a test book",
	}

	bookID, _ := repo.CreateBook(book)

	foundBook, err := repo.GetBookByID(bookID)
	if err != nil{
		t.Fatalf("GetBookByID() for existing ID failed: %v", err)
	}

	if foundBook.Title != "Test Book" {
		t.Errorf("Expected title 'Test Book', but got '%s'", foundBook.Title)
	}

	_, err = repo.GetBookByID(9999)
	if err != ErrBookNotFound{
		t.Errorf("Expected ErrBookNotFound for non-existing ID, but got: %v", err)
	}
}

func TestUpdateBook(t *testing.T){
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBookRepository(db)

	// test data
	book := model.Book{
		Title: "Original Title",
		Author: "Original Author",
		Description: "Original Description",
	}

	bookID, _ := repo.CreateBook(book)

	book.ID = bookID
	book.Title = "Update Title"
	book.Author = "Update Author"
	book.Description = "Update Description"

	err := repo.UpdateBook(bookID, book)
	if err != nil{
		t.Fatalf("UpdateBook() failed: %v", err)
	}

	updateBook, err := repo.GetBookByID(bookID)
	if err != nil{
		t.Fatalf("GetBookByID() after update failed: %v", err)
	}

	if updateBook.Title != "Update Title" || updateBook.Author != "Update Author" || updateBook.Description != "Update Description"{
		t.Errorf("Book not updated correctly. Got: %+v", updateBook)
	}
}

func TestDeleteBook(t *testing.T){
	db := setupTestDB(t)
	defer db.Close()

	repo := NewBookRepository(db)

	// test data
	book := model.Book{
		Title: "To be deleted",
		Author: "To be deleted",
		Description: "To be deleted",
	}

	bookID, _ := repo.CreateBook(book)
	err := repo.DeleteBook(bookID)
	if err != nil{
		t.Fatalf("DeleteBook() failed: %v", err)
	}

	_, err = repo.GetBookByID(bookID)
	if err != ErrBookNotFound{
		t.Errorf("Expected ErrBookNotFound after deletion, but got: %v", err)
	}
	
}