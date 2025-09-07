package main

import (
	"bookstore-api/handler"
	"bookstore-api/repository"
	"database/sql"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/jackc/pgx/v5/stdlib"
	"github.com/joho/godotenv"
)

func LoggerMiddleware() gin.HandlerFunc{
	return func (c *gin.Context){
		startTime := time.Now()
		c.Next()
		duration := time.Since(startTime)
		log.Printf("Request: %s %s - Status: %d - Duration:%s", c.Request.Method, c.Request.RequestURI, c.Writer.Status(), duration)
	}
}

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

	// For Books
	bookRepo := repository.NewBookRepository(db)
	bookHandler := handler.NewBookHandler(bookRepo)
	
	// For Users
	userRepo := repository.NewUserRepository(db)
	userHandler := handler.NewUserHandler(userRepo)

	router := gin.Default()
	router.Use(LoggerMiddleware())

	// Book routes
	router.POST("/books", bookHandler.CreateBookHandler)
	router.GET("/books", bookHandler.GetBooksHandler)
	router.GET("/books/:id", bookHandler.GetBookByIDHandler)
	router.PUT("/books/:id", bookHandler.UpdateBookHandler)
	router.DELETE("/books/:id", bookHandler.DeleteBookHandler)

	// User routes
	router.POST("/register", userHandler.RegisterUserHandler)


	fmt.Println("Starting server on port 8080...")
	router.Run(":8080")
}