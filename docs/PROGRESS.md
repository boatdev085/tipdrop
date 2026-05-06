# Progress Log

This file tracks development progress for TipDrop.

---

## 2026-05-06

### Completed
- Created Flutter-first monorepo structure for mobile, web, API, worker, infra, docs, and issues.
- Added Go workspace for `services/api` and `services/worker`.
- Added Go Fiber API skeleton with `/healthz`, `/readyz`, `/api` group, and graceful shutdown.
- Added Go worker skeleton with signal-aware shutdown.
- Added Next.js support web skeleton for leaderboard, discover, public profile previews, privacy, and terms.
- Added Flutter primary app shell with stub screens for splash, OAuth login, home, profile setup, payment verification, donor tip flow, worker dashboard, leaderboard, discover, and QR scanner.
- Added local infrastructure for PostgreSQL, Redis, RabbitMQ, and MinIO.
- Added environment and event contract documentation.

### In Progress
- Foundation hardening: dependency lockfiles, repeatable checks, backend dependency connections, and config validation.

### Next
- Install and lock web dependencies.
- Implement real API readiness checks for PostgreSQL, Redis, RabbitMQ, and S3.
- Add database schema migrations for OAuth users and payment-profile verification.
- Implement Google OAuth login API before Facebook OAuth.
- Implement OTP only for worker payment-profile verification, not login.

---

## Superseded 2026-05-03 Notes

The original plan included OTP login tasks. That direction is superseded by `issues/TASKS_MASTER_V2.md`:
- Login uses Google and Facebook OAuth.
- OTP is used only to verify a worker payment profile before the worker can receive tips.

---

## Status Legend

- PLANNED: Task defined but not started
- IN PROGRESS: Work ongoing
- BLOCKED: Waiting for dependency or decision
- DONE: Completed and merged

---

## Notes

- Use `issues/TASKS_MASTER_V2.md` as the active task source of truth.
- Follow `AGENTS.md` and `docs/AI_TASK_GUIDE.md` for all implementations.
- Workers receive 100% of tips; the platform must never hold funds.
