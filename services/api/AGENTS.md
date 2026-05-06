# AGENTS.md - Backend API

## Scope
These instructions apply to the Go Fiber API in `services/api`.

## Goal
Maintain the backend as the source of truth for TipDrop accounts, tip flows, signed slip uploads, and status transitions.

## Stack and Libraries
- Language/runtime: Go 1.22.
- HTTP framework: Go Fiber v2.52.13 (`github.com/gofiber/fiber/v2`).
- Module path: `github.com/boatdev085/tipdrop/services/api`.
- Expected integrations from the platform architecture: PostgreSQL for persistence, Redis for cache/session-like data, RabbitMQ for events, and S3-compatible storage for signed slip uploads. Add concrete client libraries to `go.mod` only when implementation work requires them.

## API Contract
- Document every new or changed HTTP endpoint, request body, response body, status code, and error shape.
- Keep response payloads explicit and stable for Flutter and Next.js clients.
- Validate all client input in the API before calling storage, queue, or database layers.
- Preserve the direct PromptPay model: never add wallet balances, escrow, custody, or platform-held funds.
- Slip image flows must issue or consume signed S3 upload references only; do not accept arbitrary public URLs as proof of upload.

## DB Changes
- Put schema changes in `infra/migrations` and make migrations safe to run once in order.
- Update Go models, queries, and tests in the same task as the migration.
- Do not store sensitive data without encryption or a documented hashing/tokenization strategy.
- Keep tip status updates aligned with the defined status machine and make transitions auditable.

## Files Changed
For each task, summarize changed API files, related migrations, and any client contract changes required downstream.

## Test Coverage
- Run `go test ./...` from `services/api` for API changes.
- Add or update unit tests for config parsing, validation, status transitions, and handler behavior when applicable.
- If a DB migration is added, include migration verification notes or a local command that exercises it.

## Output Used By Next Task
Leave enough detail for the next agent to update Flutter, Next.js, worker consumers, or docs from the API contract without reverse engineering code.
