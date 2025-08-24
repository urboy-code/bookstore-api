package main

import (
	"bookstore-api/handler"
	"bookstore-api/repository"
	"fmt"
	"log"
	"database/sql"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
)

func main(){
	connStr := "postgres://postgres:postgres@localhost:5432/bookstore?sslmode=disable"

	db, err := sql.Open("pgx", connStr)
	if err != nil {
		log.Fatalf("Faild to connect to the database: %v", err)
	}
	defer db.Close()

	err = db.Ping()
	if err != nil{
		log.Fatalf("Failed to ping the database: %v", err)
	}

	fmt.Println("Successfully connected to the database!")

	bookRepo := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)

	router := gin.Default()
	router.POST("/books", bookHandler.CreateBookHandler)
	router.GET("/books", bookHandler.GetBooksHandler)
	router.GET("/books/:id", bookHandler.GetBookByIDHandler)

	fmt.Println("Starting server on port 8080...")
	router.Run(":8080")
}