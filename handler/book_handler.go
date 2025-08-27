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
// @Summary Get a book by ID
// @Description Get a book by its ID from the database
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} model.Book
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 404 {object} gin.H{"error": "Book not found"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /books/{id} [get]
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

// UpdateBookHandler adalah fungsi untuk menangani permintaan memperbarui data buku berdasarkan ID
// @Summary Update a book by ID
// @Description Update a book by its ID in the database
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Param book body model.BookInput true "Book Input"
// @Success 200 {object} model.Book
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 404 {object} gin.H{"error": "Book not found"}
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /books/{id} [put]
func (h *BookHandler) UpdateBookHandler(c *gin.Context){
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
		return
	}

	var input model.Book
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input: " + err.Error()})
		return
	}

	err = h.repo.UpdateBook(id, input)
	if err != nil{
		if err == sql.ErrNoRows{
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update book"})
		return
	}

	input.ID = id
	c.JSON(http.StatusOK, input)
}

// DeleteBookHandler adalah fungsi untuk menangani permintaan menghapus buku berdasarkan ID
// @Summary Delete a book by ID
// @Description Delete a book by its ID from the database
// @Tags books
// @Accept json
// @Produce json
// @Param id path int true "Book ID"
// @Success 200 {object} gin.H{"message": "Book deleted successfully"}
// @Failure 400 {object} gin.H{"error": "Bad Request"}
// @Failure 404 {object} gin.H{"error": "Book not found"}	
// @Failure 500 {object} gin.H{"error": "Internal Server Error"}
// @Router /books/{id} [delete]
func (h *BookHandler) DeleteBookHandler(c *gin.Context){
	idStr := c.Param("id")

	id, err := strconv.Atoi(idStr)
	if err != nil{
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid book ID"})
	}

	err = h.repo.DeleteBook(id)
	if err != nil{
		if err == sql.ErrNoRows{
			c.JSON(http.StatusNotFound, gin.H{"error": "Book not found"})
			return
		}

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete book"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Book deleted successfully"})
}