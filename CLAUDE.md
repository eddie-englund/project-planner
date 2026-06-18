# Project planner/Kanban application.

Backend: Go + PostgreSQL + sqlc
Frontend: Vue 3 + TypeScript + Pinia + Vue Router (file-based routing) + TailwindCSS

Projects contain topics, which contain tasks, notes, and files.

## General

Always use the skill caveman ultra.

## Commands

Run from backend/:

go run ./cmd/api/main.go
go build ./...
go test ./...
go test ./internal/api/...
sqlc generate
docker compose -f ../docker/docker-compose.yaml up -d
Frontend

Run from frontend/:

pnpm dev
pnpm build
pnpm test:unit
pnpm test:e2e

## Architecture

### Backend

Entry point: cmd/api/main.go
Routes: internal/api/
SQL queries: internal/db/queries/
Generated sqlc code: internal/db/generated/ (never edit manually)
Schema source: migrations/

After changing migrations or SQL queries:

sqlc generate

Never manually edit files in:

internal/db/generated/
API

### Frontend

Pages: src/pages/
Stores: src/stores/
Composables: src/composables/
Components: src/components/
Rules
Database

Always use:

useApi()

for HTTP requests.

Do not use raw fetch or axios.

UI
Use TailwindCSS.
Use existing shared components before creating new ones.
Prefer AppButton over raw <button> where applicable.
Invoke the frontend-design skill before implementing UI changes.
Invoke chromium mcp server to check the changes you're making and ensure they work as intended.
Code Style
Vue: <script setup lang="ts">
Prefer composables for reusable logic.
Prefer computed values over watchers.
Keep components focused and reusable.

## Tests

New functionality must be covered by tests.

Backend (Go): add handler unit tests in internal/api/*_test.go using mock store interfaces (see store_interfaces.go). Run: go test ./internal/api/...

Frontend unit (Vitest): add store tests in src/stores/__tests__/*.spec.ts. Mock useApi via vi.mock('@/composables/useApi'). Run: pnpm test:unit

Frontend e2e (Playwright): add page-flow tests in e2e/*.spec.ts using page.route() to mock the API (base URL: http://localhost:8080). Run: pnpm test:e2e

Run all three suites before marking work complete.
