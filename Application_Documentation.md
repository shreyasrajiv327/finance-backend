# Finance Backend API Documentation

## Overview

The Finance Backend API provides a comprehensive RESTful interface for managing personal finance records, including income and expense tracking, categorization, and financial summaries.


---

## Table of Contents

- [Authentication](#authentication)
- [Authorization & Roles](#authorization--roles)
- [Endpoints Summary](#endpoints-summary)
- [Error Handling](#error-handling)
- [API Endpoints](#api-endpoints)
  - [Authentication](#1-authentication)
  - [Records Management](#records-management)
  - [Analytics & Summaries](#analytics--summaries)
- [Data Models](#data-models)
- [Status Codes](#status-codes)

---

## Authentication

All endpoints (except `/login`) require authentication using JWT (JSON Web Token) bearer tokens.

### Authentication Header Format

```http
Authorization: Bearer <your_jwt_token>
```

### Obtaining a Token

To obtain an authentication token, use the [Login endpoint](#1-login).

**Token Expiration:** Tokens are valid for 24 hours by default.

---

## Authorization & Roles

The API implements role-based access control (RBAC) with three permission levels:

| Role | Permissions |
|------|-------------|
| `admin` | Full access to all operations |
| `editor` | Can create, read, update, and delete records |
| `viewer` | Read-only access (cannot modify data) |

---

## Endpoints Summary

| Method | Endpoint | Description | Auth Required | Roles |
|--------|----------|-------------|---------------|-------|
| POST | `/login` | Authenticate user and obtain JWT token | No | All |
| POST | `/records` | Create a new financial record | Yes | `editor`, `admin` |
| GET | `/records` | Retrieve all records for authenticated user | Yes | All |
| GET | `/records/:id` | Retrieve a specific record by ID | Yes | All |
| PUT | `/records/:id` | Update an existing record | Yes | `editor`, `admin` |
| DELETE | `/records/:id` | Delete a record | Yes | `editor`, `admin` |
| GET | `/records/summary` | Get monthly income/expense summary | Yes | All |
| GET | `/records/category-summary` | Get spending breakdown by category | Yes | All |
| GET | `/records/recent` | Get most recent records | Yes | All |
| GET | `/records/filter` | Filter records with query parameters | Yes | All |

---

## Error Handling

The API uses standard HTTP status codes and returns error details in JSON format.

### Error Response Format

```json
{
  "error": "Error message describing what went wrong",
  "code": "ERROR_CODE",
  "details": {
    "field": "Additional context if applicable"
  }
}
```

### Common Error Responses

| Status Code | Description |
|-------------|-------------|
| `400 Bad Request` | Invalid request payload or parameters |
| `401 Unauthorized` | Missing or invalid authentication token |
| `403 Forbidden` | Insufficient permissions for the operation |
| `404 Not Found` | Requested resource does not exist |
| `500 Internal Server Error` | Server-side error occurred |

---

## API Endpoints

### 1. Login

Authenticate a user and receive a JWT token for subsequent requests.

**Endpoint:** `POST /login`  
**Authentication Required:** No

#### Request Body

```json
{
  "email": "editor@test.com",
  "password": "123456"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `email` | string | Yes | User's email address |
| `password` | string | Yes | User's password |

#### Success Response

**Status Code:** `200 OK`

```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "user": {
    "id": 2,
    "email": "editor@test.com",
    "role": "editor"
  },
  "expires_at": "2026-04-06T12:22:07Z"
}
```

#### Error Responses

**Status Code:** `401 Unauthorized`

```json
{
  "error": "Invalid credentials",
  "code": "INVALID_CREDENTIALS"
}
```

#### Example Request

```bash
curl -X POST http://localhost:8080/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "editor@test.com",
    "password": "123456"
  }'
```

---

## Records Management

### 2. Create Record

Create a new income or expense record.

**Endpoint:** `POST /records`  
**Authentication Required:** Yes  
**Authorized Roles:** `editor`, `admin`

#### Request Body

```json
{
  "amount": 200.50,
  "type": "expense",
  "category": "food",
  "date": "2026-04-05T12:22:07Z",
  "notes": "Lunch at downtown restaurant"
}
```

| Field | Type | Required | Description |
|-------|------|----------|-------------|
| `amount` | number | Yes | Transaction amount (positive number) |
| `type` | string | Yes | Transaction type: `income` or `expense` |
| `category` | string | Yes | Category name (e.g., "food", "rent", "salary") |
| `date` | string (ISO 8601) | No | Transaction date (defaults to current time) |
| `notes` | string | No | Additional notes or description |

#### Success Response

**Status Code:** `201 Created`

```json
{
  "id": 15,
  "user_id": 2,
  "amount": 200.50,
  "type": "expense",
  "category": "food",
  "date": "2026-04-05T12:22:07Z",
  "notes": "Lunch at downtown restaurant",
  "created_at": "2026-04-05T12:22:07Z",
  "updated_at": "2026-04-05T12:22:07Z"
}
```

#### Error Responses

**Status Code:** `400 Bad Request`

```json
{
  "error": "Invalid input",
  "code": "VALIDATION_ERROR",
  "details": {
    "amount": "Amount must be a positive number",
    "type": "Type must be either 'income' or 'expense'"
  }
}
```

**Status Code:** `403 Forbidden`

```json
{
  "error": "Insufficient permissions",
  "code": "FORBIDDEN"
}
```

#### Example Request

```bash
curl -X POST http://localhost:8080/records \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 200.50,
    "type": "expense",
    "category": "food",
    "notes": "Lunch at downtown restaurant"
  }'
```

---

### 3. Get All Records

Retrieve all financial records for the authenticated user.

**Endpoint:** `GET /records`  
**Authentication Required:** Yes  
**Authorized Roles:** All

#### Success Response

**Status Code:** `200 OK`

```json
[
  {
    "id": 15,
    "user_id": 2,
    "amount": 200.50,
    "type": "expense",
    "category": "food",
    "date": "2026-04-05T12:22:07Z",
    "notes": "Lunch at downtown restaurant",
    "created_at": "2026-04-05T12:22:07Z",
    "updated_at": "2026-04-05T12:22:07Z"
  },
  {
    "id": 14,
    "user_id": 2,
    "amount": 5000.00,
    "type": "income",
    "category": "salary",
    "date": "2026-04-01T00:00:00Z",
    "notes": "Monthly salary",
    "created_at": "2026-04-01T08:30:00Z",
    "updated_at": "2026-04-01T08:30:00Z"
  }
]
```

#### Example Request

```bash
curl -X GET http://localhost:8080/records \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 4. Get Record by ID

Retrieve a specific financial record by its unique identifier.

**Endpoint:** `GET /records/:id`  
**Authentication Required:** Yes  
**Authorized Roles:** All

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Unique record identifier |

#### Success Response

**Status Code:** `200 OK`

```json
{
  "id": 15,
  "user_id": 2,
  "amount": 200.50,
  "type": "expense",
  "category": "food",
  "date": "2026-04-05T12:22:07Z",
  "notes": "Lunch at downtown restaurant",
  "created_at": "2026-04-05T12:22:07Z",
  "updated_at": "2026-04-05T12:22:07Z"
}
```

#### Error Responses

**Status Code:** `403 Forbidden`

```json
{
  "error": "Access denied",
  "code": "FORBIDDEN",
  "message": "This record belongs to another user"
}
```

**Status Code:** `404 Not Found`

```json
{
  "error": "Record not found",
  "code": "NOT_FOUND"
}
```

#### Example Request

```bash
curl -X GET http://localhost:8080/records/15 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 5. Update Record

Update an existing financial record.

**Endpoint:** `PUT /records/:id`  
**Authentication Required:** Yes  
**Authorized Roles:** `editor`, `admin`

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Unique record identifier |

#### Request Body

All fields are optional. Only include fields you want to update.

```json
{
  "amount": 250.00,
  "category": "transport",
  "notes": "Updated notes"
}
```

#### Success Response

**Status Code:** `200 OK`

```json
{
  "id": 15,
  "user_id": 2,
  "amount": 250.00,
  "type": "expense",
  "category": "transport",
  "date": "2026-04-05T12:22:07Z",
  "notes": "Updated notes",
  "created_at": "2026-04-05T12:22:07Z",
  "updated_at": "2026-04-05T14:30:00Z"
}
```

#### Error Responses

**Status Code:** `400 Bad Request`

```json
{
  "error": "Invalid input",
  "code": "VALIDATION_ERROR"
}
```

**Status Code:** `403 Forbidden`

```json
{
  "error": "Access denied",
  "code": "FORBIDDEN",
  "message": "This record belongs to another user"
}
```

**Status Code:** `404 Not Found`

```json
{
  "error": "Record not found",
  "code": "NOT_FOUND"
}
```

#### Example Request

```bash
curl -X PUT http://localhost:8080/records/15 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN" \
  -H "Content-Type: application/json" \
  -d '{
    "amount": 250.00,
    "category": "transport"
  }'
```

---

### 6. Delete Record

Permanently delete a financial record.

**Endpoint:** `DELETE /records/:id`  
**Authentication Required:** Yes  
**Authorized Roles:** `editor`, `admin`

#### Path Parameters

| Parameter | Type | Description |
|-----------|------|-------------|
| `id` | integer | Unique record identifier |

#### Success Response

**Status Code:** `200 OK`

```json
{
  "message": "Record deleted successfully",
  "id": 15
}
```

#### Error Responses

**Status Code:** `403 Forbidden`

```json
{
  "error": "Access denied",
  "code": "FORBIDDEN",
  "message": "This record belongs to another user"
}
```

**Status Code:** `404 Not Found`

```json
{
  "error": "Record not found",
  "code": "NOT_FOUND"
}
```

#### Example Request

```bash
curl -X DELETE http://localhost:8080/records/15 \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Analytics & Summaries

### 7. Get Monthly Summary

Retrieve a summary of income, expenses, and balance for the current month.

**Endpoint:** `GET /records/summary`  
**Authentication Required:** Yes  
**Authorized Roles:** All

#### Success Response

**Status Code:** `200 OK`

```json
{
  "period": "2026-04",
  "total_income": 161000.00,
  "total_expense": 38300.50,
  "balance": 122699.50,
  "transaction_count": {
    "income": 12,
    "expense": 48
  }
}
```

| Field | Type | Description |
|-------|------|-------------|
| `period` | string | Month in YYYY-MM format |
| `total_income` | number | Sum of all income transactions |
| `total_expense` | number | Sum of all expense transactions |
| `balance` | number | Net balance (income - expense) |
| `transaction_count` | object | Number of income and expense records |

#### Example Request

```bash
curl -X GET http://localhost:8080/records/summary \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 8. Get Category Summary

Retrieve spending breakdown by category for the current month.

**Endpoint:** `GET /records/category-summary`  
**Authentication Required:** Yes  
**Authorized Roles:** All

#### Success Response

**Status Code:** `200 OK`

```json
[
  {
    "category": "rent",
    "total": 36000.00,
    "percentage": 93.99,
    "transaction_count": 1
  },
  {
    "category": "food",
    "total": 2000.50,
    "percentage": 5.22,
    "transaction_count": 15
  },
  {
    "category": "transport",
    "total": 300.00,
    "percentage": 0.78,
    "transaction_count": 8
  }
]
```

Categories are sorted by total amount in descending order.

#### Example Request

```bash
curl -X GET http://localhost:8080/records/category-summary \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 9. Get Recent Records

Retrieve the most recently created financial records.

**Endpoint:** `GET /records/recent`  
**Authentication Required:** Yes  
**Authorized Roles:** All

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `limit` | integer | No | 10 | Number of records to return (max: 100) |

#### Success Response

**Status Code:** `200 OK`

Returns an array of records sorted by creation time (newest first).

```json
[
  {
    "id": 20,
    "user_id": 2,
    "amount": 50.00,
    "type": "expense",
    "category": "entertainment",
    "date": "2026-04-05T18:00:00Z",
    "notes": "Movie tickets",
    "created_at": "2026-04-05T18:15:00Z",
    "updated_at": "2026-04-05T18:15:00Z"
  },
  {
    "id": 19,
    "user_id": 2,
    "amount": 150.00,
    "type": "expense",
    "category": "groceries",
    "date": "2026-04-05T15:30:00Z",
    "notes": "Weekly shopping",
    "created_at": "2026-04-05T16:00:00Z",
    "updated_at": "2026-04-05T16:00:00Z"
  }
]
```

#### Example Request

```bash
curl -X GET "http://localhost:8080/records/recent?limit=5" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

### 10. Get Filtered Records

Filter records based on multiple criteria with pagination support.

**Endpoint:** `GET /records/filter`  
**Authentication Required:** Yes  
**Authorized Roles:** All

#### Query Parameters

| Parameter | Type | Required | Default | Description |
|-----------|------|----------|---------|-------------|
| `type` | string | No | - | Filter by type: `income` or `expense` |
| `category` | string | No | - | Filter by category name |
| `startDate` | string | No | - | Start date (YYYY-MM-DD format) |
| `endDate` | string | No | - | End date (YYYY-MM-DD format) |
| `limit` | integer | No | 20 | Number of records per page (max: 100) |
| `offset` | integer | No | 0 | Number of records to skip (for pagination) |

#### Success Response

**Status Code:** `200 OK`

```json
{
  "records": [
    {
      "id": 15,
      "user_id": 2,
      "amount": 200.50,
      "type": "expense",
      "category": "food",
      "date": "2026-04-05T12:22:07Z",
      "notes": "Lunch",
      "created_at": "2026-04-05T12:22:07Z",
      "updated_at": "2026-04-05T12:22:07Z"
    }
  ],
  "pagination": {
    "total": 48,
    "limit": 20,
    "offset": 0,
    "has_more": true
  }
}
```

#### Example Requests

**Filter by type and category:**
```bash
curl -X GET "http://localhost:8080/records/filter?type=expense&category=food" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**Filter by date range:**
```bash
curl -X GET "http://localhost:8080/records/filter?startDate=2026-04-01&endDate=2026-04-05" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

**With pagination:**
```bash
curl -X GET "http://localhost:8080/records/filter?limit=10&offset=20" \
  -H "Authorization: Bearer YOUR_JWT_TOKEN"
```

---

## Data Models

### Record Object

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Unique record identifier |
| `user_id` | integer | ID of the user who owns this record |
| `amount` | number | Transaction amount (always positive) |
| `type` | string | Transaction type: `income` or `expense` |
| `category` | string | Category name (user-defined) |
| `date` | string (ISO 8601) | Transaction date and time |
| `notes` | string | Optional notes or description |
| `created_at` | string (ISO 8601) | Record creation timestamp |
| `updated_at` | string (ISO 8601) | Last update timestamp |

### User Object

| Field | Type | Description |
|-------|------|-------------|
| `id` | integer | Unique user identifier |
| `email` | string | User's email address |
| `role` | string | User role: `admin`, `editor`, or `viewer` |

---

## Status Codes

The API uses standard HTTP status codes:

| Code | Description |
|------|-------------|
| `200 OK` | Request succeeded |
| `201 Created` | Resource successfully created |
| `400 Bad Request` | Invalid request parameters or body |
| `401 Unauthorized` | Authentication required or failed |
| `403 Forbidden` | Insufficient permissions |
| `404 Not Found` | Resource not found |
| `500 Internal Server Error` | Server error occurred |

---
