package model

// AppError represents a structured error with a message and an HTTP status code.
type AppError struct{
	Code int		`json:"code"`
	Message string `json:"message"`
}

