# 🧾 WidaTech Technical Challenge - Backend (Golang Fiber)
![Go Version](https://img.shields.io/badge/go-1.24-blue)
![Drone Build Status](https://drone-tencent.aldera.space/api/badges/kiminodare/HOVARLAY-BE/status.svg?branch=dev)

Backend service for Hovarlay, built using [Golang Fiber](https://gofiber.io/), [Entgo](https://entgo.io/), and PostgreSQL.

## 🚀 Getting Started

### ✅ Prerequisites

- [Go](https://golang.org/)
- [Fiber](https://gofiber.io/)
- [Entgo](https://entgo.io/)
- [PostgreSQL](https://www.postgresql.org/)
- [Argon2](https://github.com/elithrar/argon2)
- [Crypto](https://pkg.go.dev/golang.org/x/crypto)
- [UUID](https://pkg.go.dev/github.com/google/uuid)

### 1. Setup & Run Backend

```bash
go mod tidy         # Download dependencies
go generate ./ent   # Generate Entgo schema
go run cmd/migrate/main.go  # Run database migration
go run cmd/server/main.go   # Run server
```

> ⚙️ Make sure the `.env` file is created and points to the PostgreSQL database. \
> 💡 The application will run on `http://localhost:9888` (or the port specified in the `.env` file).

Example `.env`:

```dotenv
# Database
DB_HOST=localhost
DB_PORT=5432
DB_USER=hovarlay_admin
DB_PASSWORD=password
DB_NAME=hovarlay
DB_SSLMODE=disable

# Server
SERVER_PORT=9888
APP_ENV=development  # local/development/production
ALLOWED_ORIGINS=http://localhost:4004,http://localhost:3000

# JWT
JWT_SECRET=hovarlay_secret

#Aes
AES_KEY=hovarlay_aes_key

```

---

## ✨ Features

- 🔐 Authentication and authorization using JWT
- 🔑 Encrypt password using Argon2
- 🔒 Encrypt sensitive data using AES
- 🔍 Search and filter data
- 🔄 Pagination and lazy loading

---

## 🛠 Tech Stack

- **Runtime:** Go 1.24
- **Libraries:**
  - Golang Fiber (Framework)
  - Entgo (ORM)
  - PostgreSQL (Database)
  - Argon2
  - Crypto
  - UUID
  - JWT
  - AES
---

## 👨‍💻 Author

**Kiminodare**  
Fullstack Engineer
