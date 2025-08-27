# 📚 Bookstore API

A simple RESTful API for a bookstore built with Go, Gin, and PostgreSQL. This project is part of my personal learning journey to master backend development.

---

## ✨ Tech Stack

- **Go (Golang)**: Language
- **Gin Gonic**: Web Framework
- **PostgreSQL**: SQL Database
- **Docker & Docker Compose**: Containerization
- **Swagger**: API Documentation

---

## 🚀 API Endpoints

This API provides full CRUD functionality for managing books
| Method | Endpoint      | Description           |
|--------|---------------|-----------------------|
| `POST` | `/books`      | Create a new book     |
| `GET`  | `/books`      | Get a list of all books|
| `GET`  | `/books/:id`  | Get a single book by ID|
| `PUT`  | `/books/:id`  | Update a book by ID   |
| `DELETE`| `/books/:id` | Delete a book by ID    |

### Endpoint Details

#### `POST /books`
Creates a new book in the database/
- **Request Body:**
    ```json
    {
        "title": "New Book Title",
        "author": "Author Name",
        "description": "A great description."
    }
**Success Response (201 Created)**

#### `GET /books`
Retrieves a list of all books
- **Body**
    ```json
        [
            {
                "id": 1,
                "title": "New Book Title",
                "author": "Author Name",
                "description": "A great description."
            },
            {
                "id": 2,
                "title": "Another Book",
                "author": "Another Author",
                "description": "Another description."
            }
        ]
**Success Response (200 OK)**

#### `GET /books/:id`
Retrieves a single book by its unique ID.
- **Body**
    ```json
    {
        "id": 1,
        "title": "New Book Title",
        "author": "Author Name",
        "description": "A great description."
    }
- **Success Response (200 OK)**

#### `PUT /books/:id`
Updates the details of an existing book.
- **Request Body**
    ```json
        {
            "title": "Updated Book Title",
            "author": "Updated Author Name",
            "description": "An updated description."
        }
- **Success Response (200 OK)**

#### `DELETE /books/:id`
Deletes a book from the database.
- **Body**
    ```json
    {
        "message": "Book deleted successfully"
    }
- **Success Response (200 OK)**

## ⚙️ Setup & Installation

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/urboy-code/bookstore-api.git
    cd bookstore-api
    ```

2.  **Run the database using Docker:**
    ```bash
    docker compose up -d
    ```

3.  **Install dependencies and run the server:**
    ```bash
    go mod tidy
    go run main.go
    ```
    The server will be running on `http://localhost:8080`.