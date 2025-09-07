package handler

import (
	"bookstore-api/model"
	"bookstore-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

type UserHandler struct {
	repo *repository.UserRepository
}

// RegisterUser handles user registration 
func NewUserHandler(repo *repository.UserRepository) *UserHandler{
	return &UserHandler{repo: repo}
}

// RegisterUserHandler handles user registration
func (h *UserHandler) RegisterUserHandler(c *gin.Context){
	var input model.User

	// Bind JSON and Validate input (name, email, password)
	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, model.AppError{
			Code: http.StatusBadRequest,
			Message: "Invalid input" + err.Error(),
		})
		return
	}

	// Call repository to create new user
	userID, err := h.repo.CreateUser(&input)
	if err != nil{
		if err == repository.ErrMailExists{
			c.JSON(http.StatusConflict, model.AppError{
				Code: http.StatusConflict,
				Message: "Email already exists",
			})
			return
		}

		// If any other error, use ErrorHandler
		ErrorHandler(c, err)
		return
	}

	// send success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user_id": userID,
	})
}

