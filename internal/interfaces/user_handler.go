package interfaces

import (
	"encoding/json"
	"net/http"
	"strconv"

	"hello-world/internal/domain"
	"hello-world/internal/interfaces/dto"
	"hello-world/internal/interfaces/mapper"
)

// UserHandler handles HTTP requests for user operations
type UserHandler struct {
	userService domain.UserService
	mapper      *mapper.UserMapper
}

// NewUserHandler creates a new UserHandler
func NewUserHandler(userService domain.UserService) *UserHandler {
	return &UserHandler{
		userService: userService,
		mapper:      mapper.NewUserMapper(),
	}
}

// APIResponse represents a generic API response
type APIResponse struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

// ErrorResponse represents an error response
type ErrorResponse struct {
	Error string `json:"error"`
}

// @Summary Hello World
// @Description Get a hello world message
// @Tags general
// @Accept json
// @Produce json
// @Success 200 {object} dto.APIResponse
// @Router / [get]
func (h *UserHandler) HelloHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	response := dto.APIResponse{Message: "Hello world"}
	json.NewEncoder(w).Encode(response)
}

// @Summary Register User
// @Description Register a new user
// @Tags auth
// @Accept json
// @Produce json
// @Param user body dto.CreateUserRequest true "User registration data"
// @Success 201 {object} dto.APIResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 409 {object} dto.ErrorResponse
// @Router /register [post]
func (h *UserHandler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.CreateUserRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if err := h.validateRegisterRequest(req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, err.Error())
		return
	}

	// Parse request using mapper
	email, password, firstName, lastName, phone, birthday, err := h.mapper.ParseCreateUserRequest(req)
	if err != nil {
		if domainErr, ok := err.(domain.DomainError); ok {
			h.sendErrorResponseWithCode(w, http.StatusBadRequest, domainErr.Message, domainErr.Code)
		} else {
			h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request data")
		}
		return
	}

	user, err := h.userService.Register(r.Context(), email, password, firstName, lastName, phone, birthday)
	if err != nil {
		if domainErr, ok := err.(domain.DomainError); ok {
			switch domainErr.Code {
			case "USER_ALREADY_EXISTS":
				h.sendErrorResponseWithCode(w, http.StatusConflict, domainErr.Message, domainErr.Code)
			case "INVALID_BIRTHDAY":
				h.sendErrorResponseWithCode(w, http.StatusBadRequest, domainErr.Message, domainErr.Code)
			default:
				h.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			}
		} else {
			h.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
		}
		return
	}

	userResponse := h.mapper.ToUserResponse(user)
	h.sendSuccessResponse(w, http.StatusCreated, "User registered successfully", userResponse)
}

// @Summary Login User
// @Description Login with email and password
// @Tags auth
// @Accept json
// @Produce json
// @Param credentials body dto.LoginRequest true "Login credentials"
// @Success 200 {object} dto.LoginResponse
// @Failure 400 {object} dto.ErrorResponse
// @Failure 401 {object} dto.ErrorResponse
// @Router /login [post]
func (h *UserHandler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var req dto.LoginRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		h.sendErrorResponse(w, http.StatusBadRequest, "Invalid request body")
		return
	}

	// Validate required fields
	if req.Email == "" || req.Password == "" {
		h.sendErrorResponse(w, http.StatusBadRequest, "Email and password are required")
		return
	}

	token, user, err := h.userService.Login(r.Context(), req.Email, req.Password)
	if err != nil {
		if domainErr, ok := err.(domain.DomainError); ok {
			switch domainErr.Code {
			case "INVALID_CREDENTIALS":
				h.sendErrorResponseWithCode(w, http.StatusUnauthorized, domainErr.Message, domainErr.Code)
			default:
				h.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "Login failed")
		return
	}

	loginResponse := h.mapper.ToLoginResponse(token, user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(loginResponse)
}

// @Summary Get Current User
// @Description Get current user information from JWT token
// @Tags auth
// @Accept json
// @Produce json
// @Security ApiKeyAuth
// @Success 200 {object} dto.UserResponse
// @Failure 401 {object} dto.ErrorResponse
// @Failure 404 {object} dto.ErrorResponse
// @Router /me [get]
func (h *UserHandler) MeHandler(w http.ResponseWriter, r *http.Request) {
	userID, err := h.getUserIDFromContext(r)
	if err != nil {
		h.sendErrorResponse(w, http.StatusUnauthorized, "Invalid user context")
		return
	}

	user, err := h.userService.GetUserProfile(r.Context(), userID)
	if err != nil {
		if domainErr, ok := err.(domain.DomainError); ok {
			switch domainErr.Code {
			case "USER_NOT_FOUND":
				h.sendErrorResponseWithCode(w, http.StatusNotFound, domainErr.Message, domainErr.Code)
			default:
				h.sendErrorResponse(w, http.StatusInternalServerError, "Internal server error")
			}
			return
		}
		h.sendErrorResponse(w, http.StatusInternalServerError, "Failed to get user profile")
		return
	}

	userResponse := h.mapper.ToUserResponse(user)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(userResponse)
}

// Helper methods

func (h *UserHandler) sendErrorResponse(w http.ResponseWriter, statusCode int, message string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message})
}

func (h *UserHandler) sendErrorResponseWithCode(w http.ResponseWriter, statusCode int, message, code string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.ErrorResponse{Error: message, Code: code})
}

func (h *UserHandler) sendSuccessResponse(w http.ResponseWriter, statusCode int, message string, data interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(statusCode)
	json.NewEncoder(w).Encode(dto.APIResponse{Message: message, Data: data})
}

func (h *UserHandler) validateRegisterRequest(req dto.CreateUserRequest) error {
	if req.Email == "" {
		return domain.DomainError{Code: "VALIDATION_ERROR", Message: "Email is required"}
	}
	if req.Password == "" {
		return domain.DomainError{Code: "VALIDATION_ERROR", Message: "Password is required"}
	}
	if req.FirstName == "" {
		return domain.DomainError{Code: "VALIDATION_ERROR", Message: "First name is required"}
	}
	if req.LastName == "" {
		return domain.DomainError{Code: "VALIDATION_ERROR", Message: "Last name is required"}
	}
	if req.Phone == "" {
		return domain.DomainError{Code: "VALIDATION_ERROR", Message: "Phone is required"}
	}
	if req.Birthday == "" {
		return domain.DomainError{Code: "VALIDATION_ERROR", Message: "Birthday is required"}
	}
	return nil
}

func (h *UserHandler) getUserIDFromContext(r *http.Request) (int, error) {
	userIDValue := r.Context().Value("user_id")
	if userIDValue == nil {
		return 0, domain.ErrUnauthorized
	}

	switch v := userIDValue.(type) {
	case int:
		return v, nil
	case string:
		return strconv.Atoi(v)
	default:
		return 0, domain.ErrUnauthorized
	}
}
