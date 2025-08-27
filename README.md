# Hello World Go API

This is a simple Go web API built with the [Chi router](https://github.com/go-chi/chi) that returns a JSON "Hello world" message.

## Features

- Built with Go and Chi router
- Returns JSON response
- Includes middleware for logging, request ID, and recovery
- Based on the Chi hello-world example

## Installation

1. Clone or download this project
2. Install dependencies:
   ```bash
   go mod tidy
   ```

## Running the Application

1. Start the server:
   ```bash
   go run main.go
   ```

2. The server will start on port 3333

## API Endpoints

### GET /

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

## Dependencies

- [github.com/go-chi/chi](https://github.com/go-chi/chi) - Lightweight, idiomatic and composable router for building Go HTTP services
