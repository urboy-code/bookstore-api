package handler

import (
	"bookstore-api/model"
	"bookstore-api/repository"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"bytes"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func setupTestRouter() (*gin.Engine, *repository.BookRepository){
	gin.SetMode(gin.TestMode)

	connStr := "postgres://postgres:postgres@localhost:5433/bookstore_db_test?sslmode=disable"
	db, err := sql.Open("pgx", connStr)
	if err != nil{
		panic("Failed to connect to test database")
	}
	db.Exec("DELETE FROM books")

	repo := repository.NewBookRepository(db)
	handler := NewBookHandler(repo)

	router := gin.Default()
	router.GET("/books", handler.GetBooksHandler)
	router.GET("/books/:id", handler.GetBookByIDHandler)
	router.POST("/books", handler.CreateBookHandler)

	return router, repo
}

// TestGetBooksHandler tests the GetBooksHandler function
// It tests the following cases:
// - an empty list of books is returned when the database is empty
// - a single book is returned when there is only one book in the database
//
// The test uses the setupTestRouter function to create a test router and database
// connection. It then creates a test request to the "/books" endpoint, and uses
// the httptest.NewRecorder to record the response. The response code and body
// are then checked to ensure they match the expected values.
func TestGetBooksHandler(t *testing.T){
	router, repo := setupTestRouter()

	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodGet, "/books", nil)

	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusOK{
		t.Errorf("Expected status code %d but got %d", http.StatusOK, recorder.Code)
	}

	if recorder.Body.String() != "[]"{
		t.Errorf("Expected empty array '[]' but got %s", recorder.Body.String())
	}

	repo.CreateBook(model.Book{
		Title: "Test Book",
		Author: "Test Author",
	})

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodGet, "/books", nil)

	router.ServeHTTP(recorder, request)

	books := make([]model.Book, 0)
	json.Unmarshal(recorder.Body.Bytes(), &books)

	if len(books) != 1{
		t.Errorf("Expected 1 book but got %d", len(books))
	}

	if books[0].Title != "Test Book"{
		t.Errorf("Expected book title 'Test Book' but got '%s'", books[0].Title)
	}
}

func TestGetBookByIDHandler(t *testing.T){
	router, repo := setupTestRouter()

	// Test case 1: Book Found
	createBook := model.Book{
		Title: "Test Book",
		Author: "Test Author",
	}
	bookID, _ := repo.CreateBook(createBook)

	// Create request for existing book
	recorder := httptest.NewRecorder()
	// Format the URL with the book ID
	requestURL := fmt.Sprintf("/books/%d", bookID)
	request, _ := http.NewRequest(http.MethodGet, requestURL, nil)

	// Run the request
	router.ServeHTTP(recorder, request)

	// Chect the result of the request
	if recorder.Code != http.StatusOK{
		t.Errorf("Expected status code 200 but got %d", recorder.Code)
	}

	var foundBook model.Book
	json.Unmarshal(recorder.Body.Bytes(), &foundBook)
	if foundBook.Title != "Test Book" {
		t.Errorf("Expected book title 'Test Book', but got '%s'", foundBook.Title)
	}

	// Test case 2: Book Not Found
	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodGet, "/books/999", nil)
	
	router.ServeHTTP(recorder, request)

	if recorder.Code != http.StatusNotFound{
		t.Errorf("Expected status code 404 but got %d", recorder.Code)
	}
}

func TestCreateBookHandler(t *testing.T){
	// This function can be implemented similarly to the other test functions
	// by setting up the test router, creating a POST request with a book payload,
	// and checking the response for correctness.

	router, _ := setupTestRouter()
	
	// Test case: Valid Book Creation
	bookPayload := `{"title": "Test Book", "author": "Test Author", "description": "Test Description"}`
	bodyHeader := bytes.NewReader([]byte(bookPayload))

	// Create request Post with body
	recorder := httptest.NewRecorder()
	request, _ := http.NewRequest(http.MethodPost, "/books", bodyHeader)
	request.Header.Set("Content-Type", "application/json")

	// Run the request
	router.ServeHTTP(recorder, request)

	// Check the result of the request
	if recorder.Code != http.StatusCreated{
		t.Errorf("Expected status code 201 but got %d", recorder.Code)
	}

	var createdBook model.Book
	json.Unmarshal(recorder.Body.Bytes(), &createdBook)
	if createdBook.Title != "Test Book"{
		t.Errorf("Expected book title 'Test Book', but got '%s'", createdBook.Title)
	}

	// Test case: Invalid Book Creation (missing title)
	invalidPayload := `{"author": "Test Author", "description": "Test Description"}`
	bodyHeader = bytes.NewReader([]byte(invalidPayload))

	recorder = httptest.NewRecorder()
	request, _ = http.NewRequest(http.MethodPost, "/books", bodyHeader)
	request.Header.Set("Content-Type", "application/json")

	router.ServeHTTP(recorder, request)
	if recorder.Code != http.StatusBadRequest{
		t.Errorf("Expected status code 400 but got %d", recorder.Code)
	}
	
}