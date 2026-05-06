# TipDrop

TipDrop is a Flutter-first direct PromptPay tipping app with a Go Fiber API, PostgreSQL, Redis, RabbitMQ, S3-compatible storage, a Go background worker, and a small Next.js support web app.

## Core Rules

- Workers receive 100% of tips.
- The platform never holds money or implements a wallet balance.
- The backend is the source of truth for tip state transitions.
- Slip images must use the signed S3 upload flow.

## Repository Layout

- `apps/mobile` - Flutter primary app
- `apps/web` - Next.js support web app
- `services/api` - Go Fiber API
- `services/worker` - Go background worker
- `infra/docker` - local Docker Compose infrastructure
- `infra/migrations` - SQL database migrations
- `docs` - developer documentation
- `issues` - planning and task documents

## Local Development

1. Copy `.env.example` to `.env` and fill local values.
2. Start local dependencies:

```sh
make docker-up
```

Equivalent command:

```sh
docker compose -f infra/docker/docker-compose.yml up -d
```

3. Run the API once Go is installed:

```sh
cd services/api
go run ./cmd/api
```

4. Run the worker once Go is installed:

```sh
cd services/worker
go run ./cmd/worker
```

5. Install and run the web app:

```sh
npm install
npm --workspace apps/web run dev
```

6. Generate Flutter platform folders and run the mobile app once Flutter is installed:

```sh
cd apps/mobile
flutter create .
flutter pub get
flutter run
```

## Checks

Use the root `Makefile` for repeatable checks:

```sh
make test-go
make test-web
make build-web
```

The root itself is not a Go module, so prefer `make test-go` instead of `go test ./...` from the repository root.

## Active Task Plan

Use `issues/TASKS_MASTER_V2.md` as the active source of truth. The V2 plan makes Flutter the primary product surface, keeps Next.js as support pages, uses Google/Facebook OAuth for login, and reserves OTP for worker payment-profile verification.

## Documentation

- `docs/ARCHITECTURE.md` - system architecture and responsibility boundaries
- `docs/AI_TASK_GUIDE.md` - required task format and AI guardrails
- `docs/ENVIRONMENT.md` - environment variable contract
- `docs/EVENTS.md` - RabbitMQ event envelope and event names
- `docs/DEPLOYMENT.md` - staging and production checklist placeholder
