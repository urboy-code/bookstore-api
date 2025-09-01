package model

type User struct {
	ID			 int64  `json:"id"`
	Name		 string `json:"name" binding:"required"` 
	Email		 string `json:"email" binding:"required,email"`
	Password	 string `json:"password,omitempty" binding:"required,min=6"`
	PasswordHash string `json:"-"`
}