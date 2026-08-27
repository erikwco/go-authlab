# go-authlab

Nota: esto es un laboratorio de aprendizaje, no un sistema de autenticación en producción.

A public **learning lab** for an Auth-shaped HTTP service in Go + Postgres. It is an original portfolio project by Erik Chacón (Software Engineer — Auth at Supabase interview). The product outline is inspired by GoTrue (signup, sessions, JWT, OAuth later), but this is **not** a GoTrue fork, **not** affiliated with Supabase, and **not production Auth experience**.

Do not use this to protect real users or real credentials.

## Honesty

This repo does **not** implement authentication. There is no signup, login, password hashing, JWT, sessions, refresh tokens, or OAuth. Milestone 1 is a runnable skeleton: process, Postgres, migrations, and a health check.

## What this milestone is

- A Go module (`github.com/erikwco/go-authlab`) that starts an HTTP server.
- Docker Compose: Postgres 16 + the app (app waits for a healthy DB, runs migrations, then serves HTTP on `:8080`).
- One SQL migration that creates `users` (empty table; unused until signup).
- `GET /health` that pings the database.
- Unit tests, Makefile targets, a multi-stage Dockerfile, and GitHub Actions (`go test` + `go vet`).

## What this milestone is not

- Not a production auth service.
- Not signup / login / JWT / OAuth / sessions (those are later milestones).
- Not Kubernetes, not Terraform, not a hosted IdP.

## Run locally

Prerequisites: Docker and Docker Compose.

```bash
cp .env.example .env   # dummy local values only; do not commit secrets
make up
```

In another terminal:

```bash
curl -sS http://localhost:8080/health
# {"status":"ok","db":"ok"}
```

Stop:

```bash
make down
```

`make up` builds the app image, starts Postgres (named volume `authlab_postgres_data`, healthcheck via `pg_isready`), waits until Postgres is healthy, then starts the app. On boot the app runs embedded SQL migrations and listens on `:8080`.

### Without the app container

If `DATABASE_URL` is set (Postgres already running, e.g. `docker compose up -d postgres`):

```bash
export DATABASE_URL=postgres://authlab:authlab@localhost:5432/authlab?sslmode=disable
make run
```

`HTTP_ADDR` defaults to `:8080`. `DATABASE_URL` is required.

### Tests

```bash
make test
```

## Architecture

HTTP handlers receive requests, call a thin service-style dependency (today: a DB ping), and Postgres is the only datastore. There is no extra framework: stdlib `net/http`, a small `internal/db` pool, and SQL files under `migrations/`. Later auth logic should stay out of handlers and talk to Postgres through that same layering: **HTTP → service → postgres**.

## Stack

- Go 1.22+
- Postgres 16
- Docker Compose

No Terraform. No Kubernetes.

## Config

| Variable       | Required | Default | Notes                                      |
|----------------|----------|---------|--------------------------------------------|
| `HTTP_ADDR`    | no       | `:8080` | Listen address                             |
| `DATABASE_URL` | yes      | —       | Postgres URL (`postgres://…`)              |

See `.env.example` for dummy local values. Never commit `.env` or real secrets.

## Migrations

A tiny embedded migrator in `internal/db` applies `migrations/*.up.sql` at startup (sorted by filename). Applied versions are stored in `schema_migrations`. It is **not** golang-migrate: statements are split on `;`, so keep migrations simple (no semicolons inside string literals). Down files are present for reference but are not applied automatically.

The first migration creates `users` with `id`, `email` (CITEXT unique), `password_hash` (unused until signup / hito 2), and `created_at`.

## Layout

```
cmd/authlab/main.go    # process entry: config, migrate, HTTP server
internal/http/         # handlers, health, timeouts
internal/db/           # pool, ping, migrate
migrations/            # SQL
```

## Next milestones (not in this PR)

These are future work, listed so the lab has a product shape. They are **not** implemented here:

1. Signup / login (passwords, `users.password_hash`)
2. Sessions + JWT
3. Refresh token rotation
4. OAuth2 authorization code + PKCE

## License

Use this as a learning project. It is not a security product.
