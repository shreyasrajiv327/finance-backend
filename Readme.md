#  Finance Backend API

A backend service for managing personal financial records, including income and expenses.  
Built with **Go (Golang)** using the **Gin framework** and **PostgreSQL**, with **JWT authentication** and **Role-Based Access Control (RBAC)**.

---

##  Features

-  User authentication using JWT
-  CRUD operations for financial records
-  Summary APIs (total income, expenses, balance)
-  Category-wise expense breakdown
-  Recent transactions
-  Advanced filtering (type, category, date range, pagination)
-  Role-Based Access Control (admin, editor, viewer)
-  Optimized SQL queries with filtering support

---

##  Tech Stack

- **Language:** Go (Golang)
- **Framework:** Gin
- **Database:** PostgreSQL
- **Authentication:** JWT
- **Testing:** Bash script using `curl`

---

##  Project Structure

```
.
├── internal/handlers/                 # HTTP handlers (controllers)
├── internal/handlers/repository/      # Database queries
├── internal/handlers/models/          # Data models
├── internal/handlers/middleware/      # JWT & RBAC middleware
├── internal/handlers/routes/          # Route definitions
├── cmd/             # Entry point (contains main.go)
├── testing/         # API testing script files
└── API_DOCUMENTATION.md
└── Readme.md
```

---

##  Setup Instructions

### 1. Clone the repository

```bash
git clone git@github.com:shreyasrajiv327/finance-backend.git
cd <project-folder>
```

### 2. Setup PostgreSQL

Create a database:

```sql
CREATE DATABASE finance;
```

Update your DB credentials in the config file.

### 3. Install dependencies

```bash
go mod tidy
```

### 4. Run the server

```bash
go run cmd/main.go
```

Server will start at:

```
http://localhost:8080
```

---

##  Authentication

Login to get a JWT token:

```http
POST /login
```

Use the token in headers:

```http
Authorization: Bearer <your_token>
```

---

##  Testing APIs

Run the test script available in the testing folder:

```bash
chmod +x test.sh
./test.sh
```

Results will be saved in:

```
results.txt
```

---

##  Roles & Permissions

| Role | Permissions |
|------|-------------|
| `admin` | Full access |
| `editor` | Create, update, delete records |
| `viewer` | Read-only access |

---

##  Key APIs

- `POST /login` → Authenticate user
- `POST /records` → Create record
- `GET /records` → Get all records
- `GET /records/summary` → Financial summary
- `GET /records/category-summary` → Category breakdown
- `GET /records/recent` → Recent transactions
- `GET /records/filter` → Filter records
- `PUT /records/:id` → Update record
- `DELETE /records/:id` → Delete record

**For full details, see:**  
 [API_DOCUMENTATION.md](./API_Documentation.md)

---

##  Validation & Error Handling

- Input validation for amount, type, and query params
- Proper HTTP status codes (`200`, `201`, `400`, `403`, `404`, `500`)
- RBAC enforcement for restricted actions
- Protection against unauthorized data access

---

## ⚡ Design Decisions

- Repository pattern used for database interaction
- Handlers kept lightweight and focused on request/response
- SQL queries optimized with optional filters
- JWT used for stateless authentication
- Simplicity preferred over over-engineering

---

##  License

Copyright © 2026. All rights reserved.

---



