# Calculator

Full-stack calculator application — React (TypeScript) frontend + Go REST API backend.

## Architecture

```
app-web/          React SPA (Vite + Ant Design)
services/api/     Go REST API (Chi router)
Dockerfile        Multi-stage build (Go + Node + Caddy)
```

**Backend** — A single `POST /calculate` endpoint handles all operations. The core calculation logic lives in `pkg/calculator/` as a pure function with no dependencies, making it easy to test. The HTTP handler in `pkg/controller/calculator/` deals only with request parsing, error mapping, and JSON responses. Structured logging via `slog` is the only injected dependency.

**Frontend** — A `useCalculator` hook manages all calculator state (operands, pending operation, display, history). It calls the backend API for every calculation so results are always server-authoritative. The UI is split into four components: `Calculator` (container), `Display`, `Keypad`, and `HistoryPanel`.

### Design decisions

- **Single endpoint** instead of one-per-operation (`/add`, `/subtract`, etc.) — fewer routes, one request/response contract, and the operation is just data in the payload. Adding a new operation means adding a case to the switch, not wiring up a new route.
- **Vite over Next.js** — no SSR needed for a calculator. Vite gives a faster dev loop and simpler config.
- **Ant Design** — provides accessible, responsive components out of the box so I could focus on logic instead of styling from scratch.
- **Caddy in Docker** — serves the SPA static files and reverse-proxies `/api/*` to the Go binary, all in a single container. No nginx config files to maintain.
- **Server-side calculation** — the frontend delegates all math to the API. This keeps the frontend thin and ensures a single source of truth for how operations behave.

## Prerequisites

- Go 1.26+
- Node.js 22+
- Or just Docker

## Makefile

Common tasks are available via `make`:

| Target | Description |
|--------|-------------|
| `make dev-api` | Start the Go API on `:4000` |
| `make dev-web` | Start the React dev server on `:5173` |
| `make test-api` | Run Go tests |
| `make test-web` | Run frontend tests |
| `make test` | Run all tests |
| `make coverage-api` | Go tests with per-function coverage report |
| `make coverage-web` | Frontend tests with coverage report |
| `make coverage` | Both coverage reports |
| `make vet` | Run `go vet` |
| `make lint` | Run `golangci-lint` |
| `make check` | Run vet + lint + all tests |
| `make build` | `docker compose build` |
| `make up` | `docker compose up` |
| `make down` | `docker compose down` |

## Running locally

**Backend** (runs on `:4000`):

```bash
cd services/api
go run ./cmd
```

**Frontend** (runs on `:5173`, proxies API calls to `:4000`):

```bash
cd app-web
npm install
npm run dev
```

Open http://localhost:5173.

**With Docker** (runs everything on `:8080`):

```bash
docker compose up --build
```

Open http://localhost:8080.

## Running tests

```bash
# Go tests
cd services/api && go test ./...

# Frontend tests
cd app-web && npm test
```

### Coverage

```bash
make coverage        # both layers
make coverage-api    # Go: writes services/api/coverage.out + prints per-function summary
make coverage-web    # frontend: vitest --coverage (v8), prints per-file table
```

For an HTML view of the Go report: `cd services/api && go tool cover -html=coverage.out`.

Current numbers:

| Layer | Coverage | Notes |
|-------|----------|-------|
| Go — `pkg/calculator` | 100% | All operations, error paths, and expression formatting |
| Go — `pkg/controller/calculator` | 96% | HTTP handler; uncovered lines are the JSON-encode failure branch |
| Go — total | 55% | Remainder is wiring (`main`, router, config, logging middleware) with no branching logic |
| Frontend — `useCalculator` hook | 100% | All state logic: digit entry, chaining, unary ops, errors, clear/backspace, result formatting |
| Frontend — components | 74% | `Display`, `Keypad`, `HistoryPanel` fully covered; `Calculator` container is wiring only |
| Frontend — total | 74% | Remaining gap is the API client's fetch wrapper and app-mount wiring |

## API

### `POST /calculate`

Performs a calculation and returns the result with a formatted expression.

**Request body:**

| Field       | Type   | Required | Description                                                                 |
|-------------|--------|----------|-----------------------------------------------------------------------------|
| `operation` | string | yes      | One of: `add`, `subtract`, `multiply`, `divide`, `power`, `sqrt`, `percentage` |
| `a`         | number | yes      | First operand                                                               |
| `b`         | number | for binary ops | Second operand (not needed for `sqrt`)                                |

**Response (200):**

```json
{ "result": 8, "expression": "5 + 3" }
```

**Error responses:**

| Status | When                              |
|--------|-----------------------------------|
| 400    | Invalid JSON, unknown operation, missing operand |
| 422    | Division by zero, square root of negative number |

### Examples

```bash
# Addition
curl -X POST http://localhost:4000/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "add", "a": 5, "b": 3}'
# {"result":8,"expression":"5 + 3"}

# Division
curl -X POST http://localhost:4000/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 10, "b": 4}'
# {"result":2.5,"expression":"10 / 4"}

# Square root
curl -X POST http://localhost:4000/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "sqrt", "a": 16}'
# {"result":4,"expression":"√16"}

# Percentage (what is 15% of 200?)
curl -X POST http://localhost:4000/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "percentage", "a": 15, "b": 200}'
# {"result":30,"expression":"15% of 200"}

# Division by zero → 422
curl -X POST http://localhost:4000/calculate \
  -H "Content-Type: application/json" \
  -d '{"operation": "divide", "a": 1, "b": 0}'
# {"error":"division by zero"}
```

## AI tooling

This project was built with [Claude Code](https://docs.anthropic.com/en/docs/claude-code) (Claude Opus). I used it for scaffolding components, writing tests, debugging edge cases, and iterating on the API design. All code was reviewed and understood before committing.

The prompts used are collected in [PROMPTS.md](PROMPTS.md).
