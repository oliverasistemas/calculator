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
- Docker

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

## Edge cases covered

All of these are exercised by tests (`services/api/pkg/calculator/calculator_test.go`, `app-web/src/hooks/useCalculator.test.ts`).

### Numeric precision (backend)

Add, subtract, multiply, divide, and percentage are computed in **exact rational arithmetic** (`big.Rat`) — floats only enter when the result is formatted:

- **No binary-float artifacts** — `0.1 + 0.2` returns `0.3`, not `0.30000000000000004`; `0.3 - 0.1` returns `0.2`.
- **Catastrophic cancellation** — `1.0000000001 - 1` returns exactly `1e-10`.
- **Operands beyond 2⁵³ survive** — operands are decoded from their literal JSON digits (never through `float64`), so `9007199254740993` stays `9007199254740993` instead of silently becoming `…992`.
- **Exact integer powers** — integer exponents use `big.Int`, so `3 ^ 35 = 50031545098999707` and `2 ^ 100` are digit-for-digit exact where `math.Pow` would be off. Exponents over 4096 or bases over 256 bits fall back to floats to keep hostile requests from forcing huge allocations.
- **Honest rounding** — non-terminating decimals round to 12 significant digits (`1 / 3 → 0.333333333333`); exact results longer than 50 digits round instead of printing rounding artifacts (`10 ^ 60 → 1e+60`).
- **`resultText` is authoritative** — the response carries the result as a decimal string alongside the compatibility `result` float, because a `float64` cannot represent integers beyond 2⁵³. The frontend displays and chains from `resultText` verbatim.

### Domain errors (422)

- Division by zero, including `0` raised to a negative power.
- Square root of a negative number (`√0` is fine; `0 ^ 0` returns `1`).
- Non-finite results are rejected rather than breaking JSON encoding: `(-1) ^ 0.5` (NaN), `10 ^ 1000` and `1e308 × 10` (overflow to ±Inf).

### Input validation (400)

- Unknown operations and missing operands (`a` always, `b` for binary ops).
- Malformed numbers rejected: `abc`, `1.2.3`, `1e`, `0x10`, `Inf`, `NaN`, `--5`, empty strings.
- Calculator-display forms accepted and normalized: `5.`, `.5`, leading `+`.
- Operand literals capped at 512 characters so hostile requests can't feed arbitrarily large rationals into the exact-arithmetic paths.

### UI state (frontend)

- Operands are kept and sent as **strings** — no `parseFloat` on the client, so chained results beyond 2⁵³ don't get corrupted between steps.
- A second decimal point in the same number is ignored (`1.5.2` → `1.52`).
- The leading `0` is replaced by the first digit typed; after a result, typing starts a fresh number.
- Chaining: pressing a second operator (`5 + 3 ×`) computes the pending operation and carries the result forward; re-selecting an operator before typing just replaces the pending one without firing a request.
- `=` with no pending operation is a no-op.
- On API errors the display and history are kept, the error message is shown, and typing any digit clears it; network failures fall back to a generic message.
- Backspace bottoms out at `0` and collapses a single negated digit (`-5` → `0`); toggle-sign leaves `0` alone.
- `C` clears the current entry but keeps history; `AC` wipes both. History is capped at 50 entries.

## AI tooling

This project was built with [Claude Code](https://docs.anthropic.com/en/docs/claude-code) (Claude Opus). I used it for scaffolding components, writing tests, debugging edge cases, and iterating on the API design. All code was reviewed and understood before committing.

The prompts used are collected in [PROMPTS.md](PROMPTS.md).
