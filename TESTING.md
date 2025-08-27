# API Testing Scripts

This directory contains comprehensive testing scripts for the Go Chi API backend.

## Available Scripts

### 1. `test_apis.sh` - Basic API Testing
A comprehensive test script that validates all API endpoints with various scenarios.

**Features:**
- ✅ Tests all CRUD operations
- ✅ Validates authentication flows
- ✅ Tests error handling
- ✅ JWT token validation
- ✅ Colored output for easy reading
- ✅ Detailed test results

**Usage:**
```bash
./test_apis.sh
```

**Test Coverage:**
- GET / (Hello World)
- POST /register (Success, Duplicate, Invalid data)
- POST /login (Valid, Invalid credentials, Non-existent user)
- GET /me (Valid token, Invalid token, No authorization)
- GET /swagger/ (Documentation)
- Invalid endpoints (404 testing)

### 2. `advanced_test.sh` - Advanced Testing Suite
Extended testing with edge cases, performance testing, and concurrent requests.

**Features:**
- 🧪 Input validation testing
- 🔐 Authentication & authorization edge cases
- 🌐 HTTP method validation
- 📝 Content-Type validation
- 📊 Large payload testing
- ⚡ Concurrent request testing
- 📈 Performance metrics
- 📊 Detailed statistics

**Usage:**
```bash
./advanced_test.sh
```

**Test Suites:**
1. **Input Validation** - Empty fields, invalid formats, edge cases
2. **Authentication & Authorization** - JWT handling, malformed tokens
3. **HTTP Methods** - Method not allowed scenarios
4. **Content Type Validation** - Missing/wrong headers
5. **Large Payload** - Stress testing with large data
6. **Concurrent Requests** - Multiple simultaneous requests
7. **Performance Testing** - Response time measurements

### 3. `cleanup_test.sh` - Test Data Cleanup
Utility script to reset the testing environment.

**Features:**
- 🧹 Removes SQLite database
- 🗑️ Cleans temporary files
- 🔄 Prepares for fresh testing

**Usage:**
```bash
./cleanup_test.sh
```

## Prerequisites

### Required Tools
- `curl` - For making HTTP requests
- `jq` - For JSON parsing (optional but recommended)
- `bash` - Shell environment

### Server Requirements
- Go API server running on `http://localhost:3333`
- SQLite database support
- All API endpoints implemented

## Quick Start

1. **Start the API server:**
   ```bash
   go run .
   ```

2. **Run basic tests:**
   ```bash
   ./test_apis.sh
   ```

3. **For comprehensive testing:**
   ```bash
   ./advanced_test.sh
   ```

4. **To reset test data:**
   ```bash
   ./cleanup_test.sh
   ```

## Test Output

### Color Coding
- 🔵 **Blue** - Test information and headers
- 🟢 **Green** - Successful tests and positive results
- 🔴 **Red** - Failed tests and errors
- 🟡 **Yellow** - Warnings and informational messages
- 🟣 **Purple** - Test suite headers (advanced testing)

### Status Codes
- ✅ **PASSED** - Test executed successfully with expected result
- ❌ **FAILED** - Test failed or returned unexpected result

## Example Output

```bash
🚀 Starting API Tests for Go Chi Backend
Base URL: http://localhost:3333

📋 Test 1: Hello World Endpoint
=== GET / ===
Expected Status: 200
Actual Status: 200
✅ PASSED
Response: {"message":"Hello world"}
```

## Troubleshooting

### Common Issues

1. **Server not running:**
   ```
   ❌ Server is not running at http://localhost:3333
   ```
   **Solution:** Start the server with `go run .`

2. **Permission denied:**
   ```
   bash: ./test_apis.sh: Permission denied
   ```
   **Solution:** Make script executable with `chmod +x test_apis.sh`

3. **jq not found:**
   ```
   ⚠️ jq is not installed. Some tests may have limited output parsing.
   ```
   **Solution:** Install jq with `brew install jq` (macOS) or your system's package manager

### Database Issues

If tests are failing due to existing data:
1. Run `./cleanup_test.sh` to reset the database
2. Restart the server
3. Run tests again

## Test Scenarios Covered

### Authentication Flow
1. User registration with valid data
2. Duplicate email handling
3. Invalid input validation
4. User login with correct credentials
5. Login with wrong password
6. Login with non-existent user
7. JWT token generation and validation
8. Protected endpoint access
9. Invalid token handling
10. Missing authorization header

### Error Handling
- 400 Bad Request (invalid input)
- 401 Unauthorized (auth failures)
- 404 Not Found (invalid endpoints)
- 409 Conflict (duplicate resources)
- 405 Method Not Allowed (wrong HTTP methods)

### Edge Cases
- Empty fields
- Large payloads
- Malformed JSON
- Wrong content types
- Concurrent requests
- Performance under load

## Contributing

To add new tests:

1. **Basic tests:** Add to `test_apis.sh`
2. **Advanced tests:** Add to `advanced_test.sh`
3. **Follow the existing pattern:**
   ```bash
   run_test "Test Name" "curl command" "expected_status" "description"
   ```

## Notes

- Tests are designed to be non-destructive when possible
- Some tests may create test data (users, etc.)
- Use `cleanup_test.sh` between test runs for consistent results
- Advanced tests include performance metrics and may take longer to complete
