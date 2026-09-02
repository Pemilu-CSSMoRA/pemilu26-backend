# Pemilu 2026 Backend

Backend API untuk sistem Pemilu internal CSSMoRA 2026.

Project ini dibangun menggunakan **Go**, **Gin**, **GORM**, dan **PostgreSQL**, dengan PostgreSQL dijalankan menggunakan Docker.

---

## Tech Stack

- Go
- Gin
- GORM
- PostgreSQL
- Docker
- JWT
- UUID

---

## Project Structure

```text
pemilu26-backend/
├── config/
├── constants/
├── internal/
│   ├── controller/
│   ├── dto/
│   ├── entity/
│   ├── handler/
│   ├── middleware/
│   ├── repository/
│   ├── route/
│   ├── service/
│   └── utils/
├── migration/
│   ├── migrate.go
│   └── seeder.go
├── cmd/
├── .env
├── .env.example
├── .gitignore
├── docker-compose.yaml
├── go.mod
├── go.sum
└── main.go
