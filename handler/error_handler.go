package handler

import (
	"bookstore-api/model"
	"bookstore-api/repository"
	"net/http"

	"github.com/gin-gonic/gin"
)

// ErrorHandler will handle errors and send appropriate HTTP responses
func ErrorHandler(c *gin.Context, err error){
	var appErr model.AppError

	switch err{
	case repository.ErrBookNotFound:
			appErr = model.AppError{
				Code: http.StatusNotFound,
				Message: err.Error(),
			}
		default:
			appErr = model.AppError{
				Code: http.StatusInternalServerError,
				Message: "Internal Server Error",
			}
	}
	c.JSON(appErr.Code, appErr)
}