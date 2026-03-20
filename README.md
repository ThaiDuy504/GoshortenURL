# snip — URL Shortener Service

> A fast, self-hosted URL shortener built with Go, backed by PostgreSQL and Redis.

![Go](https://img.shields.io/badge/Go-1.25-00ADD8?logo=go&logoColor=white)
![Gin](https://img.shields.io/badge/Gin-v1.12-brightgreen)
![PostgreSQL](https://img.shields.io/badge/PostgreSQL-16-4169E1?logo=postgresql&logoColor=white)
![Redis](https://img.shields.io/badge/Redis-latest-DC382D?logo=redis&logoColor=white)
![Docker](https://img.shields.io/badge/Docker-compose-2496ED?logo=docker&logoColor=white)

---

## Features

- Shorten any URL to a 6-character base-62 code
- Redirect with click-count tracking (async, non-blocking)
- Redis cache with 24-hour TTL — zero DB hits on warm cache
- PostgreSQL persistence with auto schema migration on startup
- Health check endpoint
- Multi-stage Docker build (~15 MB final image)

---

## Tech Stack

| Layer | Technology |
|---|---|
| Language | Go 1.25 |
| HTTP Framework | Gin v1.12 |
| Database | PostgreSQL 16 (pgx/v5) |
| Cache | Redis (go-redis/v9) |
| Config | cleanenv |
| Containerization | Docker + Compose |

---

## Architecture

```
HTTP Request
     │
     ▼
┌─────────────┐
│   Handler   │  ← Gin routes, HTML templates
└──────┬──────┘
       │
       ▼
┌─────────────┐
│   Service   │  ← Business logic, short code generation
└──────┬──────┘
       │
       ▼
┌─────────────────────────────┐
│        Repository           │
│                             │
│  Redis ──► cache hit        │
│    │                        │
│    └──► miss ──► Postgres   │
└─────────────────────────────┘
```

**Request flow (redirect):**
1. Look up `shortCode` in Redis
2. Cache hit → fire `incrementClickCount` in a goroutine → redirect immediately
3. Cache miss → query Postgres → populate Redis → goroutine click count → redirect

---

## Project Structure

```
Go_shortenURL/
├── cmd/api/
│   └── main.go              # Entrypoint: wire dependencies, register routes
├── internal/
│   ├── handler/             # HTTP layer: parse requests, render HTML
│   ├── service/             # Business logic: shorten & resolve URLs
│   └── repository/          # Data layer: Redis cache + PostgreSQL
├── pkg/shortener/           # Short code generator (base-62, 6 chars)
├── configs/
│   ├── config.go            # Config struct, loaded via cleanenv
│   └── config.example.env   # Environment variable template
├── db/migrations/           # SQL migration files
├── templates/               # HTML templates (Gin)
├── docker-compose.yml       # App + PostgreSQL + Redis
├── dockerfile               # Multi-stage build
├── go.mod
└── go.sum
```

---

## Getting Started

### Prerequisites

- [Go 1.25+](https://go.dev/dl/)
- [Docker & Docker Compose](https://docs.docker.com/get-docker/) — for containerized setup
- PostgreSQL 16 + Redis — if running locally without Docker

### Run with Docker (recommended)

```bash
# Clone the repo
git clone https://github.com/your-username/Go_shortenURL.git
cd Go_shortenURL

# Start all services
docker compose up --build
```

App is available at [http://localhost:8080](http://localhost:8080).

### Run locally

```bash
# Copy and configure environment variables
cp configs/config.example.env configs/config.env
# Edit configs/config.env to match your local Postgres/Redis

go run ./cmd/api
```

> Requires a running PostgreSQL and Redis instance matching the values in `config.env`.

---

## Environment Variables

The app reads config from environment variables (Docker Compose) or `configs/config.env` (local).

| Variable | Description | Example |
|---|---|---|
| `APP_PORT` | HTTP server port | `8080` |
| `DB_HOST` | PostgreSQL host | `localhost` |
| `DB_PORT` | PostgreSQL port | `5432` |
| `DB_USER` | PostgreSQL user | `postgres` |
| `DB_PASSWORD` | PostgreSQL password | `secret` |
| `DB_NAME` | PostgreSQL database | `shortener` |
| `REDIS_ADDR` | Redis address | `localhost:6379` |
| `REDIS_PASSWORD` | Redis password | `123` |
| `REDIS_DB` | Redis DB index | `0` |

> Never commit `configs/config.env` — it is git-ignored. Use `config.example.env` as the template.

---

## API Reference

| Method | Path | Description |
|---|---|---|
| `GET` | `/` | Serve the shorten form (HTML) |
| `POST` | `/shorten` | Submit form field `url` → returns short URL |
| `GET` | `/:shortCode` | Redirect to original URL |
| `GET` | `/health` | Health check → `{"status":"ok"}` |

### Example

```bash
# Shorten a URL
curl -X POST http://localhost:8080/shorten \
  -d "url=https://example.com/very/long/path"

# Redirect (follow the redirect)
curl -L http://localhost:8080/aB3xYz
```

---

## How Short Codes Are Generated

Short codes are 6-character random strings sampled from a base-62 alphabet (`0-9a-zA-Z`), giving 62⁶ ≈ **56 billion** possible codes.

```
ALPHABET = "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
```

Collision on insert is handled by `ON CONFLICT (short_code) DO UPDATE` in PostgreSQL.

---

## Docker Services

| Service | Image | Port | Notes |
|---|---|---|---|
| `app` | local build | `8080` | Depends on Redis health check |
| `postgres` | `postgres:16-alpine` | `5432` | Volume: `postgres_data` |
| `redis` | `redis:latest` | `6379` | AOF persistence, password protected |

---

## License

MIT
