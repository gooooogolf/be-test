# Clean Architecture Implementation

This project has been refactored to follow Clean Architecture principles as defined by Robert C. Martin (Uncle Bob).

## Architecture Overview

The application is structured in the following layers:

### 1. Domain Layer (`internal/domain/`)
- **Entities**: Core business objects (`User`)
- **Interfaces**: Contracts for external dependencies (`UserRepository`, `AuthService`, `UserService`)
- **Domain Errors**: Business rule violations and domain-specific errors
- **Business Rules**: Core business logic embedded in entities

### 2. Use Case Layer (`internal/usecase/`)
- **Application Services**: Orchestrates domain entities and external services
- **Business Logic**: Application-specific business rules
- **Implements Domain Interfaces**: Use cases implement the domain service interfaces

### 3. Interface Adapters Layer (`internal/interfaces/`)
- **Controllers/Handlers**: HTTP request handlers (`UserHandler`)
- **DTOs**: Data Transfer Objects for external communication
- **Mappers**: Convert between domain entities and DTOs
- **Presenters**: Format data for external consumers

### 4. Infrastructure Layer (`internal/infrastructure/`)
- **Database Repositories**: Data persistence implementations
- **External Services**: Third-party service integrations
- **Frameworks & Drivers**: Database drivers, HTTP frameworks, etc.

### 5. Application Layer (`internal/app/`)
- **Dependency Injection**: Wiring of all dependencies
- **Application Configuration**: App-level setup and teardown

## Key Principles Implemented

### 1. Dependency Inversion
- High-level modules (use cases) don't depend on low-level modules (infrastructure)
- Both depend on abstractions (interfaces in domain layer)
- Dependencies point inward toward the domain

### 2. Interface Segregation
- Small, focused interfaces
- Clients depend only on interfaces they use

### 3. Single Responsibility
- Each layer has a clear responsibility
- Entities focus on business rules
- Use cases focus on application workflows
- Handlers focus on HTTP concerns

### 4. Open/Closed Principle
- Easy to extend with new features
- Closed for modification of existing code

## Dependency Flow

```
main.go
    ↓
app/container.go (DI Container)
    ↓
interfaces/ (Controllers, DTOs, Mappers)
    ↓
usecase/ (Application Business Logic)
    ↓
domain/ (Entities, Interfaces, Domain Logic)
    ↑
infrastructure/ (Repositories, External Services)
```

## Benefits of This Architecture

1. **Testability**: Easy to unit test business logic in isolation
2. **Maintainability**: Clear separation of concerns
3. **Flexibility**: Easy to swap implementations (e.g., database)
4. **Independence**: Business logic independent of frameworks
5. **Scalability**: Easy to add new features without affecting existing code

## File Structure

```
internal/
├── app/
│   └── container.go          # Dependency injection container
├── domain/
│   └── user.go              # User entity, interfaces, domain errors
├── usecase/
│   └── user_usecase.go      # Application business logic
├── interfaces/
│   ├── dto/
│   │   └── user_dto.go      # Data transfer objects
│   ├── mapper/
│   │   └── user_mapper.go   # Entity/DTO conversion
│   ├── user_handler.go      # HTTP handlers
│   ├── auth_middleware.go   # Authentication middleware
│   └── router.go           # Route configuration
└── infrastructure/
    ├── database.go          # Database connection
    ├── user_repository.go   # User data persistence
    └── auth_service.go      # JWT authentication service
```

## Testing Strategy

- **Unit Tests**: Test domain entities and use cases in isolation
- **Integration Tests**: Test repository implementations with real database
- **Contract Tests**: Test that implementations satisfy interfaces
- **End-to-End Tests**: Test complete workflows through HTTP API

## Adding New Features

1. **Define Domain**: Add entities and interfaces to `domain/`
2. **Implement Use Cases**: Add business logic to `usecase/`
3. **Create Interfaces**: Add DTOs, mappers, and handlers to `interfaces/`
4. **Implement Infrastructure**: Add repositories/services to `infrastructure/`
5. **Wire Dependencies**: Update `app/container.go`
