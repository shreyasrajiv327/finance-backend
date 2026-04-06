#!/bin/bash

BASE_URL="http://localhost:8080"
OUTPUT_FILE="results.txt"

# Clear file first
> $OUTPUT_FILE

echo " Logging in users..." >> $OUTPUT_FILE

ADMIN_TOKEN=$(curl -s -X POST "$BASE_URL/login" \
-H "Content-Type: application/json" \
-d '{"email":"admin@test.com","password":"123456"}' | jq -r '.token')

EDITOR_TOKEN=$(curl -s -X POST "$BASE_URL/login" \
-H "Content-Type: application/json" \
-d '{"email":"editor@test.com","password":"123456"}' | jq -r '.token')

VIEWER_TOKEN=$(curl -s -X POST "$BASE_URL/login" \
-H "Content-Type: application/json" \
-d '{"email":"viewer@test.com","password":"123456"}' | jq -r '.token')

echo "Tokens fetched" >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE
echo " Testing APIs..." >> $OUTPUT_FILE
echo "==================================" >> $OUTPUT_FILE

# 🔹 Summary
echo "" >> $OUTPUT_FILE
echo "--- SUMMARY (EDITOR) ---" >> $OUTPUT_FILE
curl -s "$BASE_URL/records/summary" \
-H "Authorization: Bearer $EDITOR_TOKEN" >> $OUTPUT_FILE

# 🔹 Category Summary
echo "" >> $OUTPUT_FILE
echo "--- CATEGORY SUMMARY ---" >> $OUTPUT_FILE
curl -s "$BASE_URL/records/category-summary" \
-H "Authorization: Bearer $EDITOR_TOKEN" >> $OUTPUT_FILE

#  Recent
echo "" >> $OUTPUT_FILE
echo "--- RECENT ---" >> $OUTPUT_FILE
curl -s "$BASE_URL/records/recent" \
-H "Authorization: Bearer $EDITOR_TOKEN" >> $OUTPUT_FILE

#  RBAC TEST (Viewer should FAIL)
echo "" >> $OUTPUT_FILE
echo "--- RBAC TEST (VIEWER CREATE RECORD) ---" >> $OUTPUT_FILE
curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/records" \
-H "Authorization: Bearer $VIEWER_TOKEN" \
-H "Content-Type: application/json" \
-d '{"amount":100,"type":"expense","category":"food"}' >> $OUTPUT_FILE

#  RBAC TEST (Editor should PASS)
echo "" >> $OUTPUT_FILE
echo "--- RBAC TEST (EDITOR CREATE RECORD) ---" >> $OUTPUT_FILE
curl -s -o /dev/null -w "%{http_code}" -X POST "$BASE_URL/records" \
-H "Authorization: Bearer $EDITOR_TOKEN" \
-H "Content-Type: application/json" \
-d '{"amount":200,"type":"expense","category":"food"}' >> $OUTPUT_FILE

echo "" >> $OUTPUT_FILE
echo " Tests completed." >> $OUTPUT_FILE

#  Filtered Records
echo "" >> $OUTPUT_FILE
echo "--- FILTERED RECORDS (EDITOR) ---" >> $OUTPUT_FILE
curl -s "$BASE_URL/records/filtered?type=expense&category=food&startDate=2026-01-01&endDate=2026-12-31&limit=5&offset=0" \
-H "Authorization: Bearer $EDITOR_TOKEN" >> $OUTPUT_FILE