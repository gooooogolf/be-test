# Go Chi API with Authentication

This is a comprehensive Go web API built with the [Chi router](https://github.com/go-chi/chi) that includes user authentication, JWT tokens, SQLite database, and Swagger documentation.

## Features

- ✅ Built with Go and Chi router
- ✅ User registration and authentication
- ✅ JWT token-based authentication
- ✅ SQLite database for data persistence
- ✅ Password hashing with bcrypt
- ✅ Swagger/OpenAPI documentation
- ✅ RESTful API design
- ✅ Middleware for logging, request ID, and recovery

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Application

1. Start the server:
   ```bash
   go run .
   ```

2. The server will start on port 3333
3. Swagger UI will be available at: http://localhost:3333/swagger/

## API Endpoints

### Public Endpoints

#### GET /
Returns a JSON greeting message.

**Response:**
```json
{
  "message": "Hello world"
}
```

**Example:**
```bash
curl http://localhost:3333/
```

#### POST /register
Register a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123",
  "firstname": "John",
  "lastname": "Doe",
  "phone": "0812345678",
  "birthday": "1990-01-01"
}
```

**Response (201 Created):**
```json
{
  "message": "User registered successfully",
  "data": {
    "user_id": 1
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:3333/register \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123",
    "firstname": "John",
    "lastname": "Doe",
    "phone": "0812345678",
    "birthday": "1990-01-01"
  }'
```

#### POST /login
Login with email and password to get JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response (200 OK):**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 1,
    "email": "user@example.com",
    "firstname": "John",
    "lastname": "Doe",
    "phone": "0812345678",
    "birthday": "1990-01-01T00:00:00Z",
    "created_at": "2025-08-27T14:00:00Z",
    "updated_at": "2025-08-27T14:00:00Z"
  }
}
```

**Example:**
```bash
curl -X POST http://localhost:3333/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

### Protected Endpoints (Require JWT Token)

#### GET /me
Get current user information from JWT token.

**Headers:**
```
Authorization: Bearer <your-jwt-token>
```

**Response (200 OK):**
```json
{
  "id": 1,
  "email": "user@example.com",
  "firstname": "John",
  "lastname": "Doe",
  "phone": "0812345678",
  "birthday": "1990-01-01T00:00:00Z",
  "created_at": "2025-08-27T14:00:00Z",
  "updated_at": "2025-08-27T14:00:00Z"
}
```

**Example:**
```bash
curl -X GET http://localhost:3333/me \
  -H "Authorization: Bearer <your-jwt-token>"
```

## Error Responses

All endpoints may return error responses in the following format:

```json
{
  "error": "Error message description"
}
```

Common HTTP status codes:
- `400 Bad Request` - Invalid request data
- `401 Unauthorized` - Missing or invalid authentication
- `404 Not Found` - Resource not found
- `409 Conflict` - Resource already exists (e.g., email already registered)
- `500 Internal Server Error` - Server error

## Database

The application uses SQLite database with the following schema:

**Users Table:**
- `id` (INTEGER PRIMARY KEY)
- `email` (TEXT UNIQUE NOT NULL)
- `password` (TEXT NOT NULL) - bcrypt hashed
- `firstname` (TEXT NOT NULL)
- `lastname` (TEXT NOT NULL)
- `phone` (TEXT NOT NULL)
- `birthday` (DATE NOT NULL)
- `created_at` (DATETIME)
- `updated_at` (DATETIME)

## Authentication

The API uses JWT (JSON Web Tokens) for authentication:

1. Register a new account with `/register`
2. Login with `/login` to get a JWT token
3. Include the token in the `Authorization` header for protected endpoints:
   ```
   Authorization: Bearer <your-jwt-token>
   ```

## Environment Variables

- `JWT_SECRET`: Secret key for JWT token signing (default: "your-secret-key")

## Swagger Documentation

Interactive API documentation is available at:
```
http://localhost:3333/swagger/
```

The Swagger UI provides:
- Complete API reference
- Interactive endpoint testing
- Request/response examples
- Authentication testing

## Dependencies

- [github.com/go-chi/chi](https://github.com/go-chi/chi) - HTTP router
- [github.com/golang-jwt/jwt/v5](https://github.com/golang-jwt/jwt) - JWT tokens
- [github.com/mattn/go-sqlite3](https://github.com/mattn/go-sqlite3) - SQLite driver
- [golang.org/x/crypto/bcrypt](https://golang.org/x/crypto/bcrypt) - Password hashing
- [github.com/swaggo/http-swagger](https://github.com/swaggo/http-swagger) - Swagger UI

## Development

To regenerate Swagger documentation after API changes:
```bash
swag init
```

## Testing

This project includes comprehensive testing scripts to validate all API functionality.

### Quick Testing
```bash
# Make scripts executable (first time only)
chmod +x *.sh

# Run basic API tests
./test_apis.sh

# Run advanced tests with edge cases
./advanced_test.sh

# Clean test data for fresh testing
./cleanup_test.sh
```

### Test Coverage
- ✅ All API endpoints (GET, POST)
- ✅ Authentication flows (register, login, protected routes)
- ✅ Error handling (400, 401, 404, 409)
- ✅ JWT token validation
- ✅ Input validation and edge cases
- ✅ Concurrent request handling
- ✅ Performance testing

For detailed testing documentation, see [TESTING.md](TESTING.md).
