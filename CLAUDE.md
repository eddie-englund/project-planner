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

### Json

Always defined camelCase json definitions like

```go
type createProjectRequest struct {
	Title    string  `json:"title"     validate:"required"`
	Color    string  `json:"color"     validate:"required"`
	ImageURL *string `json:"image_url" validate:"omitempty"`
}
```

### Database

PostgreSQL via pgx/v5. Connection credentials from environment (`.env`). Docker Compose config is at `docker/docker-compose.yaml` (postgres:18, db `app`, user/pass `postgres`).

### sqlc workflow

Schema is read from `migrations/`, queries from `internal/db/queries/`. After changing either, run `sqlc generate` to regenerate `internal/db/generated/`. The generated `Queries` struct wraps a `DBTX` interface (pgx-compatible), so it can be used with both a pool and a transaction.

### Data model

- `users` — id (UUID), email, email_verified, created_at
- `projects` — id, title, color, image_url, created_by (FK → users)
- `project_topics` — id, project_id (FK → projects), index, title, color, image_url

## Frontend (`frontend/`)

Vue 3 + Vite + TypeScript + Pinia + Vue Router 5 (file-based routing) + Tailwindcss.

Commands run from `frontend/`:

```bash
pnpm dev          # dev server
pnpm build        # type-check + build
pnpm test:unit    # vitest
pnpm test:e2e     # playwright
```

### File-based routing

Routes live in `src/pages/`. Convention:

| File                        | Route        |
| --------------------------- | ------------ |
| `src/pages/index.vue`       | `/`          |
| `src/pages/about.vue`       | `/about`     |
| `src/pages/users/[id].vue`  | `/users/:id` |
| `src/pages/users/index.vue` | `/users`     |

`typed-router.d.ts` is auto-generated on `pnpm dev` — do not edit manually. Router config is `src/router/index.ts`; routes import from `vue-router/auto-routes`.

**Best practices:**

- All new pages go in `src/pages/` — never add manual routes to `router/index.ts`
- Use `<RouterView />` in layout components, not hardcoded page content
- Nested routes → nested folders (e.g. `src/pages/projects/[id]/topics.vue`)
- Use `definePage()` macro in SFCs to set route meta (title, requiresAuth, etc.)

### Vue 3 best practices

- Always use `<script setup lang="ts">` — no Options API
- Prefer `defineProps` / `defineEmits` with TypeScript generics, not runtime declarations
- Extract reusable logic into composables in `src/composables/` (`use` prefix)
- Use `computed()` for derived state; avoid watchers unless reacting to external side effects
- Keep components small and focused — one responsibility per component
- Prefer `v-bind` shorthand (`:prop`) and `v-on` shorthand (`@event`)
- Before writing a new `<button>`, use `AppButton` (`src/components/AppButton.vue`) with `variant="primary|ghost|secondary|outline"`. Only use raw `<button>` for highly specific UI with no shared pattern (color pickers, card tiles, icon-only utilities).
- If the same markup appears in 2+ places, extract a component. Prefer extending an existing component over creating a new one.

### Pinia best practices

- Stores live in `src/stores/` — one store per domain (e.g. `useProjectsStore`, `useTopicsStore`)
- Use `defineStore` with the setup syntax (preferred over options syntax):
  ```ts
  export const useProjectsStore = defineStore('projects', () => {
    const items = ref<Project[]>([])
    // actions are plain functions
    async function fetchAll() { ... }
    return { items, fetchAll }
  })
  ```
- Store IDs must be unique strings matching the variable name
- Keep API calls inside store actions, not in components
- Use `storeToRefs()` when destructuring reactive state from a store

### API calls

**Always use `useApi` from `src/composables/useApi.ts`** for all HTTP requests unless the user specifies otherwise. It is a `createFetch` instance from `@vueuse/core` pre-configured with `VITE_API_BASE_URL`.

```ts
import { useApi } from "@/composables/useApi";

const { data, error, isFetching } = useApi("/projects").json<Project[]>();
```

- Base URL set via `VITE_API_BASE_URL` in `.env.local` (default: `http://localhost:8080`)
- Never use raw `fetch` or `axios` — always go through `useApi`
- Invoke inside store actions, not directly in components

### UI changes

**Always use tailwindcss unless explicitly allowed by the user**
**Always invoke the `frontend-design` skill** before implementing any UI. Run `/frontend-design` and describe the feature — do not make visual/layout decisions without it.
