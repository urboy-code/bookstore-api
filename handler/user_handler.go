package handler

import (
	"bookstore-api/model"
	"bookstore-api/repository"
	"net/http"
	"os"
	"time"

	"github.com/golang-jwt/jwt/v5"
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

// LoginUserHandler handles user login and JWT token generation
func (h *UserHandler) LoginUserHandler(c *gin.Context){
	var input struct{
		Email		string `json:"email" binding:"required,email"`
		Password	string `json:"password" binding:"required,min=6"`
	}

	if err := c.ShouldBindJSON(&input); err != nil{
		c.JSON(http.StatusBadRequest, model.AppError{
			Code: http.StatusBadRequest,
			Message: "Invalid input: " + err.Error(),
		})
		return
	}

	// Verify user credentials to repository
	user, err := h.repo.GetUserByEmail(input.Email)
	if err != nil{
		c.JSON(http.StatusUnauthorized, model.AppError{
			Code: http.StatusUnauthorized,
			Message: "Invalid email or password",
		})
		return
	}

	// Generate JWT TOKEN
	// Create "claims" which include data we want to include in the token
	claims := jwt.MapClaims{
		"user_id" : user.ID,
		"email" : user.Email,
		"exp" : time.Now().Add(time.Hour * 24).Unix(), // Token expires in 24 hours
	}

	// Create the token using the claims
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// Sign the token with a secret key
	tokenString, err := token.SignedString([]byte(os.Getenv("JWT_SECRET_KEY")))
	if err != nil{
		ErrorHandler(c, err)
		return
	}

	// Send the token in response
	c.JSON(http.StatusOK, gin.H{
		"token": tokenString,
	})
}
