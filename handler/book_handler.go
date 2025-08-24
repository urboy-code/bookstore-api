package handler

import (
	"bookstore-api/model"
	"bookstore-api/repository"
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type BookHandler struct {
	repo *repository.BookRepository
}

func NewBookHandler(repo *repository.BookRepository) *BookHandler{
	return &BookHandler{repo: repo}
}

// CreateBookHandler adalah fungsi untuk menangani permintaan pembuatan buku baru
// @Summary Create a new book
// @Description Create a new book with the provided details to the database
// @Tags books
// @Accept json
// @Produce json
// @Param book body model.BookInput true "Book Input"
// @Success 201 {object} model.Book
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /books [post]
func (h *BookHandler) CreateBookHandler(c *gin.Context){
	var input model.Book

	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
	}

	bookID, err := h.repo.CreateBook(input)
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to create book"})
		return
	}

	input.ID = bookID

	c.JSON(http.StatusCreated, input)

}

// GetBooksHandler adalah fungsi untuk menangani permintaan mendapatkan semua daftar buku
// @Summary Get all books
// @Description Get all books from the database
// @Tags books
// @Accept json
// @Produce json
// @Success 200 {array} model.Book
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /books [get]
func (h *BookHandler) GetBooksHandler(c *gin.Context){
	books, err := h.repo.GetBooks()
	if err != nil{
		c.JSON(http.StatusInternalServerError, gin.H{"error" : "Failed to retrieve books"})
		return
	}

	c.JSON(http.StatusOK, books)
}	

// GetBookByIDHandler adalah fungsi untuk menangani permintaan mendapatkan buku berdasarkan ID

func (h *BookHandler) GetBookByIDHandler(c *gin.Context){
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if (err != nil){
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	book, err := h.repo.GetBookByID(id)
	if err != nil{
		if err == sql.ErrNoRows{
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve book"})
		return
	}

	c.JSON(http.StatusOK, book)
}