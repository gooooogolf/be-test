#!/bin/bash

# Advanced API Testing Script
# This script includes comprehensive testing scenarios

# Colors for output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
PURPLE='\033[0;35m'
NC='\033[0m' # No Color

# Configuration
BASE_URL="http://localhost:3333"
TEST_EMAIL="advanced.test@example.com"
TEST_PASSWORD="securePassword123!"
CONCURRENT_REQUESTS=5

# Counters
PASSED=0
FAILED=0
TOTAL=0

# Function to increment test counters
count_test() {
    local status="$1"
    TOTAL=$((TOTAL + 1))
    if [ "$status" = "PASSED" ]; then
        PASSED=$((PASSED + 1))
    else
        FAILED=$((FAILED + 1))
    fi
}

# Function to run test with detailed output
run_advanced_test() {
    local test_name="$1"
    local curl_command="$2"
    local expected_status="$3"
    local description="$4"
    
    echo -e "\n${PURPLE}🧪 $test_name${NC}"
    echo -e "${BLUE}Description: $description${NC}"
    
    # Execute curl command and capture response
    eval "$curl_command" > /tmp/test_output.json 2>/dev/null
    local exit_code=$?
    local actual_status=$(tail -n1 /tmp/test_output.json)
    local response_body=$(head -n -1 /tmp/test_output.json)
    
    echo -e "Expected: $expected_status | Actual: $actual_status"
    
    if [ "$actual_status" = "$expected_status" ] && [ $exit_code -eq 0 ]; then
        echo -e "${GREEN}✅ PASSED${NC}"
        count_test "PASSED"
    else
        echo -e "${RED}❌ FAILED${NC}"
        count_test "FAILED"
    fi
    
    if [ -n "$response_body" ] && [ "$response_body" != "$actual_status" ]; then
        echo -e "${YELLOW}Response: $response_body${NC}"
    fi
}

echo -e "${BLUE}🚀 Advanced API Testing Suite${NC}"
echo -e "${BLUE}==============================${NC}\n"

# Check if jq is available for JSON parsing
if ! command -v jq &> /dev/null; then
    echo -e "${YELLOW}⚠️  jq is not installed. Some tests may have limited output parsing.${NC}\n"
fi

# Test Suite 1: Input Validation
echo -e "${PURPLE}📋 Test Suite 1: Input Validation${NC}"

run_advanced_test "Empty Email Registration" \
    "curl -s -w '%{http_code}' -X POST $BASE_URL/register -H 'Content-Type: application/json' -d '{\"email\":\"\",\"password\":\"pass123\",\"firstname\":\"John\",\"lastname\":\"Doe\",\"phone\":\"0812345678\",\"birthday\":\"1990-01-01\"}'" \
    "400" \
    "Test registration with empty email field"

run_advanced_test "Invalid Email Format" \
    "curl -s -w '%{http_code}' -X POST $BASE_URL/register -H 'Content-Type: application/json' -d '{\"email\":\"invalid-email\",\"password\":\"pass123\",\"firstname\":\"John\",\"lastname\":\"Doe\",\"phone\":\"0812345678\",\"birthday\":\"1990-01-01\"}'" \
    "201" \
    "Test registration with invalid email format (should still work as no email validation implemented)"

run_advanced_test "Invalid Birthday Format" \
    "curl -s -w '%{http_code}' -X POST $BASE_URL/register -H 'Content-Type: application/json' -d '{\"email\":\"test.birthday@example.com\",\"password\":\"pass123\",\"firstname\":\"John\",\"lastname\":\"Doe\",\"phone\":\"0812345678\",\"birthday\":\"invalid-date\"}'" \
    "400" \
    "Test registration with invalid birthday format"

run_advanced_test "Very Long Password" \
    "curl -s -w '%{http_code}' -X POST $BASE_URL/register -H 'Content-Type: application/json' -d '{\"email\":\"longpass@example.com\",\"password\":\"$(printf '%.500s' 'a')\",\"firstname\":\"John\",\"lastname\":\"Doe\",\"phone\":\"0812345678\",\"birthday\":\"1990-01-01\"}'" \
    "201" \
    "Test registration with very long password"

# Test Suite 2: Authentication & Authorization
echo -e "\n${PURPLE}📋 Test Suite 2: Authentication & Authorization${NC}"

# First, register a test user
curl -s -X POST "$BASE_URL/register" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\",\"firstname\":\"Test\",\"lastname\":\"User\",\"phone\":\"0812345678\",\"birthday\":\"1990-01-01\"}" > /dev/null

# Login to get token
JWT_TOKEN=$(curl -s -X POST "$BASE_URL/login" \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"$TEST_EMAIL\",\"password\":\"$TEST_PASSWORD\"}" | \
    jq -r '.token' 2>/dev/null)

run_advanced_test "Malformed JWT Token" \
    "curl -s -w '%{http_code}' -X GET $BASE_URL/me -H 'Authorization: Bearer malformed.jwt.token'" \
    "401" \
    "Test access with malformed JWT token"

run_advanced_test "Empty Bearer Token" \
    "curl -s -w '%{http_code}' -X GET $BASE_URL/me -H 'Authorization: Bearer '" \
    "401" \
    "Test access with empty bearer token"

run_advanced_test "Wrong Authorization Format" \
    "curl -s -w '%{http_code}' -X GET $BASE_URL/me -H 'Authorization: $JWT_TOKEN'" \
    "401" \
    "Test access without 'Bearer' prefix"

# Test Suite 3: HTTP Methods
echo -e "\n${PURPLE}📋 Test Suite 3: HTTP Methods${NC}"

run_advanced_test "GET on POST-only Endpoint" \
    "curl -s -w '%{http_code}' -X GET $BASE_URL/register" \
    "405" \
    "Test GET method on registration endpoint (should be Method Not Allowed)"

run_advanced_test "PUT on POST Endpoint" \
    "curl -s -w '%{http_code}' -X PUT $BASE_URL/login -H 'Content-Type: application/json' -d '{\"email\":\"test@example.com\",\"password\":\"pass123\"}'" \
    "405" \
    "Test PUT method on login endpoint"

run_advanced_test "DELETE on GET Endpoint" \
    "curl -s -w '%{http_code}' -X DELETE $BASE_URL/" \
    "405" \
    "Test DELETE method on hello endpoint"

# Test Suite 4: Content Type Validation
echo -e "\n${PURPLE}📋 Test Suite 4: Content Type Validation${NC}"

run_advanced_test "Missing Content-Type" \
    "curl -s -w '%{http_code}' -X POST $BASE_URL/register -d '{\"email\":\"noheader@example.com\",\"password\":\"pass123\",\"firstname\":\"John\",\"lastname\":\"Doe\",\"phone\":\"0812345678\",\"birthday\":\"1990-01-01\"}'" \
    "400" \
    "Test POST without Content-Type header"

run_advanced_test "Wrong Content-Type" \
    "curl -s -w '%{http_code}' -X POST $BASE_URL/register -H 'Content-Type: text/plain' -d 'invalid data'" \
    "400" \
    "Test POST with wrong Content-Type"

# Test Suite 5: Large Payload
echo -e "\n${PURPLE}📋 Test Suite 5: Large Payload Testing${NC}"

# Create large JSON payload
LARGE_FIRSTNAME=$(printf '%.1000s' 'A')
run_advanced_test "Large Payload" \
    "curl -s -w '%{http_code}' -X POST $BASE_URL/register -H 'Content-Type: application/json' -d '{\"email\":\"large@example.com\",\"password\":\"pass123\",\"firstname\":\"$LARGE_FIRSTNAME\",\"lastname\":\"Doe\",\"phone\":\"0812345678\",\"birthday\":\"1990-01-01\"}'" \
    "201" \
    "Test registration with large firstname field"

# Test Suite 6: Concurrent Requests
echo -e "\n${PURPLE}📋 Test Suite 6: Concurrent Requests${NC}"

echo -e "${YELLOW}Running $CONCURRENT_REQUESTS concurrent requests...${NC}"
concurrent_results=()

for i in $(seq 1 $CONCURRENT_REQUESTS); do
    (
        response=$(curl -s -w "%{http_code}" -X GET "$BASE_URL/" 2>/dev/null)
        echo "$response" > "/tmp/concurrent_$i.tmp"
    ) &
done

wait

concurrent_passed=0
for i in $(seq 1 $CONCURRENT_REQUESTS); do
    if [ -f "/tmp/concurrent_$i.tmp" ]; then
        status=$(cat "/tmp/concurrent_$i.tmp" | tail -c 4)
        if [ "$status" = "200" ]; then
            concurrent_passed=$((concurrent_passed + 1))
        fi
        rm "/tmp/concurrent_$i.tmp"
    fi
done

echo -e "Concurrent requests: $concurrent_passed/$CONCURRENT_REQUESTS passed"
if [ $concurrent_passed -eq $CONCURRENT_REQUESTS ]; then
    echo -e "${GREEN}✅ Concurrent test PASSED${NC}"
    count_test "PASSED"
else
    echo -e "${RED}❌ Concurrent test FAILED${NC}"
    count_test "FAILED"
fi

# Test Suite 7: API Performance
echo -e "\n${PURPLE}📋 Test Suite 7: Performance Testing${NC}"

echo -e "${YELLOW}Testing response times...${NC}"
for endpoint in "/" "/swagger/" "/non-existent"; do
    response_time=$(curl -s -w "%{time_total}" -o /dev/null "$BASE_URL$endpoint" 2>/dev/null)
    echo -e "$endpoint: ${response_time}s"
done

# Cleanup
rm -f /tmp/test_output.json /tmp/concurrent_*.tmp

# Final Summary
echo -e "\n${BLUE}📊 Advanced Test Results${NC}"
echo -e "${BLUE}========================${NC}"
echo -e "${GREEN}✅ Passed: $PASSED${NC}"
echo -e "${RED}❌ Failed: $FAILED${NC}"
echo -e "${YELLOW}📊 Total: $TOTAL${NC}"

if [ $FAILED -eq 0 ]; then
    echo -e "\n${GREEN}🎉 All tests passed!${NC}"
    exit 0
else
    echo -e "\n${YELLOW}⚠️  Some tests failed. Review the output above.${NC}"
    exit 1
fi
