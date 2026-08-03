# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Build & Run Commands

| Task | Command |
|------|---------|
| Run Go API (`:4000`) | `make dev-api` |
| Run React dev server (`:5173`) | `make dev-web` |
| Run all tests | `make test` |
| Run Go tests only | `make test-api` |
| Run frontend tests only | `make test-web` |
| Run a single Go test | `cd services/api && go test ./pkg/calculator/ -run TestName` |
| Run a single frontend test | `cd app-web && npx vitest run src/components/Display.test.tsx` |
| Frontend tests in watch mode | `cd app-web && npm run test:watch` |
| Vet + lint + all tests | `make check` |
| Docker build & run (`:8080`) | `make up` |

## Architecture

Monorepo with two apps: a Go REST API (`services/api/`) and a React SPA (`app-web/`).

**Backend** (`services/api/`): Single `POST /calculate` endpoint. Chi router, structured logging via `slog`.
- `cmd/main.go` — entrypoint, wires config + logger + router
- `pkg/calculator/` — pure calculation logic (no dependencies), all math lives here
- `pkg/controller/calculator/` — HTTP handler: request parsing, error mapping, JSON responses
- `pkg/router/` — route registration and middleware
- `pkg/types/` — shared request/response structs
- `pkg/logger/` — slog setup
- `pkg/utils/` — number formatting helpers

**Frontend** (`app-web/`): Vite + React 19 + TypeScript + Ant Design.
- `src/hooks/useCalculator.ts` — all calculator state management; calls backend API for every calculation (server-authoritative)
- `src/components/` — `Calculator` (container), `Display`, `Keypad`, `HistoryPanel`
- `src/api/` — API client for the backend
- `@` path alias maps to `./src/` (configured in `vite.config.ts`)

**Key design choice**: The frontend delegates all math to the API — there is no client-side calculation. Adding a new operation means adding a case in `pkg/calculator/` and updating the frontend UI.

## Testing

- **Go**: Standard `go test` with table-driven tests. Tests for calculator logic and HTTP handler.
- **Frontend**: Vitest + React Testing Library + jsdom. Setup file at `src/test/setup.ts`. Tests are colocated with components.

## Docker

Single-container deployment: multi-stage Dockerfile builds both Go binary and React static assets, serves via Caddy (reverse-proxies `/api/*` to the Go binary, serves SPA for everything else).
