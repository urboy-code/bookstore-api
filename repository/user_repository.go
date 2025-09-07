package repository

import (
	"database/sql"
	"bookstore-api/model"
	"errors"
	"strings"

	"golang.org/x/crypto/bcrypt"
)

// define custom errors for user
var ErrMailExists = errors.New("email already exists")

type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository{
	return &UserRepository{db: db}
}

// CreateUser for hashing password and storing user in db
func (r *UserRepository) CreateUser(user *model.User) (int, error){
	// hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	if err != nil{
		return 0, err
	}

	var userID int

	// Save the user to the database with hashed password
	query := `INSERT INTO users (name, email, password_hash) values ($1, $2, $3) RETURNING id`
	err = r.db.QueryRow(query, user.Name, user.Email, string(hashedPassword)).Scan(&userID)
	if err != nil {
		// Check error if any existing email (violates unique constraint)
		if strings.Contains(err.Error(), "unique constraint"){
			return 0, ErrMailExists
		}
		return 0, err
	}
	return userID, nil
}

// GetUserByEmail fetch user by email
func (r *UserRepository) GetUserByEmail(email string) (model.User, error){
	var user model.User
	query := `SELECT id, name, email, password_hash FROM users WHERE email=$1`

	err := r.db.QueryRow(query, email).Scan(&user.ID, &user.Name, &user.Email, &user.PasswordHash)
	if err != nil{
		if err == sql.ErrNoRows{
			return user, ErrBookNotFound
		}
		return user, err
	}
	return user, nil
}