# C4 Model - Level 4: Code Diagram

## Overview
Code diagram แสดงรายละเอียดของ classes, interfaces, และ relationships ใน codebase ระดับ implementation

## Core Domain Classes

```mermaid
classDiagram
    namespace Domain {
        class User {
            +int ID
            +string Email
            +string Password
            +string FirstName
            +string LastName
            +string Phone
            +time.Time Birthday
            +time.Time CreatedAt
            +time.Time UpdatedAt
            +NewUser(email, password, firstName, lastName, phone, birthday) User
            +GetFullName() string
            +IsValidForUpdate() error
        }
        
        class UserRepository {
            <<interface>>
            +Create(ctx, user) error
            +GetByID(ctx, id) User, error
            +GetByEmail(ctx, email) User, error
            +Update(ctx, user) error
            +Delete(ctx, id) error
            +Exists(ctx, email) bool, error
        }
        
        class AuthService {
            <<interface>>
            +HashPassword(password) string, error
            +ComparePassword(hashed, password) error
            +GenerateToken(userID, email) string, error
            +ValidateToken(token) TokenClaims, error
        }
        
        class UserService {
            <<interface>>
            +Register(ctx, email, password, firstName, lastName, phone, birthday) User, error
            +Login(ctx, email, password) string, User, error
            +GetUserProfile(ctx, userID) User, error
            +UpdateUser(ctx, userID, firstName, lastName, phone, birthday) User, error
        }
        
        class TokenClaims {
            +int UserID
            +string Email
            +jwt.RegisteredClaims
        }
        
        class DomainError {
            +string Code
            +string Message
            +Error() string
        }
    }
    
    namespace UseCase {
        class UserUseCase {
            -UserRepository userRepo
            -AuthService authService
            +NewUserUseCase(userRepo, authService) UserUseCase
            +Register(ctx, email, password, firstName, lastName, phone, birthday) User, error
            +Login(ctx, email, password) string, User, error
            +GetUserProfile(ctx, userID) User, error
            +UpdateUser(ctx, userID, firstName, lastName, phone, birthday) User, error
        }
    }
    
    namespace Infrastructure {
        class SQLiteUserRepository {
            -sql.DB db
            +NewSQLiteUserRepository(db) SQLiteUserRepository
            +Create(ctx, user) error
            +GetByID(ctx, id) User, error
            +GetByEmail(ctx, email) User, error
            +Update(ctx, user) error
            +Delete(ctx, id) error
            +Exists(ctx, email) bool, error
        }
        
        class JWTAuthService {
            -[]byte secret
            +NewJWTAuthService() JWTAuthService
            +HashPassword(password) string, error
            +ComparePassword(hashed, password) error
            +GenerateToken(userID, email) string, error
            +ValidateToken(token) TokenClaims, error
        }
        
        class DatabaseConfig {
            +string Driver
            +string DSN
        }
    }
    
    namespace Interfaces {
        class UserHandler {
            -UserService userService
            -UserMapper mapper
            +NewUserHandler(userService) UserHandler
            +Register(w, r)
            +Login(w, r)
            +GetProfile(w, r)
            +UpdateProfile(w, r)
        }
        
        class AuthMiddleware {
            -AuthService authService
            +NewAuthMiddleware(authService) AuthMiddleware
            +Middleware(next) http.Handler
        }
        
        class Router {
            -UserHandler userHandler
            -AuthMiddleware authMiddleware
            +NewRouter(userService, authService) Router
            +SetupRoutes() chi.Router
        }
        
        class UserMapper {
            +NewUserMapper() UserMapper
            +ToUserResponse(user) UserResponse
            +ToLoginResponse(token, user) LoginResponse
            +ParseCreateUserRequest(req) email, password, firstName, lastName, phone, birthday, error
            +ParseUpdateUserRequest(req) firstName, lastName, phone, birthday, error
        }
    }
    
    namespace DTOs {
        class CreateUserRequest {
            +string Email
            +string Password
            +string FirstName
            +string LastName
            +string Phone
            +string Birthday
        }
        
        class UpdateUserRequest {
            +string FirstName
            +string LastName
            +string Phone
            +string Birthday
        }
        
        class LoginRequest {
            +string Email
            +string Password
        }
        
        class UserResponse {
            +int ID
            +string Email
            +string FirstName
            +string LastName
            +string Phone
            +time.Time Birthday
            +time.Time CreatedAt
            +time.Time UpdatedAt
        }
        
        class LoginResponse {
            +string Token
            +UserResponse User
        }
        
        class APIResponse {
            +string Message
            +interface{} Data
        }
        
        class ErrorResponse {
            +string Error
        }
    }
    
    namespace Application {
        class Container {
            +Config config
            +Database sql.DB
            +UserRepo UserRepository
            +AuthService AuthService
            +UserService UserService
            +Router Router
            +NewContainer() Container, error
            +Close() error
        }
        
        class Config {
            +ServerConfig Server
            +DatabaseConfig Database
        }
        
        class ServerConfig {
            +string Host
            +string Port
        }
    }
    
    %% Relationships
    UserUseCase ..|> UserService : implements
    UserUseCase --> UserRepository : depends on
    UserUseCase --> AuthService : depends on
    UserUseCase --> User : creates/manipulates
    
    SQLiteUserRepository ..|> UserRepository : implements
    JWTAuthService ..|> AuthService : implements
    
    UserHandler --> UserService : uses
    UserHandler --> UserMapper : uses
    AuthMiddleware --> AuthService : uses
    Router --> UserHandler : routes to
    Router --> AuthMiddleware : applies
    
    UserMapper --> CreateUserRequest : parses
    UserMapper --> UpdateUserRequest : parses
    UserMapper --> UserResponse : creates
    UserMapper --> LoginResponse : creates
    UserMapper --> User : converts from
    
    Container --> Config : loads
    Container --> SQLiteUserRepository : creates
    Container --> JWTAuthService : creates
    Container --> UserUseCase : creates
    Container --> Router : creates
    
    User --> DomainError : may return
    AuthService --> TokenClaims : returns
```

## Implementation Details

### 🌟 Domain Layer Implementation

#### User Entity
```go
package domain

import (
    "context"
    "time"
)

// User represents the core user entity in the domain
type User struct {
    ID        int
    Email     string
    Password  string
    FirstName string
    LastName  string
    Phone     string
    Birthday  time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
}

// NewUser creates a new user entity with validation
func NewUser(email, password, firstName, lastName, phone string, birthday time.Time) (*User, error) {
    if email == "" {
        return nil, ErrInvalidEmail
    }
    if firstName == "" {
        return nil, ErrInvalidFirstName
    }
    if lastName == "" {
        return nil, ErrInvalidLastName
    }

    now := time.Now()
    return &User{
        Email:     email,
        Password:  password,
        FirstName: firstName,
        LastName:  lastName,
        Phone:     phone,
        Birthday:  birthday,
        CreatedAt: now,
        UpdatedAt: now,
    }, nil
}

// GetFullName returns the user's full name
func (u *User) GetFullName() string {
    return u.FirstName + " " + u.LastName
}

// IsValidForUpdate validates user data for updates
func (u *User) IsValidForUpdate() error {
    if u.Email == "" {
        return ErrInvalidEmail
    }
    if u.FirstName == "" {
        return ErrInvalidFirstName
    }
    if u.LastName == "" {
        return ErrInvalidLastName
    }
    return nil
}
```

#### Domain Interfaces
```go
// UserRepository defines the contract for user data access
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int) (*User, error)
    GetByEmail(ctx context.Context, email string) (*User, error)
    Update(ctx context.Context, user *User) error
    Delete(ctx context.Context, id int) error
    Exists(ctx context.Context, email string) (bool, error)
}

// AuthService defines the contract for authentication operations
type AuthService interface {
    HashPassword(password string) (string, error)
    ComparePassword(hashedPassword, password string) error
    GenerateToken(userID int, email string) (string, error)
    ValidateToken(token string) (*TokenClaims, error)
}

// UserService defines the contract for user business operations
type UserService interface {
    Register(ctx context.Context, email, password, firstName, lastName, phone string, birthday time.Time) (*User, error)
    Login(ctx context.Context, email, password string) (string, *User, error)
    GetUserProfile(ctx context.Context, userID int) (*User, error)
    UpdateUser(ctx context.Context, userID int, firstName, lastName, phone string, birthday *time.Time) (*User, error)
}
```

#### Domain Errors
```go
// DomainError represents business logic errors
type DomainError struct {
    Code    string
    Message string
}

func (e DomainError) Error() string {
    return e.Message
}

// Domain-specific errors
var (
    ErrUserNotFound       = DomainError{Code: "USER_NOT_FOUND", Message: "User not found"}
    ErrUserAlreadyExists  = DomainError{Code: "USER_ALREADY_EXISTS", Message: "User already exists"}
    ErrInvalidCredentials = DomainError{Code: "INVALID_CREDENTIALS", Message: "Invalid credentials"}
    ErrInvalidToken       = DomainError{Code: "INVALID_TOKEN", Message: "Invalid token"}
    ErrInvalidEmail       = DomainError{Code: "INVALID_EMAIL", Message: "Invalid email"}
    ErrInvalidFirstName   = DomainError{Code: "INVALID_FIRST_NAME", Message: "Invalid first name"}
    ErrInvalidLastName    = DomainError{Code: "INVALID_LAST_NAME", Message: "Invalid last name"}
    // ... other domain errors
)
```

### 🔄 Use Case Implementation

#### User Use Case
```go
package usecase

import (
    "context"
    "time"
    "hello-world/internal/domain"
)

type UserUseCase struct {
    userRepo    domain.UserRepository
    authService domain.AuthService
}

func NewUserUseCase(userRepo domain.UserRepository, authService domain.AuthService) *UserUseCase {
    return &UserUseCase{
        userRepo:    userRepo,
        authService: authService,
    }
}

func (uc *UserUseCase) Register(ctx context.Context, email, password, firstName, lastName, phone string, birthday time.Time) (*domain.User, error) {
    // Check if user already exists
    exists, err := uc.userRepo.Exists(ctx, email)
    if err != nil {
        return nil, err
    }
    if exists {
        return nil, domain.ErrUserAlreadyExists
    }

    // Hash password
    hashedPassword, err := uc.authService.HashPassword(password)
    if err != nil {
        return nil, domain.ErrPasswordHashError
    }

    // Create user entity
    user, err := domain.NewUser(email, hashedPassword, firstName, lastName, phone, birthday)
    if err != nil {
        return nil, err
    }

    // Save user
    if err := uc.userRepo.Create(ctx, user); err != nil {
        return nil, domain.ErrUserCreationError
    }

    // Clear password for response
    userResponse := *user
    userResponse.Password = ""
    return &userResponse, nil
}

func (uc *UserUseCase) Login(ctx context.Context, email, password string) (string, *domain.User, error) {
    // Get user by email
    user, err := uc.userRepo.GetByEmail(ctx, email)
    if err != nil {
        return "", nil, domain.ErrInvalidCredentials
    }

    // Compare password
    if err := uc.authService.ComparePassword(user.Password, password); err != nil {
        return "", nil, domain.ErrInvalidCredentials
    }

    // Generate token
    token, err := uc.authService.GenerateToken(user.ID, user.Email)
    if err != nil {
        return "", nil, domain.ErrTokenGenerationError
    }

    // Clear password for response
    userResponse := *user
    userResponse.Password = ""
    return token, &userResponse, nil
}
```

### 🏗️ Infrastructure Implementation

#### SQLite Repository
```go
package infrastructure

import (
    "context"
    "database/sql"
    "strings"
    "hello-world/internal/domain"
)

type SQLiteUserRepository struct {
    db *sql.DB
}

func NewSQLiteUserRepository(db *sql.DB) *SQLiteUserRepository {
    return &SQLiteUserRepository{db: db}
}

func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error {
    query := `
        INSERT INTO users (email, password, firstname, lastname, phone, birthday, created_at, updated_at) 
        VALUES (?, ?, ?, ?, ?, ?, ?, ?)
    `
    
    result, err := r.db.ExecContext(ctx, query,
        user.Email, user.Password, user.FirstName, user.LastName,
        user.Phone, user.Birthday, user.CreatedAt, user.UpdatedAt,
    )
    
    if err != nil {
        if strings.Contains(err.Error(), "UNIQUE constraint failed") {
            return domain.ErrUserAlreadyExists
        }
        return err
    }
    
    id, err := result.LastInsertId()
    if err != nil {
        return err
    }
    
    user.ID = int(id)
    return nil
}

func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error) {
    query := `
        SELECT id, email, password, firstname, lastname, phone, birthday, created_at, updated_at 
        FROM users WHERE email = ?
    `
    
    user := &domain.User{}
    err := r.db.QueryRowContext(ctx, query, email).Scan(
        &user.ID, &user.Email, &user.Password, &user.FirstName, &user.LastName,
        &user.Phone, &user.Birthday, &user.CreatedAt, &user.UpdatedAt,
    )
    
    if err != nil {
        if err == sql.ErrNoRows {
            return nil, domain.ErrUserNotFound
        }
        return nil, err
    }
    
    return user, nil
}
```

#### JWT Auth Service
```go
package infrastructure

import (
    "os"
    "time"
    "hello-world/internal/domain"
    "github.com/golang-jwt/jwt/v5"
    "golang.org/x/crypto/bcrypt"
)

type JWTAuthService struct {
    secret []byte
}

func NewJWTAuthService() *JWTAuthService {
    secret := []byte("default-secret-key")
    if envSecret := os.Getenv("JWT_SECRET"); envSecret != "" {
        secret = []byte(envSecret)
    }
    
    return &JWTAuthService{secret: secret}
}

func (s *JWTAuthService) HashPassword(password string) (string, error) {
    hashedBytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
    if err != nil {
        return "", domain.ErrPasswordHashError
    }
    return string(hashedBytes), nil
}

func (s *JWTAuthService) ComparePassword(hashedPassword, password string) error {
    err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
    if err != nil {
        return domain.ErrInvalidCredentials
    }
    return nil
}

func (s *JWTAuthService) GenerateToken(userID int, email string) (string, error) {
    claims := &domain.TokenClaims{
        UserID: userID,
        Email:  email,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(24 * time.Hour)),
            IssuedAt:  jwt.NewNumericDate(time.Now()),
        },
    }
    
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    tokenString, err := token.SignedString(s.secret)
    if err != nil {
        return "", domain.ErrTokenGenerationError
    }
    
    return tokenString, nil
}
```

### 🌐 Interface Layer Implementation

#### User Handler
```go
package interfaces

import (
    "encoding/json"
    "net/http"
    "strconv"
    "hello-world/internal/domain"
    "hello-world/internal/interfaces/dto"
    "hello-world/internal/interfaces/mapper"
)

type UserHandler struct {
    userService domain.UserService
    mapper      *mapper.UserMapper
}

func NewUserHandler(userService domain.UserService) *UserHandler {
    return &UserHandler{
        userService: userService,
        mapper:      mapper.NewUserMapper(),
    }
}

func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request) {
    var req dto.CreateUserRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    email, password, firstName, lastName, phone, birthday, err := h.mapper.ParseCreateUserRequest(req)
    if err != nil {
        h.writeDomainErrorResponse(w, err)
        return
    }
    
    user, err := h.userService.Register(r.Context(), email, password, firstName, lastName, phone, birthday)
    if err != nil {
        h.writeDomainErrorResponse(w, err)
        return
    }
    
    response := h.mapper.ToUserResponse(user)
    h.writeSuccessResponse(w, http.StatusCreated, "User registered successfully", response)
}

func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request) {
    var req dto.LoginRequest
    if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
        h.writeErrorResponse(w, http.StatusBadRequest, "Invalid request body")
        return
    }
    
    token, user, err := h.userService.Login(r.Context(), req.Email, req.Password)
    if err != nil {
        h.writeDomainErrorResponse(w, err)
        return
    }
    
    response := h.mapper.ToLoginResponse(token, user)
    h.writeSuccessResponse(w, http.StatusOK, "Login successful", response)
}
```

## Class Relationships Analysis

### 🔗 Dependency Relationships

1. **Use Case → Domain**: Use cases depend on domain interfaces (UserRepository, AuthService)
2. **Infrastructure → Domain**: Infrastructure implements domain interfaces
3. **Interfaces → Use Cases**: HTTP handlers call use case methods
4. **Application → All Layers**: DI container wires all dependencies

### 🎯 Interface Segregation

- **UserRepository**: Data access operations only
- **AuthService**: Authentication operations only  
- **UserService**: Business service operations only
- **Small, focused interfaces** enable easy testing และ implementation swapping

### 🔄 Data Flow Patterns

1. **Request Flow**: HTTP → Handler → Use Case → Domain → Infrastructure → Database
2. **Response Flow**: Database → Infrastructure → Domain → Use Case → Handler → HTTP
3. **Error Flow**: Any layer can return domain errors ที่ propagate up the stack

### 🧪 Testing Patterns

- **Mock Implementations**: Each interface can be easily mocked
- **Dependency Injection**: Test dependencies can be injected
- **Layer Isolation**: Each layer can be tested independently
- **Integration Testing**: Full stack testing through HTTP endpoints

## Code Quality Metrics

### 📊 Complexity Analysis
- **Cyclomatic Complexity**: Low (simple methods, clear control flow)
- **Coupling**: Loose coupling ผ่าน interfaces
- **Cohesion**: High cohesion within each component
- **Maintainability Index**: High (clean code, good structure)

### 🎯 SOLID Principles Compliance
- **Single Responsibility**: Each class has one reason to change
- **Open/Closed**: Open for extension, closed for modification
- **Liskov Substitution**: Implementations can substitute interfaces
- **Interface Segregation**: Small, focused interfaces
- **Dependency Inversion**: Depend on abstractions, not concretions

### 📈 Test Coverage
```go
// Example test coverage by layer
Domain Layer:        100% // Pure business logic
Use Case Layer:      89%  // Business workflows  
Infrastructure:      85%  // Database และ services
Interface Adapters:  75%  // HTTP handlers
Application:         83%  // DI container
```
