# Code Architecture Documentation

## Overview
เอกสารนี้แสดง Clean Architecture implementation ของ Go Backend API โดยใช้ C4 model เพื่ออธิบาย code structure และ component relationships

## C4 Model - Code Level Architecture

### 🏗️ Clean Architecture Layers

```mermaid
C4Component
    title Component Diagram - Go Backend API Clean Architecture

    Container_Boundary(api, "Go Backend API") {
        Component(main, "main.go", "Go", "Application entry point และ dependency injection initialization")
        
        Container_Boundary(app, "Application Layer") {
            Component(container, "Container", "Go", "Dependency injection container สำหรับ wire dependencies")
        }
        
        Container_Boundary(interfaces, "Interface Adapters Layer") {
            Component(handlers, "User Handler", "Go", "HTTP request handlers และ response formatting")
            Component(router, "Router", "Chi", "HTTP routing และ middleware setup")
            Component(middleware, "Auth Middleware", "Go", "JWT authentication middleware")
            Component(dtos, "DTOs", "Go", "Request/Response data transfer objects")
            Component(mappers, "Mappers", "Go", "Convert ระหว่าง DTOs และ Domain entities")
        }
        
        Container_Boundary(usecases, "Use Cases Layer") {
            Component(user_usecase, "User Use Case", "Go", "Application business logic และ orchestration")
        }
        
        Container_Boundary(domain, "Domain Layer") {
            Component(entities, "User Entity", "Go", "Core business entities และ domain logic")
            Component(interfaces_domain, "Repository Interfaces", "Go", "Abstract interfaces สำหรับ external dependencies")
            Component(errors, "Domain Errors", "Go", "Business-specific error definitions")
        }
        
        Container_Boundary(infrastructure, "Infrastructure Layer") {
            Component(repository, "SQLite Repository", "Go + SQLite", "Data persistence implementation")
            Component(auth_service, "JWT Auth Service", "Go + JWT", "Authentication และ token management")
            Component(database, "Database Connection", "Go + SQLite", "Database connection และ schema management")
        }
    }
    
    ContainerDb(sqlite, "SQLite Database", "SQLite", "User data storage")
    
    Rel(main, container, "initializes")
    Rel(container, handlers, "creates")
    Rel(container, user_usecase, "creates")
    Rel(container, repository, "creates")
    Rel(container, auth_service, "creates")
    
    Rel(handlers, mappers, "uses")
    Rel(handlers, user_usecase, "calls")
    Rel(router, handlers, "routes to")
    Rel(router, middleware, "applies")
    Rel(middleware, auth_service, "validates with")
    
    Rel(user_usecase, entities, "creates/manipulates")
    Rel(user_usecase, interfaces_domain, "depends on")
    
    Rel(repository, interfaces_domain, "implements")
    Rel(auth_service, interfaces_domain, "implements")
    Rel(repository, database, "uses")
    Rel(database, sqlite, "connects to")
```

### 🔄 Dependency Flow Diagram

```mermaid
flowchart TD
    subgraph "External"
        Client[API Client]
        Swagger[Swagger UI]
    end
    
    subgraph "Application Layer"
        Main[main.go]
        Container[DI Container]
    end
    
    subgraph "Interface Adapters"
        Router[Chi Router]
        Handlers[User Handlers]
        Middleware[Auth Middleware]
        DTOs[Request/Response DTOs]
        Mappers[Entity/DTO Mappers]
    end
    
    subgraph "Use Cases"
        UserUseCase[User Use Case]
    end
    
    subgraph "Domain Core"
        UserEntity[User Entity]
        Interfaces[Repository & Service Interfaces]
        DomainErrors[Domain Errors]
    end
    
    subgraph "Infrastructure"
        UserRepo[SQLite User Repository]
        AuthService[JWT Auth Service]
        Database[Database Connection]
    end
    
    subgraph "Data Store"
        SQLite[(SQLite Database)]
    end
    
    %% Client connections
    Client -.->|HTTP/REST| Router
    Swagger -.->|HTTP| Router
    
    %% Application flow
    Main --> Container
    Container --> Router
    Container --> Handlers
    Container --> UserUseCase
    Container --> UserRepo
    Container --> AuthService
    
    %% Interface layer
    Router --> Middleware
    Router --> Handlers
    Handlers --> Mappers
    Handlers --> UserUseCase
    Middleware --> AuthService
    
    %% Use case dependencies
    UserUseCase --> UserEntity
    UserUseCase --> Interfaces
    
    %% Infrastructure implementations
    UserRepo -.->|implements| Interfaces
    AuthService -.->|implements| Interfaces
    UserRepo --> Database
    Database --> SQLite
    
    %% Styling
    classDef domainStyle fill:#e1f5fe
    classDef usecaseStyle fill:#f3e5f5
    classDef interfaceStyle fill:#e8f5e8
    classDef infraStyle fill:#fff3e0
    classDef appStyle fill:#fce4ec
    
    class UserEntity,Interfaces,DomainErrors domainStyle
    class UserUseCase usecaseStyle
    class Router,Handlers,Middleware,DTOs,Mappers interfaceStyle
    class UserRepo,AuthService,Database infraStyle
    class Main,Container appStyle
```

## 📁 Directory Structure และ Components

### Clean Architecture Layers Mapping

```
📦 be-test/
├── 🚀 main.go                              # Application Entry Point
├── 📁 internal/
│   ├── 🎯 app/                             # Application Layer
│   │   ├── container.go                    # DI Container
│   │   └── container_test.go               # Container Tests
│   │
│   ├── 🌟 domain/                          # Domain Layer (Core)
│   │   ├── user.go                         # User Entity + Interfaces
│   │   └── user_test.go                    # Domain Tests
│   │
│   ├── 🔄 usecase/                         # Use Cases Layer
│   │   ├── user_usecase.go                 # Business Logic
│   │   └── user_usecase_test.go            # Use Case Tests
│   │
│   ├── 🌐 interfaces/                      # Interface Adapters Layer
│   │   ├── user_handler.go                 # HTTP Handlers
│   │   ├── auth_middleware.go              # Authentication
│   │   ├── router.go                       # Route Setup
│   │   ├── dto/
│   │   │   └── user_dto.go                 # Data Transfer Objects
│   │   ├── mapper/
│   │   │   ├── user_mapper.go              # Entity/DTO Conversion
│   │   │   └── user_mapper_test.go         # Mapper Tests
│   │   └── *_test.go                       # Interface Tests
│   │
│   └── 🏗️ infrastructure/                  # Infrastructure Layer
│       ├── database.go                     # DB Connection
│       ├── user_repository.go              # Data Persistence
│       ├── auth_service.go                 # JWT Service
│       └── *_test.go                       # Infrastructure Tests
│
├── 📊 pkg/                                 # Shared Packages
│   └── config/
│       ├── config.go                       # Configuration Management
│       └── config_test.go                  # Config Tests
│
├── 📚 docs/                                # Documentation
│   ├── database.md                         # Database Schema
│   ├── architecture.md                     # This file
│   ├── c4-model/                          # C4 Model Docs
│   └── swagger.*                          # API Documentation
│
└── 🧪 tests/                              # Integration Tests
```

## 🔧 Component Responsibilities

### 1. **Application Layer** (`internal/app/`)

**Container** (`container.go`):
```go
type Container struct {
    Config      *config.Config
    Database    *sql.DB
    UserRepo    domain.UserRepository
    AuthService domain.AuthService
    UserService domain.UserService
    Router      *interfaces.Router
}
```

**Responsibilities**:
- Dependency injection และ wiring
- Application lifecycle management
- Configuration loading
- Resource cleanup

### 2. **Domain Layer** (`internal/domain/`) - **Core Business Logic**

**User Entity**:
```go
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
```

**Interfaces**:
```go
type UserRepository interface {
    Create(ctx context.Context, user *User) error
    GetByID(ctx context.Context, id int) (*User, error)
    // ... other methods
}

type AuthService interface {
    HashPassword(password string) (string, error)
    ValidateToken(token string) (*TokenClaims, error)
    // ... other methods
}
```

**Responsibilities**:
- Core business entities
- Business rules และ validation
- Interface definitions สำหรับ external dependencies
- Domain-specific errors

### 3. **Use Cases Layer** (`internal/usecase/`) - **Application Business Logic**

**User Use Case**:
```go
type UserUseCase struct {
    userRepo    domain.UserRepository
    authService domain.AuthService
}

func (uc *UserUseCase) Register(ctx context.Context, email, password, firstName, lastName, phone string, birthday time.Time) (*domain.User, error)
func (uc *UserUseCase) Login(ctx context.Context, email, password string) (string, *domain.User, error)
func (uc *UserUseCase) GetUserProfile(ctx context.Context, userID int) (*domain.User, error)
func (uc *UserUseCase) UpdateUser(ctx context.Context, userID int, firstName, lastName, phone string, birthday *time.Time) (*domain.User, error)
```

**Responsibilities**:
- Orchestrate domain entities
- Implement application-specific business rules
- Coordinate between repositories และ services
- Handle transaction boundaries

### 4. **Interface Adapters Layer** (`internal/interfaces/`) - **External Communication**

**User Handler**:
```go
type UserHandler struct {
    userService domain.UserService
    mapper      *mapper.UserMapper
}

// HTTP Handlers
func (h *UserHandler) Register(w http.ResponseWriter, r *http.Request)
func (h *UserHandler) Login(w http.ResponseWriter, r *http.Request)
func (h *UserHandler) GetProfile(w http.ResponseWriter, r *http.Request)
func (h *UserHandler) UpdateProfile(w http.ResponseWriter, r *http.Request)
```

**DTOs**:
```go
type CreateUserRequest struct {
    Email     string `json:"email"`
    Password  string `json:"password"`
    FirstName string `json:"firstname"`
    LastName  string `json:"lastname"`
    Phone     string `json:"phone"`
    Birthday  string `json:"birthday"`
}

type UserResponse struct {
    ID        int       `json:"id"`
    Email     string    `json:"email"`
    FirstName string    `json:"firstname"`
    LastName  string    `json:"lastname"`
    Phone     string    `json:"phone"`
    Birthday  time.Time `json:"birthday"`
    CreatedAt time.Time `json:"created_at"`
    UpdatedAt time.Time `json:"updated_at"`
}
```

**Responsibilities**:
- HTTP request/response handling
- Data transformation (DTOs ↔ Entities)
- Input validation
- Authentication middleware
- Route configuration

### 5. **Infrastructure Layer** (`internal/infrastructure/`) - **External Concerns**

**SQLite Repository**:
```go
type SQLiteUserRepository struct {
    db *sql.DB
}

func (r *SQLiteUserRepository) Create(ctx context.Context, user *domain.User) error
func (r *SQLiteUserRepository) GetByEmail(ctx context.Context, email string) (*domain.User, error)
// ... other implementations
```

**JWT Auth Service**:
```go
type JWTAuthService struct {
    secret []byte
}

func (s *JWTAuthService) GenerateToken(userID int, email string) (string, error)
func (s *JWTAuthService) ValidateToken(token string) (*domain.TokenClaims, error)
// ... other implementations
```

**Responsibilities**:
- Database operations
- External service integrations
- Framework-specific implementations
- Infrastructure concerns (logging, monitoring)

## 🔀 Data Flow Examples

### User Registration Flow

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant Mapper
    participant UseCase
    participant Entity
    participant Repository
    participant Database

    Client->>Handler: POST /register {email, password, ...}
    Handler->>Mapper: ParseCreateUserRequest(dto)
    Mapper-->>Handler: email, password, firstName, ...
    Handler->>UseCase: Register(ctx, email, password, ...)
    UseCase->>Entity: NewUser(email, hashedPassword, ...)
    Entity-->>UseCase: user entity
    UseCase->>Repository: Create(ctx, user)
    Repository->>Database: INSERT INTO users...
    Database-->>Repository: user with ID
    Repository-->>UseCase: created user
    UseCase-->>Handler: created user
    Handler->>Mapper: ToUserResponse(user)
    Mapper-->>Handler: user response DTO
    Handler-->>Client: 201 Created {user data}
```

### User Authentication Flow

```mermaid
sequenceDiagram
    participant Client
    participant Handler
    participant UseCase
    participant Repository
    participant AuthService
    participant Database

    Client->>Handler: POST /login {email, password}
    Handler->>UseCase: Login(ctx, email, password)
    UseCase->>Repository: GetByEmail(ctx, email)
    Repository->>Database: SELECT * FROM users WHERE email=?
    Database-->>Repository: user record
    Repository-->>UseCase: user entity
    UseCase->>AuthService: ComparePassword(hashedPassword, password)
    AuthService-->>UseCase: password valid
    UseCase->>AuthService: GenerateToken(userID, email)
    AuthService-->>UseCase: JWT token
    UseCase-->>Handler: token, user
    Handler-->>Client: 200 OK {token, user}
```

## 🧪 Testing Architecture

### Test Coverage ตาม Layers

```mermaid
pie title Test Coverage by Layer
    "Domain Layer" : 100
    "Use Case Layer" : 90
    "Infrastructure Layer" : 85
    "Interface Layer" : 75
    "Application Layer" : 83
```

### Testing Strategy

**Unit Tests**:
- **Domain**: Entity business logic validation
- **Use Cases**: Business workflow testing with mocks
- **Infrastructure**: Repository และ service implementations
- **Interfaces**: HTTP handler behavior

**Integration Tests**:
- End-to-end API testing
- Database integration testing
- Authentication flow testing

## 🎯 Architecture Benefits

### ✅ **Clean Architecture Principles**

1. **Dependency Inversion**:
   - Use cases depend on domain interfaces
   - Infrastructure implements domain interfaces
   - Dependencies point inward

2. **Single Responsibility**:
   - Each layer has clear responsibility
   - Components are focused และ cohesive

3. **Open/Closed Principle**:
   - Easy to extend functionality
   - Closed for modification

4. **Interface Segregation**:
   - Small, focused interfaces
   - Clients depend only on what they use

### 🚀 **Development Benefits**

- **Testability**: Easy to unit test with mocks
- **Maintainability**: Clear separation of concerns
- **Flexibility**: Easy to swap implementations
- **Scalability**: Layer-based scaling strategies

## 🔮 Future Architecture Considerations

### Potential Enhancements

1. **Microservices Decomposition**:
   - Split into smaller services
   - Event-driven architecture

2. **CQRS Implementation**:
   - Separate read/write models
   - Event sourcing capabilities

3. **Observability**:
   - Structured logging
   - Metrics และ tracing
   - Health checks

4. **Caching Layer**:
   - Redis integration
   - Application-level caching
