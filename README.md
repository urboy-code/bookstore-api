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

| Method | Endpoint      | Description           |
|--------|---------------|-----------------------|
| `POST` | `/books`      | Create a new book     |
| `GET`  | `/books`      | Get a list of all books|
| `GET`  | `/books/:id`  | Get a single book by ID|
| `PUT`  | `/books/:id`  | Update a book by ID   |
| `DELETE`| `/books/:id` | Delete a book by ID    |

---

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