package repository

import (
	"bookstore-api/model"
	"database/sql"
	"errors"
)

var ErrBookNotFound = errors.New("book not found")

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

func (r *BookRepository) GetBooks() ([]model.Book, error){

	query := `SELECT * FROM books`
	rows, err := r.db.Query(query)
	if err != nil{
		return nil, err
	}
	defer rows.Close()

	// for zero values
	books := make([]model.Book, 0)

	for rows.Next(){
		var book model.Book

		if err := rows.Scan(&book.ID, &book.Title, &book.Author, &book.Description); err != nil{
			return nil, err
		}

		books = append(books, book)
	}
	return books, nil
}

func (r *BookRepository) GetBookByID(id int) (model.Book, error){
	var book model.Book

	query := `SELECT * FROM books WHERE id = $1`

	err := r.db.QueryRow(query, id).Scan(&book.ID, &book.Title, &book.Author, &book.Description)
	if(err != nil){
		if err == sql.ErrNoRows{
			return book, ErrBookNotFound
		}
		return book, err
	}

	return book, nil
}

func (r *BookRepository) UpdateBook(id int, book model.Book) error{
	query := `UPDATE books SET title = $1, author = $2, description = $3 WHERE id = $4`

	result, err := r.db.Exec(query, book.Title, book.Author, book.Description, id)

	if err != nil{
		return err
	}

	rowsAffected, err := result.RowsAffected()

	if err != nil{
		return err
	}

	if rowsAffected == 0{
		return ErrBookNotFound
	}

	return nil
}

func (r *BookRepository) DeleteBook(id int) error{
	query := `DELETE FROM books WHERE id = $1`

	result, err := r.db.Exec(query, id)
	if err != nil{
		return err
	}

	rowsAffected, err := result.RowsAffected()
	if err != nil{
		return err
	}

	if rowsAffected == 0{
		return ErrBookNotFound
	}

	return nil
}