package handler

import (
	"bookstore-api/model"
	"bookstore-api/repository"
	"database/sql"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

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

	return router, repo
}

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