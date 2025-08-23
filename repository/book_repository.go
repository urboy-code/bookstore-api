package repository

import (
	"bookstore-api/model"
	"database/sql"
)

type BookRepository struct {
	db *sql.DB 
}

func NewBookRepository(db *sql.DB) *BookRepository{
	return &BookRepository{db: db}
}

func (r *BookRepository) CreateBook(book model.Book) (int, error){
	var bookID int

	query := `INSERT INTO books (title, author, description) VALUES ($1, $2, $3) RETURNING id`

	err := r.db.QueryRow(query, book.Title, book.Author, book.Description).Scan(&bookID)

	if(err != nil){
		return 0, err
	}

	return bookID, nil
}