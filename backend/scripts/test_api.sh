#!/bin/bash

# Colors for output
GREEN='\033[0;32m'
RED='\033[0;31m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

BASE_URL="http://localhost:8080/api/v1"
BASE_PATH="http://localhost:8080"

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}   JLPT5 Backend API Test Suite${NC}"
echo -e "${BLUE}========================================${NC}\n"

# Test 1: Health Check
echo -e "${BLUE}Test 1: Health Check${NC}"
response=$(curl -s ${BASE_PATH}/health)
if echo "$response" | grep -q "healthy"; then
    echo -e "${GREEN}✅ PASSED${NC}: Health check endpoint working\n"
else
    echo -e "${RED}❌ FAILED${NC}: Health check failed\n"
    exit 1
fi

# Test 2: User Registration
echo -e "${BLUE}Test 2: User Registration${NC}"
TIMESTAMP=$(date +%s)
response=$(curl -s -X POST ${BASE_URL}/auth/register \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"user${TIMESTAMP}@test.com\",\"username\":\"user${TIMESTAMP}\",\"password\":\"password123\"}")
if echo "$response" | grep -q "access_token"; then
    echo -e "${GREEN}✅ PASSED${NC}: User registration successful\n"
    TOKEN=$(echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
else
    echo -e "${RED}❌ FAILED${NC}: Registration failed\n"
    echo "Response: $response"
    exit 1
fi

# Test 3: User Login
echo -e "${BLUE}Test 3: User Login${NC}"
response=$(curl -s -X POST ${BASE_URL}/auth/login \
    -H "Content-Type: application/json" \
    -d '{"email":"test@example.com","password":"password123"}')
if echo "$response" | grep -q "access_token"; then
    echo -e "${GREEN}✅ PASSED${NC}: Login successful\n"
    TOKEN=$(echo "$response" | grep -o '"access_token":"[^"]*"' | cut -d'"' -f4)
else
    echo -e "${RED}❌ FAILED${NC}: Login failed\n"
    exit 1
fi

# Test 4: List Vocabulary
echo -e "${BLUE}Test 4: List Vocabulary${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/vocabulary)
if echo "$response" | grep -q '"items"'; then
    count=$(echo "$response" | grep -o '"id":' | wc -l)
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved $count vocabulary items\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to list vocabulary\n"
    exit 1
fi

# Test 5: Get Vocabulary Item
echo -e "${BLUE}Test 5: Get Vocabulary Item${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/vocabulary/1)
if echo "$response" | grep -q '"word"'; then
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved vocabulary item\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to get vocabulary item\n"
    exit 1
fi

# Test 6: List Grammar Lessons
echo -e "${BLUE}Test 6: List Grammar Lessons${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/grammar)
if echo "$response" | grep -q '"items"'; then
    count=$(echo "$response" | grep -o '"id":' | wc -l)
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved $count grammar lessons\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to list grammar lessons\n"
    exit 1
fi

# Test 7: Get Grammar Lesson
echo -e "${BLUE}Test 7: Get Grammar Lesson${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/grammar/1)
if echo "$response" | grep -q '"grammar_point"'; then
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved grammar lesson with examples\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to get grammar lesson\n"
    exit 1
fi

# Test 8: Mark Grammar Lesson as Completed
echo -e "${BLUE}Test 8: Mark Grammar Lesson as Completed${NC}"
response=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"notes":"Great lesson!"}' \
    ${BASE_URL}/grammar/1/complete)
if echo "$response" | grep -q '"success"'; then
    echo -e "${GREEN}✅ PASSED${NC}: Grammar lesson marked as completed\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to mark grammar lesson as completed\n"
    exit 1
fi

# Test 9: List Quizzes
echo -e "${BLUE}Test 9: List Quizzes${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/quizzes)
if echo "$response" | grep -q '"items"'; then
    count=$(echo "$response" | grep -o '"id":' | wc -l)
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved $count quizzes\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to list quizzes\n"
    exit 1
fi

# Test 10: Start Quiz Session
echo -e "${BLUE}Test 10: Start Quiz Session${NC}"
response=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" ${BASE_URL}/quizzes/1/start)
if echo "$response" | grep -q '"session_id"'; then
    SESSION_ID=$(echo "$response" | grep -o '"session_id":[0-9]*' | cut -d':' -f2)
    echo -e "${GREEN}✅ PASSED${NC}: Quiz session started (Session ID: $SESSION_ID)\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to start quiz session\n"
    exit 1
fi

# Test 11: Submit Quiz Answers
echo -e "${BLUE}Test 11: Submit Quiz Answers${NC}"
response=$(curl -s -X POST -H "Authorization: Bearer $TOKEN" \
    -H "Content-Type: application/json" \
    -d '{"answers":{"1":"A","2":"B","3":"C","4":"B","5":"D"}}' \
    ${BASE_URL}/quizzes/sessions/${SESSION_ID}/submit)
if echo "$response" | grep -q '"score"'; then
    score=$(echo "$response" | grep -o '"percentage":[0-9.]*' | cut -d':' -f2)
    echo -e "${GREEN}✅ PASSED${NC}: Quiz submitted successfully (Score: ${score}%)\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to submit quiz\n"
    exit 1
fi

# Test 12: Get Quiz Results
echo -e "${BLUE}Test 12: Get Quiz Results${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/quizzes/sessions/${SESSION_ID})
if echo "$response" | grep -q '"session_id"'; then
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved quiz results\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to get quiz results\n"
    exit 1
fi

# Test 13: Get Quiz History
echo -e "${BLUE}Test 13: Get Quiz History${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/quizzes/history)
if echo "$response" | grep -q '"sessions"'; then
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved quiz history\n"
else
    echo -e "${RED}❌ FAILED${NC}: Failed to get quiz history\n"
    exit 1
fi

# Test 14: Get Progress Statistics
echo -e "${BLUE}Test 14: Get Progress Statistics${NC}"
response=$(curl -s -H "Authorization: Bearer $TOKEN" ${BASE_URL}/progress/stats)
if echo "$response" | grep -q '"study_streak_days"'; then
    echo -e "${GREEN}✅ PASSED${NC}: Retrieved progress statistics\n"
    echo "Statistics preview:"
    echo "$response" | grep -o '"vocabulary_learned":[0-9]*' | head -1
    echo "$response" | grep -o '"grammar_completed":[0-9]*' | head -1
    echo "$response" | grep -o '"quizzes_taken":[0-9]*' | head -1
else
    echo -e "${RED}❌ FAILED${NC}: Failed to get progress statistics\n"
    exit 1
fi

echo -e "\n${GREEN}========================================${NC}"
echo -e "${GREEN}   All Tests Passed! ✅${NC}"
echo -e "${GREEN}========================================${NC}\n"
