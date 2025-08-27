package main

import (
	"time"
)

// User represents the user model
type User struct {
	ID        int       `json:"id" db:"id"`
	Email     string    `json:"email" db:"email"`
	Password  string    `json:"-" db:"password"` // Hidden from JSON response
	FirstName string    `json:"firstname" db:"firstname"`
	LastName  string    `json:"lastname" db:"lastname"`
	Phone     string    `json:"phone" db:"phone"`
	Birthday  time.Time `json:"birthday" db:"birthday"`
	CreatedAt time.Time `json:"created_at" db:"created_at"`
	UpdatedAt time.Time `json:"updated_at" db:"updated_at"`
}

// RegisterRequest represents the registration request payload
type RegisterRequest struct {
	Email     string `json:"email" binding:"required" example:"user@example.com"`
	Password  string `json:"password" binding:"required" example:"password123"`
	FirstName string `json:"firstname" binding:"required" example:"John"`
	LastName  string `json:"lastname" binding:"required" example:"Doe"`
	Phone     string `json:"phone" binding:"required" example:"0812345678"`
	Birthday  string `json:"birthday" binding:"required" example:"1990-01-01"`
}

// LoginRequest represents the login request payload
type LoginRequest struct {
	Email    string `json:"email" binding:"required" example:"user@example.com"`
	Password string `json:"password" binding:"required" example:"password123"`
}

// LoginResponse represents the login response
type LoginResponse struct {
	Token string `json:"token" example:"eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9..."`
	User  User   `json:"user"`
}

// APIResponse represents a generic API response
type APIResponse struct {
	Message string      `json:"message" example:"Success"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error" example:"Invalid request"`
}
