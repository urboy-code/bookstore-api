package handler

import (
	"bookstore-api/model"
	"bookstore-api/repository"
	"net/http"

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