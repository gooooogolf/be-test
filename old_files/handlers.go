package main

import (
	"encoding/json"
	"net/http"
	"strings"
	"time"

	"golang.org/x/crypto/bcrypt"
)

// @Summary Hello World
// @Description Get a hello world message
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} Response
// @Router / [get]
func HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := Response{Message: "Hello world"}
	json.NewEncoder(w).Encode(response)
}

// @Summary Register User
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body RegisterRequest true "User registration data"
// @Success 201 {object} APIResponse
// @Failure 400 {object} ErrorResponse
// @Failure 409 {object} ErrorResponse
// @Router /register [post]
func RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req RegisterRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" || req.FirstName == "" || req.LastName == "" || req.Phone == "" || req.Birthday == "" {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "All fields are required"})
		return
	}

	// Parse birthday
	birthday, err := time.Parse("2006-01-02", req.Birthday)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid birthday format (YYYY-MM-DD)"})
		return
	}

	// Hash password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to hash password"})
		return
	}

	// Insert user into database
	query := `INSERT INTO users (email, password, firstname, lastname, phone, birthday) VALUES (?, ?, ?, ?, ?, ?)`
	result, err := db.Exec(query, req.Email, string(hashedPassword), req.FirstName, req.LastName, req.Phone, birthday)
	if err != nil {
		if err.Error() == "UNIQUE constraint failed: users.email" ||
			strings.Contains(err.Error(), "UNIQUE constraint failed") {
			w.WriteHeader(http.StatusConflict)
			json.NewEncoder(w).Encode(ErrorResponse{Error: "Email already exists"})
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to create user"})
		return
	}

	userID, _ := result.LastInsertId()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(APIResponse{
		Message: "User registered successfully",
		Data:    map[string]interface{}{"user_id": userID},
	})
}

// @Summary Login User
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body LoginRequest true "Login credentials"
// @Success 200 {object} LoginResponse
// @Failure 400 {object} ErrorResponse
// @Failure 401 {object} ErrorResponse
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid request body"})
		return
	}

	// Get user from database
	var user User
	var hashedPassword string
	query := `SELECT id, email, password, firstname, lastname, phone, birthday, created_at, updated_at FROM users WHERE email = ?`
	err := db.QueryRow(query, req.Email).Scan(
		&user.ID, &user.Email, &hashedPassword, &user.FirstName, &user.LastName,
		&user.Phone, &user.Birthday, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Check password
	if err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)); err != nil {
		w.WriteHeader(http.StatusUnauthorized)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Invalid email or password"})
		return
	}

	// Generate JWT token
	token, err := GenerateToken(user.ID, user.Email)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "Failed to generate token"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(LoginResponse{
		Token: token,
		User:  user,
	})
}

// @Summary Get Current User
// @Description Get current user information from JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} User
// @Failure 401 {object} ErrorResponse
// @Failure 404 {object} ErrorResponse
// @Router /me [get]
func MeHandler(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("user_id").(int)

	var user User
	query := `SELECT id, email, firstname, lastname, phone, birthday, created_at, updated_at FROM users WHERE id = ?`
	err := db.QueryRow(query, userID).Scan(
		&user.ID, &user.Email, &user.FirstName, &user.LastName,
		&user.Phone, &user.Birthday, &user.CreatedAt, &user.UpdatedAt,
	)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		json.NewEncoder(w).Encode(ErrorResponse{Error: "User not found"})
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}
