# Soccer Manager API

RESTful API for a fantasy football manager game.  
Users can sign up, manage one team, edit team/player data, list players on the transfer market, and buy players from other teams.

## Tech Stack

- Go
- Gin (HTTP framework)
- PostgreSQL
- JWT authentication
- Docker Compose (local database)

## Main Features

- User signup and login with JWT access tokens
- Protected endpoints with Bearer authentication
- Automatic team creation on signup:
  - 20 players total (3 GK, 6 DEF, 6 MID, 5 ATT)
  - Team budget starts at 5,000,000
  - Each player starts with market value 1,000,000
- View and update team data (`name`, `country`)
- View and update owned player data (`firstName`, `lastName`, `country`)
- Put owned players on transfer list with asking price
- Browse transfer market
- Buy listed players transactionally:
  - Buyer budget decreases by asking price
  - Seller budget increases by asking price
  - Player moves to buyer team
  - Player market value increases randomly by 10% to 100%

## Architecture Overview

Project follows a clean, feature-based structure under `internal/`:

- `users` - signup/login use cases
- `auth` - JWT token generation
- `teams` - team retrieval/update
- `players` - player generation and updates
- `transfers` - listing, market querying, buying
- `shared` - config, database, middleware, localization

Each feature is split into:

- `application` (business use cases)
- `infrastructure` (database repositories/views)
- `http` (Gin handlers)
- `domain` (domain types where needed)

## Prerequisites

- Go (latest stable)
- Docker + Docker Compose
- PostgreSQL client tools (`psql`) for manual migrations

## Setup

### 1) Clone and install dependencies

```bash
go mod download
```

### 2) Start PostgreSQL with Docker

```bash
docker compose up -d
```

Database defaults from `docker-compose.yml`:

- Host: `localhost`
- Port: `5433`
- User: `postgres`
- Password: `postgres`
- DB name: `soccer_manager`

### 3) Configure environment

Copy `.env.example` to `.env` and adjust values if needed.

Example values:

```env
APP_ENV=development
APP_PORT=8080

DB_HOST=localhost
DB_PORT=5433
DB_USER=postgres
DB_PASSWORD=postgres
DB_NAME=soccer_manager
DB_SSLMODE=disable

JWT_SECRET=super-secret-key
JWT_EXPIRES_IN=24h

DEFAULT_LOCALE=en
```

## Migration Instructions

Run the SQL migration manually (project includes raw SQL migration files):

```bash
psql "postgresql://postgres:postgres@localhost:5433/soccer_manager?sslmode=disable" -f migrations/000001_init_schema.up.sql
```

To drop schema:

```bash
psql "postgresql://postgres:postgres@localhost:5433/soccer_manager?sslmode=disable" -f migrations/000001_init_schema.down.sql
```

## Run the API

```bash
go run ./cmd/api
```

Health check:

```bash
curl http://localhost:8080/health
```

## Postman Collection

Collection is included at:

- `docs/postman/Soccer-Manager-API.postman_collection.json`

It includes all implemented endpoints and collection variables:

- `baseUrl`
- `token`
- `playerId`
- `transferId`

## Authentication Flow

1. `POST /auth/signup` - create user (team + players are created automatically)
2. `POST /auth/login` - get `accessToken`
3. Send `Authorization: Bearer <accessToken>` on protected endpoints

In Postman, login test script stores `accessToken` into `{{token}}`.

## Endpoint List

Public:

- `GET /health`
- `POST /auth/signup`
- `POST /auth/login`

Protected:

- `GET /me`
- `GET /team`
- `PATCH /team`
- `PATCH /players/:id`
- `POST /transfers/list`
- `GET /transfers`
- `POST /transfers/:id/buy`

## Transfer Flow

Typical flow:

1. User A lists owned player via `POST /transfers/list`
2. User B sees listing via `GET /transfers`
3. User B buys listing via `POST /transfers/:id/buy`
4. Transaction updates:
   - listing status (`ACTIVE -> SOLD`)
   - team budgets (buyer - / seller +)
   - player ownership (team change)
   - player market value increase (random 10%-100%)
   - transfer history insert

## Localization (`Accept-Language`)

The API supports two locales:

- `en` (English)
- `ka` (Georgian)

Behavior:

- Reads `Accept-Language` request header
- If header is missing/unsupported, falls back to `DEFAULT_LOCALE`
- If `DEFAULT_LOCALE` is invalid, falls back to English

Localization files:

- `locales/en.json`
- `locales/ka.json`

## Final Status Checklist

- [x] JWT auth and protected routes
- [x] One team per user, created on signup
- [x] Initial 20-player squad with required positions
- [x] Team and player update endpoints
- [x] Transfer listing, market view, and purchase flow
- [x] Transactional budget + ownership updates on transfer
- [x] Postman collection included and cleaned
- [x] English/Georgian localization via `Accept-Language`
- [ ] Unit tests (recommended, not required)