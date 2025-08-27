#!/bin/bash

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# API Base URL
BASE_URL="http://localhost:3333"

# Test data
TEST_EMAIL="test@example.com"
TEST_PASSWORD="password123"
TEST_FIRSTNAME="John"
TEST_LASTNAME="Doe"
TEST_PHONE="0812345678"
TEST_BIRTHDAY="1990-01-01"

# JWT Token variable
JWT_TOKEN=""

# Function to print test results
print_test_result() {
    local test_name="$1"
    local status_code="$2"
    local expected_code="$3"
    local response="$4"
    
    echo -e "\n${BLUE}=== $test_name ===${NC}"
    echo -e "Expected Status: $expected_code"
    echo -e "Actual Status: $status_code"
    
    if [ "$status_code" = "$expected_code" ]; then
        echo -e "${GREEN}✅ PASSED${NC}"
    else
        echo -e "${RED}❌ FAILED${NC}"
    fi
    
    echo -e "Response: ${YELLOW}$response${NC}"
}

# Function to extract JWT token from login response
extract_jwt_token() {
    local response="$1"
    JWT_TOKEN=$(echo "$response" | jq -r '.token' 2>/dev/null)
    if [ "$JWT_TOKEN" = "null" ] || [ -z "$JWT_TOKEN" ]; then
        JWT_TOKEN=""
    fi
}

echo -e "${BLUE}🚀 Starting API Tests for Go Chi Backend${NC}"
echo -e "${BLUE}Base URL: $BASE_URL${NC}\n"

# Check if server is running
echo -e "${YELLOW}Checking if server is running...${NC}"
if ! curl -s "$BASE_URL" > /dev/null; then
    echo -e "${RED}❌ Server is not running at $BASE_URL${NC}"
    echo -e "${YELLOW}Please start the server with: go run .${NC}"
    exit 1
fi
echo -e "${GREEN}✅ Server is running${NC}\n"

# Test 1: GET / - Hello World
echo -e "${BLUE}📋 Test 1: Hello World Endpoint${NC}"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json "$BASE_URL/")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "GET /" "$status_code" "200" "$response_body"

# Test 2: POST /register - User Registration (Success)
echo -e "\n${BLUE}📋 Test 2: User Registration (Success)${NC}"
registration_data='{
    "email": "'$TEST_EMAIL'",
    "password": "'$TEST_PASSWORD'",
    "firstname": "'$TEST_FIRSTNAME'",
    "lastname": "'$TEST_LASTNAME'",
    "phone": "'$TEST_PHONE'",
    "birthday": "'$TEST_BIRTHDAY'"
}'

response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d "$registration_data")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "POST /register (Success)" "$status_code" "201" "$response_body"

# Test 3: POST /register - Duplicate Email (Should fail)
echo -e "\n${BLUE}📋 Test 3: User Registration (Duplicate Email)${NC}"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d "$registration_data")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "POST /register (Duplicate)" "$status_code" "409" "$response_body"

# Test 4: POST /register - Invalid Data (Missing fields)
echo -e "\n${BLUE}📋 Test 4: User Registration (Invalid Data)${NC}"
invalid_data='{"email": "invalid@test.com"}'
response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d "$invalid_data")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "POST /register (Invalid Data)" "$status_code" "400" "$response_body"

# Test 5: POST /login - Valid Credentials
echo -e "\n${BLUE}📋 Test 5: User Login (Valid Credentials)${NC}"
login_data='{
    "email": "'$TEST_EMAIL'",
    "password": "'$TEST_PASSWORD'"
}'

response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "$login_data")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "POST /login (Valid)" "$status_code" "200" "$response_body"

# Extract JWT token for subsequent tests
if [ "$status_code" = "200" ]; then
    extract_jwt_token "$response_body"
    if [ -n "$JWT_TOKEN" ]; then
        echo -e "${GREEN}✅ JWT Token extracted successfully${NC}"
        echo -e "${YELLOW}Token: ${JWT_TOKEN:0:50}...${NC}"
    else
        echo -e "${RED}❌ Failed to extract JWT token${NC}"
    fi
fi

# Test 6: POST /login - Invalid Credentials
echo -e "\n${BLUE}📋 Test 6: User Login (Invalid Credentials)${NC}"
invalid_login_data='{
    "email": "'$TEST_EMAIL'",
    "password": "wrongpassword"
}'

response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "$invalid_login_data")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "POST /login (Invalid)" "$status_code" "401" "$response_body"

# Test 7: POST /login - Non-existent User
echo -e "\n${BLUE}📋 Test 7: User Login (Non-existent User)${NC}"
nonexistent_login_data='{
    "email": "nonexistent@example.com",
    "password": "password123"
}'

response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "$nonexistent_login_data")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "POST /login (Non-existent)" "$status_code" "401" "$response_body"

# Test 8: GET /me - Valid JWT Token
if [ -n "$JWT_TOKEN" ]; then
    echo -e "\n${BLUE}📋 Test 8: Get User Profile (Valid Token)${NC}"
    response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
        -X GET "$BASE_URL/me" \
        -H "Authorization: Bearer $JWT_TOKEN")
    status_code=$(echo "$response" | tail -c 4)
    response_body=$(cat /tmp/response.json)
    print_test_result "GET /me (Valid Token)" "$status_code" "200" "$response_body"
else
    echo -e "\n${RED}❌ Skipping GET /me test - No JWT token available${NC}"
fi

# Test 9: GET /me - Invalid JWT Token
echo -e "\n${BLUE}📋 Test 9: Get User Profile (Invalid Token)${NC}"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X GET "$BASE_URL/me" \
    -H "Authorization: Bearer invalid_token_here")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "GET /me (Invalid Token)" "$status_code" "401" "$response_body"

# Test 10: GET /me - Missing Authorization Header
echo -e "\n${BLUE}📋 Test 10: Get User Profile (No Authorization)${NC}"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json \
    -X GET "$BASE_URL/me")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "GET /me (No Auth)" "$status_code" "401" "$response_body"

# Test 11: GET /swagger/ - Swagger Documentation
echo -e "\n${BLUE}📋 Test 11: Swagger Documentation${NC}"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json "$BASE_URL/swagger/")
status_code=$(echo "$response" | tail -c 4)
if [ "$status_code" = "200" ]; then
    response_body="Swagger UI loaded successfully"
else
    response_body=$(cat /tmp/response.json)
fi
print_test_result "GET /swagger/" "$status_code" "200" "$response_body"

# Test 12: Invalid Endpoint
echo -e "\n${BLUE}📋 Test 12: Invalid Endpoint${NC}"
response=$(curl -s -w "%{http_code}" -o /tmp/response.json "$BASE_URL/invalid-endpoint")
status_code=$(echo "$response" | tail -c 4)
response_body=$(cat /tmp/response.json)
print_test_result "GET /invalid-endpoint" "$status_code" "404" "$response_body"

# Summary
echo -e "\n${BLUE}📊 Test Summary${NC}"
echo -e "${BLUE}===============${NC}"
echo -e "${GREEN}✅ All critical API endpoints tested${NC}"
echo -e "${YELLOW}📝 Manual verification recommended for:${NC}"
echo -e "   - Swagger UI functionality at $BASE_URL/swagger/"
echo -e "   - Database persistence across server restarts"
echo -e "   - JWT token expiration handling"

# Cleanup
rm -f /tmp/response.json

echo -e "\n${GREEN}🎉 API Testing Complete!${NC}"
