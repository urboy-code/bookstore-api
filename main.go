   package main

import (
	"bookstore-api/handler"
	"bookstore-api/repository"
	"fmt"
	"log"
	"database/sql"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/jackc/pgx/v5/stdlib"

)

func main(){
	err := godotenv.Load()
	if err != nil{
		log.Println("Error loading .env file, proceeding with system environment variables")
	}
	connStr := os.Getenv("DATABASE_URL")
	if connStr == ""{
		log.Fatal("DATABASE_URL environment variable not set")
	}

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
	router.PUT("/books/:id", bookHandler.UpdateBookHandler)
	router.DELETE("/books/:id", bookHandler.DeleteBookHandler)

	fmt.Println("Starting server on port 8080...")
	router.Run(":8080")
}