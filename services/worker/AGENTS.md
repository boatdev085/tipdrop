# AGENTS.md - Backend Worker

## Scope
These instructions apply to the Go background worker in `services/worker`.

## Goal
Process asynchronous TipDrop jobs safely without becoming a payment custodian or overriding API-owned business rules.

## Stack and Libraries
- Language/runtime: Go 1.22.
- Module path: `github.com/boatdev085/tipdrop/services/worker`.
- Current worker has no third-party Go runtime dependency declared in `go.mod`; keep it minimal until a queue, database, or observability client is needed.
- Expected integrations from the platform architecture: RabbitMQ for asynchronous jobs, PostgreSQL for coordinated persisted state, Redis where cache/locking is explicitly required, and S3-compatible storage only through backend-approved flows.

## API Contract
- Treat API-published events and database state as the worker contract.
- Document every new or changed RabbitMQ event name, routing key, payload field, retry behavior, and dead-letter behavior.
- Do not let worker code create tip state transitions that bypass backend validation or the status machine.

## DB Changes
- Put schema changes in `infra/migrations` and coordinate them with API code if the worker reads or writes the same tables.
- Worker writes must be idempotent where jobs can be retried.
- Do not store sensitive data without encryption or documented hashing/tokenization.

## Files Changed
For each task, summarize changed worker files, event contracts, queue configuration, and any API or migration dependencies.

## Test Coverage
- Run `go test ./...` from `services/worker` for worker changes.
- Add tests for config parsing, event payload handling, retry/idempotency logic, and failure paths when applicable.

## Output Used By Next Task
Leave clear event and operational notes so API, infra, and observability work can continue without guessing worker behavior.
