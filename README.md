# Task Manager

A full-stack project & task management application built to demonstrate production-grade backend and frontend practices — layered architecture, role-based authorization, comprehensive testing, structured logging, containerization, and CI/CD.

![Backend CI](https://github.com/DanKRT-Star/task-manager/actions/workflows/backend-ci.yml/badge.svg)

## Live Demo

- **App:** https://task-manager-phi-one-73.vercel.app
- **API:** https://task-manager-backend-t3to.onrender.com
- **API Docs (Swagger):** https://task-manager-backend-t3to.onrender.com/swagger/index.html

> Note: the backend is hosted on a free tier and may take 30–50 seconds to wake up on the first request after a period of inactivity.

## Features

- **Authentication** — JWT access tokens (15 min) + rotating refresh tokens (7 days, hashed at rest, revocable). Passwords hashed with bcrypt.
- **Projects** — Create projects with owner/member roles. Owners manage membership and settings; members collaborate on tasks.
- **Tasks** — Full CRUD, either standalone (personal) or attached to a project, optionally organized under an Epic, Milestone, and Sprint, and assigned to a project member.
- **Epics & Milestones** — Group tasks by theme (Epic) or release/deadline (Milestone) within a project.
- **Sprints** — Time-boxed iterations with status transitions (`planned` → `active` → `completed`).
- **Comments** — Threaded discussion on tasks; only the author can delete their own comment.
- **Labels** — Project-scoped tags attachable to tasks (many-to-many).
- **Activity Log** — Automatic audit trail per task (creation, status changes), viewable by anyone with access to the task.
- **Authorization model** — Two-tier project roles (`owner`/`member`); task-level permission checks ensure members only modify tasks they created or were assigned.
- **Filtering, Sorting & Pagination** — Available on both personal and project task lists.
- **Rate Limiting** — Global limiter plus a stricter one on auth endpoints to mitigate brute-force attempts.
- **Structured Logging** — JSON logs (Zerolog) with a consistent, named event catalog across every service.
- **Graceful Shutdown** — In-flight requests are given time to finish before the process exits; database connections close cleanly.
- **Automated DB Backups** — Daily backup workflow via GitHub Actions, independent of the hosting provider's own backups.
- **Containerized** — Full stack (frontend, backend, database) runs with a single `docker-compose up`.
- **CI Pipeline** — Every push runs build, vet, vulnerability scanning, race-condition testing, and Docker image builds.

## Tech Stack

**Backend**
- Go + [Fiber v3](https://github.com/gofiber/fiber) — HTTP framework
- [GORM](https://gorm.io/) + PostgreSQL — ORM and database
- JWT (`golang-jwt/jwt`) — access tokens; custom rotating refresh token store
- `bcrypt` — password hashing
- `go-playground/validator` — request validation
- [`zerolog`](https://github.com/rs/zerolog) — structured logging
- `testify` — testing (unit + integration)
- `swaggo/swag` — auto-generated API documentation (Swagger/OpenAPI)

**Frontend**
- React + TypeScript + Vite
- Tailwind CSS
- React Router — client-side routing
- React Hook Form + Zod — form handling and schema validation
- Axios — HTTP client with automatic access-token refresh on 401
- Sonner — toast notifications

**Infrastructure**
- Docker & Docker Compose
- GitHub Actions (CI + scheduled DB backup)
- Render (backend + Postgres) / Vercel (frontend)

## Architecture

The backend follows a layered architecture to separate concerns:

Route → Handler → Service → Repository → Database


- **Handler** — parses HTTP requests, validates input, formats responses
- **Service** — business logic (authorization rules, cross-entity validation, activity logging)
- **Repository** — database queries via GORM, exposed through interfaces for mockability
- **DTO** — request/response shapes, decoupled from database models
- **Logger** — named, per-domain event helpers (e.g. `logger.TaskCreated(...)`) instead of ad-hoc log statements, ensuring every logged event is discoverable and consistent

### Domain model
```
User ──< ProjectMember >── Project ──< Epic
├─< Milestone
├─< Sprint
├─< Label
└─< Task ──< Comment
│ ├─< ActivityLog
│ └─< Label (many-to-many)
├── Epic (optional)
├── Milestone (optional)
├── Sprint (optional)
└── Assignee (optional, must be a project member)
```


A task may exist standalone (no project) or within a project; Epic/Milestone/Sprint/Label assignment is only valid for project tasks, and each is validated to belong to the *same* project as the task.
```
backend/
├── cmd/ # application entrypoint (routing, graceful shutdown)
├── internal/
│ ├── apperror/ # centralized error types
│ ├── config/ # database connection, schema init/reset
│ ├── dto/ # request/response structs
│ ├── handler/ # HTTP handlers (one per domain)
│ ├── logger/ # structured logging: zerolog wrapper + per-domain event helpers
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
│ ├── components/ # reusable UI components (widgets, panels, badges)
│ ├── context/ # React Context (auth state)
│ ├── hooks/ # data-fetching hooks per domain (useProject, useTask, useEpic, ...)
│ ├── lib/ # validation schemas (Zod), shared error helpers
│ ├── pages/ # route-level components (home, projects, project detail, tasks, auth)
│ ├── services/ # API client layer (Axios, incl. auto token refresh)
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
- Authentication — registration, login, validation, duplicate email handling, **refresh token rotation** (reused tokens are rejected), logout revocation
- Task CRUD, filtering, sorting, and pagination — both personal and project-scoped
- **Authorization** — project owners vs. members; users cannot access, modify, or delete tasks/epics/milestones/sprints/comments/labels they don't have rights to
- Cross-entity validation — e.g. an epic/milestone/sprint/label from one project cannot be attached to a task in another
- Rate limiting behavior on auth endpoints
- Malformed/expired token handling

## API Overview

| Method | Endpoint | Description | Auth |
|--------|----------|--------------|:---:|
| POST | `/api/v1/auth/register` | Register a new user | ❌ |
| POST | `/api/v1/auth/login` | Log in, returns access + refresh token | ❌ |
| POST | `/api/v1/auth/refresh` | Exchange a refresh token for a new pair (rotation) | ❌ |
| POST | `/api/v1/auth/logout` | Revoke a refresh token | ❌ |
| GET | `/api/v1/auth/me` | Get the current user | ✅ |
| GET/POST | `/api/v1/tasks` | List / create personal tasks | ✅ |
| PUT/DELETE | `/api/v1/tasks/:id` | Update / delete a task | ✅ |
| GET/POST | `/api/v1/tasks/:taskId/comments` | List / add comments on a task | ✅ |
| DELETE | `/api/v1/comments/:commentId` | Delete own comment | ✅ |
| GET | `/api/v1/tasks/:taskId/activity` | View a task's activity history | ✅ |
| GET/POST/DELETE | `/api/v1/tasks/:taskId/labels/:labelId` | Attach / detach / list a task's labels | ✅ |
| GET/POST | `/api/v1/projects` | List / create projects | ✅ |
| GET/PUT/DELETE | `/api/v1/projects/:id` | Get / update / delete a project | ✅ |
| GET | `/api/v1/projects/:id/tasks` | List a project's tasks | ✅ |
| GET/POST | `/api/v1/projects/:id/members` | List / add project members | ✅ |
| DELETE | `/api/v1/projects/:id/members/:userId` | Remove a member | ✅ |
| GET/POST | `/api/v1/projects/:id/epics` | List / create epics | ✅ |
| PUT/DELETE | `/api/v1/epics/:epicId` | Update / delete an epic | ✅ |
| GET/POST | `/api/v1/projects/:id/milestones` | List / create milestones | ✅ |
| PUT/DELETE | `/api/v1/milestones/:milestoneId` | Update / delete a milestone | ✅ |
| GET/POST | `/api/v1/projects/:id/sprints` | List / create sprints | ✅ |
| PUT/DELETE | `/api/v1/sprints/:sprintId` | Update (incl. status) / delete a sprint | ✅ |
| GET/POST | `/api/v1/projects/:id/labels` | List / create labels | ✅ |
| DELETE | `/api/v1/labels/:labelId` | Delete a label | ✅ |
| GET | `/health` | Health check (includes DB connectivity) | ❌ |

Full request/response schemas are available in Swagger (see below).

## API Documentation

Interactive API documentation (Swagger UI) is available once the backend is running:

http://localhost:3000/swagger/index.html


To test authenticated endpoints: call `POST /auth/login` to obtain an access token, then click **Authorize** in the Swagger UI and enter `Bearer <access-token>`.

## CI/CD

Every push to `main` triggers a GitHub Actions pipeline that:
1. Builds the backend and runs `go vet`
2. Scans dependencies for known vulnerabilities (`govulncheck`)
3. Runs the full test suite with the race detector enabled
4. Builds both Docker images to catch containerization issues early

See [`.github/workflows/backend-ci.yml`](.github/workflows/backend-ci.yml).

## Database Backup & Restore

The production database is backed up automatically every day at **01:00 (Vietnam time / UTC+7)** via GitHub Actions.

- Workflow: [`.github/workflows/db-backup.yml`](.github/workflows/db-backup.yml)
- Backups are stored as workflow **artifacts**, retained for **30 days**
- Can also be triggered manually: **Actions** tab → **Database Backup** → **Run workflow**

> Render Postgres also provides its own automatic backups (see the **Backups** tab in the Render Dashboard). This workflow is a **second layer of protection** — e.g. in case of lost access to Render, or to keep a copy outside Render's infrastructure — not a replacement for it.

### Downloading a backup

1. Go to the **Actions** tab → select the desired **Database Backup** run
2. Scroll to **Artifacts** → download `db-backup-<run_id>`
3. Unzip to get the `.dump` file

### Restoring

Requires `pg_restore` (version equal to or newer than the target server's Postgres version).

```powershell
pg_restore --clean --if-exists -d "postgresql://user:pass@host/dbname" backup_20260817_010000.dump
```

- `-d "..."` — connection string of the **target** database to restore into (from Render Dashboard → Connections → **External Database URL**)
- `--clean --if-exists` — drops existing objects before restoring, avoiding conflicts
- Replace `backup_20260817_010000.dump` with the actual downloaded filename

### Required secret

The backup workflow reads `PROD_DATABASE_URL` from **Settings → Secrets and variables → Actions**. This must be Render's **External Database URL** (not Internal), since GitHub Actions runs outside Render's network.

## Notable Design Decisions

- **Two-tier authorization, enforced in the service layer** — every project has exactly one `owner` and any number of `member`s. Owners can manage membership and delete/update project-level entities (epics, milestones, sprints, labels); members can create and collaborate but only modify tasks they created or were assigned to. This mirrors real tools like Jira/Linear rather than a flat "everyone can do everything" model.
- **Rotating refresh tokens** — access tokens are short-lived (15 min); refresh tokens are long-lived (7 days), stored **hashed** (never plaintext), and **rotated** on every use — reusing an old refresh token is rejected, limiting the blast radius of a leaked token.
- **Cross-entity project consistency checks** — attaching an Epic, Milestone, Sprint, or Label to a Task explicitly verifies they belong to the *same* project, preventing accidental (or malicious) cross-project data mixing.
- **Dependency injection via interfaces** — repositories are injected as interfaces into services, enabling unit tests with mocked dependencies instead of a live database.
- **Defense in depth for data integrity** — task/sprint status is validated both in application code and via PostgreSQL `CHECK` constraints.
- **Structured, catalogued logging** — instead of ad-hoc `log.Println`, every significant event (success and failure paths) has a named function in `internal/logger` (e.g. `TaskCreated`, `AuthLoginInvalidPassword`), producing consistent JSON fields and making the full set of loggable events discoverable by reading the logger package.
- **Graceful shutdown** — on `SIGINT`/`SIGTERM`, the server stops accepting new requests but lets in-flight ones finish (bounded by a timeout) before closing the database connection — important for zero-downtime deploys on platforms like Render.
- **Separate rate-limit configuration for tests** — rate limiting is toggled via a parameter, avoiding flaky tests caused by shared IP-based limits in the test environment.

## License

This project is for portfolio purposes.