# Restaurant API

![CI](https://github.com/BreakDown-CS/restaurant-api/actions/workflows/ci.yml/badge.svg)

Rebuilt from my graduation project to learn how a production-style Go backend is designed, not just how to make it work.

Frontend: [restaurant-admin](https://github.com/BreakDown-CS/restaurant-admin)

## Tech stack

Go 1.27 · Gin · PostgreSQL 18 · pgx · goose (migrations) · Docker Compose · testify · GitHub Actions · golangci-lint

## Architecture

Feature-based clean architecture. Each feature owns its handler, usecase and repository.

```
HTTP ─► handler ─► usecase ─► repository ─► PostgreSQL
        (JSON,      (business   (SQL only)
         status)     rules)
```

- The usecase declares the repository **interface it needs**, so tests swap in a fake without a database.
- Every response uses one envelope: `{ success, message, data, error: { code, details } }`.
- Unexpected errors are logged on the server and returned as `500 INTERNAL_ERROR` without internal details.

## Features

- [x] Menus CRUD with validation (price ≥ 0, max 2 decimals)
- [x] CI: vet, lint and tests on every push and pull request
- [ ] Tables and bills
- [ ] Buffet packages + à la carte add-ons
- [ ] POS checkout
- [ ] Login with 5 roles (admin, manager, cashier, waiter, kitchen)
- [ ] Customer QR ordering
- [ ] Kitchen display, members & discounts, sales reports, stock

## Getting started

Requires Go 1.27, Docker and GNU Make.

```bash
cp .env.example .env
make dev        # starts PostgreSQL in Docker, runs migrations, live reload on :8080
```

```bash
curl http://localhost:8080/api/v1/menus
```

## API

| Method | Path | Success | Errors |
| --- | --- | --- | --- |
| GET | `/api/v1/menus` | 200 | |
| POST | `/api/v1/menus` | 201 | 400 |
| PUT | `/api/v1/menus/:id` | 200 | 400, 404 |
| DELETE | `/api/v1/menus/:id` | 204 | 404 |

## Testing

```bash
make check      # go vet + golangci-lint + go test
```

Usecase rules are tested with a fake repository; HTTP handlers are tested end to end with `httptest`.

## What I learned

- Why the usecase declares the repository interface itself: it lets tests use a fake instead of a real database.
- Never send internal errors to the client: log the details, return a generic 500.
- CI does not block a push by itself; it needs pull requests and branch rules to protect main.