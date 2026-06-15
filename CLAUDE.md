# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project overview

A project planner app: a landing page of color-coded tiles for projects, each project containing topics, and each topic containing tasks/notes/files. Backend is Go, frontend is TBD (currently empty).

## Commands

All commands run from `backend/`:

```bash
# Run the server
go run ./cmd/api/main.go

# Build
go build ./...

# Test
go test ./...

# Run a single test
go test ./internal/api/...

# Regenerate sqlc types after editing queries or migrations
sqlc generate

# Start the database
docker compose -f ../docker/docker-compose.yaml up -d
```

## Architecture

### Backend (`backend/`)

- `cmd/api/main.go` — entrypoint; loads `.env`, constructs the router, serves on `:8080`
- `internal/api/router.go` — registers HTTP routes using stdlib `net/http`
- `internal/db/queries/*.sql` — SQL queries annotated for sqlc
- `internal/db/generated/` — sqlc-generated Go code (do not edit manually)
- `migrations/` — raw SQL migration files; also serve as the schema source for sqlc (`sqlc.yaml` points `schema` here)

### Database

PostgreSQL via pgx/v5. Connection credentials from environment (`.env`). Docker Compose config is at `docker/docker-compose.yaml` (postgres:18, db `app`, user/pass `postgres`).

### sqlc workflow

Schema is read from `migrations/`, queries from `internal/db/queries/`. After changing either, run `sqlc generate` to regenerate `internal/db/generated/`. The generated `Queries` struct wraps a `DBTX` interface (pgx-compatible), so it can be used with both a pool and a transaction.

### Data model

- `users` — id (UUID), email, email_verified, created_at
- `projects` — id, title, color, image_url, created_by (FK → users)
- `project_topics` — id, project_id (FK → projects), index, title, color, image_url
