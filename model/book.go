package model

type Book struct {
	ID          int    `json:"id"`
	Title       string `json:"title" binding:"required"`
	Author      string `json:"author" binding:"required"`
	Description string `json:"description"`
}