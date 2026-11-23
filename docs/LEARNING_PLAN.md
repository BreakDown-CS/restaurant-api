# Learning Plan

A restaurant system (buffet + à la carte, POS, QR ordering) rebuilt from my graduation project, used as a hands-on path to production-style full-stack development.

- **Backend:** [restaurant-api](https://github.com/BreakDown-CS/restaurant-api) — Go, Gin, PostgreSQL 18, pgx, goose
- **Frontend:** [restaurant-admin](https://github.com/BreakDown-CS/restaurant-admin) — React 19, TypeScript, TanStack Router/Query, Ant Design, zod
- **How I work:** I type every line myself; an AI tutor explains concepts, reviews my code and checks each step. Every step ends with lint + tests and a Conventional Commit. Since 23 Nov 2026 all work goes through a branch, a pull request and CI; `main` is protected in both repos.

Last updated: 23 Nov 2026

---

## Done

### Frontend foundations — `restaurant-admin` (13–17 Nov 2026)

- [x] **Phase 0 — Project cleanup**
  - Learn: path aliases, TypeScript scripts, Prettier, branch naming
  - Practice: `@/*` alias, `typecheck` script, format setup
- [x] **Phase 1 — Feature-based structure**
  - Learn: organising code by feature, public `index.ts` per feature, LF line endings
  - Practice: `features/dashboard`, `.gitattributes`, format on save
- [x] **Phase 2 — Validated environment variables**
  - Learn: why config is checked at startup, why `VITE_` values must not hold secrets
  - Practice: `src/lib/env.ts` with zod, `.env.example`
- [x] **Phase 3 — Data fetching**
  - Learn: TanStack Query cache and stale time, one API client for every request, typed errors, mocking an API
  - Practice: QueryClient in router context, `api-client.ts` + `ApiError`, menu list from MSW in an Ant Design table
- [x] **Phase 4 — Forms and CRUD**
  - Learn: React Hook Form + zod resolver, mutations and cache invalidation
  - Practice: create / edit / delete menu with a modal form
- [x] **Phase 5 — Testing**
  - Learn: Vitest, Testing Library, MSW in tests, testing behaviour instead of implementation
  - Practice: tests for the API client, schemas and the menu form
- [x] **Phase 6 — Git hooks**
  - Learn: catching problems before a commit
  - Practice: lefthook (Prettier on staged files) + commitlint
- [x] **Phase 7 — CI**
  - Learn: GitHub Actions, running the same checks on a clean machine
  - Practice: typecheck, lint, format check, tests and build on every push

### Backend — `restaurant-api` (18–20 Nov 2026)

- [x] **B0 — Project setup**
  - Learn: Go modules, config validation, Docker Compose, connection pools, migrations, graceful shutdown, structured logging, linting
  - Practice: Gin router + `GET /api/v1/health` (checks the DB), `.env` config, PostgreSQL 18 in Docker, pgx pool, goose migrations embedded in the binary, `http.Server` with signal handling, slog, golangci-lint v2, Makefile, live reload with air
- [x] **B1 — Menu CRUD**
  - Learn: Handler → Usecase → Repository layering, money as decimal (never float), UUID v7, validation in two layers, sentinel errors, CORS, interfaces declared by the consumer, fakes in tests, `httptest`
  - Practice: `GET/POST/PUT/DELETE /api/v1/menus`, one response envelope `{ success, message, data, meta, error }`, 400 / 404 / 500 mapping that never leaks internal errors, CORS allow-list from env, 12 tests (usecase rules with a fake repository, HTTP handlers end to end)
- [x] **B4 — Backend CI**
  - Learn: CI does not block a push by itself; reading a failed CI log; pull requests
  - Practice: GitHub Actions running vet, golangci-lint and tests; broke a test on purpose in a PR and read the failure

### Connecting both sides (20–23 Nov 2026)

- [x] **Repo readiness**
  - Learn: what a README should show, branch protection
  - Practice: README for both repos, rulesets (no direct push to `main`, PR + passing `check` required, no force push or deletion), branch → PR → squash merge workflow
- [x] **F1 — Frontend talks to the real backend**
  - Learn: keeping a frontend/backend contract, parsing error envelopes, keeping mocks identical to the real API, matching validation rules on both sides
  - Practice: API client unwraps `data` and turns error envelopes into `ApiError { code, details }`; MSW handlers return the same envelope; price limited to 2 decimals and names trimmed like the backend; full CRUD from the browser verified in PostgreSQL

---

## Next

- [ ] **B2 — Categories**
  - Learn: table relations, foreign keys, migrating existing data
  - Practice: `categories` table, menus reference a category, category dropdown in the admin
- [ ] **BF1 — Buffet packages**
  - Learn: many-to-many relations
  - Practice: `buffet_packages` (price per person, time limit, height rule, overtime and waste fees) + `package_menus`, management pages
- [ ] **BF2 — Tables, orders and bills**
  - Learn: transactions, price snapshots, deriving state instead of storing it twice
  - Practice: open a table with guest count, buffet and à la carte bills, bill lines by type (buffet / extra / overtime / waste), table status computed from open bills, order history
- [ ] **P1 — POS checkout**
  - Learn: money calculations, closing a bill safely
  - Practice: close bill, cash with change, static PromptPay QR, receipt, daily summary
- [ ] **B3 — Login and roles**
  - Learn: authentication with JWT, role-based authorization
  - Practice: staff login, 5 fixed roles (kitchen, waiter, cashier, manager, admin), money-sensitive actions limited to manager and admin
- [ ] **BF3 — Customer QR ordering**
  - Learn: public API design, short-lived tokens, mobile-first UI
  - Practice: customer API using a per-bill token, new frontend repo `restaurant-order`
- [ ] **K1 — Kitchen display**
  - Practice: order queue with status (waiting / cooking / served)
- [ ] **M1 — Discounts and members**
  - Practice: apply discounts and member benefits at checkout
- [ ] **R1 — Sales reports**
  - Practice: reports built on the daily POS summary
- [ ] **D1 — Docker and deployment**
  - Learn: production images, environment configuration, hosting
  - Practice: containerise both apps and deploy a public demo
- [ ] **S1 — Ingredient stock**
  - Practice: recipe → ingredients, deduct stock when ordering, low-stock alerts
- [ ] **Frontend Phase 8 — Finishing touches**
  - Practice: `.nvmrc`, `engines`, error boundary
