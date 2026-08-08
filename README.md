# Task Manager

A full-stack task management application with user authentication, built to demonstrate production-grade backend and frontend practices — layered architecture, comprehensive testing, containerization, and CI/CD.

![Backend CI](https://github.com/DanKRT-Star/task-manager/actions/workflows/backend-ci.yml/badge.svg)

## Features

- **User Authentication** — Register and log in with JWT-based auth; passwords hashed with bcrypt.
- **Task Management (CRUD)** — Create, read, update, and delete tasks, scoped strictly to the authenticated user.
- **Filtering, Sorting & Pagination** — Filter tasks by status, sort by deadline, and paginate results.
- **Rate Limiting** — Auth endpoints are rate-limited to mitigate brute-force attempts.
- **Containerized** — Full stack (frontend, backend, database) runs with a single `docker-compose up`.
- **CI Pipeline** — Every push runs build, vet, vulnerability scanning, race-condition testing, and Docker image builds.

## Tech Stack

**Backend**
- Go + [Fiber v3](https://github.com/gofiber/fiber) — HTTP framework
- [GORM](https://gorm.io/) + PostgreSQL — ORM and database
- JWT (`golang-jwt/jwt`) — authentication
- `bcrypt` — password hashing
- `go-playground/validator` — request validation
- `testify` — testing (unit + integration)
- `swaggo/swag` — auto-generated API documentation (Swagger/OpenAPI)

**Frontend**
- React + TypeScript + Vite
- Tailwind CSS
- React Router — client-side routing
- React Hook Form + Zod — form handling and schema validation
- Axios — HTTP client
- Sonner — toast notifications

**Infrastructure**
- Docker & Docker Compose
- GitHub Actions (CI)

## Architecture

The backend follows a layered architecture to separate concerns:

Route → Handler → Service → Repository → Database


- **Handler** — parses HTTP requests, validates input, formats responses
- **Service** — business logic (password hashing, status validation, authorization checks)
- **Repository** — database queries via GORM, exposed through interfaces for mockability
- **DTO** — request/response shapes, decoupled from database models

```
backend/
├── cmd/ # application entrypoint
├── internal/
│ ├── apperror/ # centralized error types
│ ├── config/ # database connection setup
│ ├── dto/ # request/response structs
│ ├── handler/ # HTTP handlers
│ ├── middleware/ # JWT auth middleware
│ ├── model/ # GORM models
│ ├── repository/ # data access layer (interfaces + implementations)
│ ├── route/v1/ # route registration
│ ├── service/ # business logic
│ └── validator/ # validation setup
├── test/integration/ # end-to-end API tests
└── Dockerfile
```
```
frontend/
├── src/
│ ├── components/ # reusable UI components
│ ├── context/ # React Context (auth state)
│ ├── hooks/ # custom hooks (useAuth, useTask)
│ ├── lib/ # validation schemas (Zod)
│ ├── pages/ # route-level components
│ ├── services/ # API client layer
│ └── types/ # TypeScript types
└── Dockerfile
```

## Getting Started

### Option 1 — Run with Docker (recommended)

Requires [Docker Desktop](https://www.docker.com/products/docker-desktop/).

1. Clone the repository:
```bash
   git clone https://github.com/DanKRT-Star/task-manager.git
   cd task-manager
```

2. Create a `.env` file in the project root:
```env
   DB_PASSWORD=your_db_password
   JWT_SECRET=your_jwt_secret
```

3. Start the full stack:
```bash
   docker-compose up --build
```

4. Access the app:
   - Frontend: [http://localhost:5173](http://localhost:5173)
   - Backend API: [http://localhost:3000](http://localhost:3000)

### Option 2 — Run locally (without Docker)

**Prerequisites:** Go 1.26+, Node.js 20+, PostgreSQL 16+

**Backend:**
```bash
cd backend
cp .env.example .env   # fill in DATABASE_URL, JWT_SECRET, PORT
go mod download
go run cmd/main.go
```

**Frontend:**
```bash
cd frontend
cp .env.example .env   # set VITE_API_URL
npm install
npm run dev
```

## Running Tests

```bash
cd backend

# Unit tests (mocked repositories, no database required)
go test ./internal/... -v

# Integration tests (requires a running PostgreSQL test database)
go test ./test/... -v

# All tests with race detector
go test ./... -race -v
```

Test coverage includes:
- Authentication (registration, login, validation, duplicate email handling)
- Task CRUD, filtering, sorting, and pagination
- **Authorization** — users cannot access, modify, or delete tasks belonging to other users
- Rate limiting behavior on auth endpoints
- Malformed/expired token handling

## API Overview

| Method | Endpoint | Description | Auth Required |
|--------|----------|--------------|:---:|
| POST | `/api/v1/auth/register` | Register a new user | ❌ |
| POST | `/api/v1/auth/login` | Log in, returns JWT | ❌ |
| GET | `/api/v1/tasks` | List tasks (supports `status`, `sort`, `page`, `limit`) | ✅ |
| POST | `/api/v1/tasks` | Create a task | ✅ |
| PUT | `/api/v1/tasks/:id` | Update a task | ✅ |
| DELETE | `/api/v1/tasks/:id` | Delete a task | ✅ |
| GET | `/health` | Health check (includes DB connectivity) | ❌ |

## API Documentation

Interactive API documentation (Swagger UI) is available once the backend is running:
http://localhost:3000/swagger/index.html

To test authenticated endpoints, first call `POST /auth/login` to obtain a JWT token, then click **Authorize** in the Swagger UI and enter `Bearer <your-token>`.

## CI/CD

Every push to `main` triggers a GitHub Actions pipeline that:
1. Builds the backend and runs `go vet`
2. Scans dependencies for known vulnerabilities (`govulncheck`)
3. Runs the full test suite with the race detector enabled
4. Builds both Docker images to catch containerization issues early

See [`.github/workflows/backend-ci.yml`](.github/workflows/backend-ci.yml).

## Notable Design Decisions

- **Authorization enforced at the repository layer** — every task query includes the owning `user_id`, preventing cross-user data access even if a check is missed elsewhere in the code.
- **Dependency injection via interfaces** — repositories are injected as interfaces into services, enabling unit tests with mocked dependencies instead of a live database.
- **Defense in depth for data integrity** — task status is validated both in application code and via a PostgreSQL `CHECK` constraint.
- **Separate rate-limit configuration for tests** — rate limiting is toggled via a parameter, avoiding flaky tests caused by shared IP-based limits in the test environment.

## License

This project is for portfolio purposes.