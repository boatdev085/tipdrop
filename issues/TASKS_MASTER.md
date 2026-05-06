# TipDrop Detailed Task Master List

> [!IMPORTANT]
> This task list is superseded by `issues/TASKS_MASTER_V2.md`. Use V2 as the active source of truth for new implementation work.

This document is the source of truth for reviewing scope before creating individual GitHub Issues. Each task is written so a developer or AI agent can understand the expected implementation, dependency, output, and completion criteria.

## Global Rules

- Worker must receive 100% of the tip.
- The platform must not hold money or implement a wallet in MVP.
- Backend is the source of truth for all state transitions.
- Slip images must be uploaded through signed S3 flow.
- Every implementation must follow `AGENTS.md`, `docs/ARCHITECTURE.md`, and `docs/AI_TASK_GUIDE.md`.

## Label Strategy

Use labels in this pattern:

- `phase:0`, `phase:1`, `phase:2`, etc.
- `stack:backend`, `stack:frontend`, `stack:mobile`, `stack:worker`, `stack:infra`, `stack:docs`
- Optional: `type:feature`, `type:bug`, `type:task`, `priority:high`, `priority:medium`, `priority:low`

---

# Phase 0 — Project Foundation

## T00.01 — Setup Monorepo Structure

**Stack Context**
- Stack: infra/docs
- Tools: Git, monorepo layout

**Goal**
Create the repository structure used by all future work.

**Dependencies**
- Repository exists.

**Implementation Steps**
1. Create `apps/mobile` for Flutter.
2. Create `apps/web` for Next.js.
3. Create `services/api` for Go Fiber API.
4. Create `services/worker` for Go background worker.
5. Create `infra/migrations` for SQL migrations.
6. Create `infra/docker` for local service setup.
7. Create `docs`, `docs/flows`, `docs/mockups`, and `issues`.
8. Update README with local dev overview.

**Files/Output**
- `apps/mobile/.gitkeep`
- `apps/web/.gitkeep`
- `services/api/.gitkeep`
- `services/worker/.gitkeep`
- `infra/migrations/.gitkeep`
- `infra/docker/.gitkeep`
- `README.md`

**Output for Next Task**
All later setup tasks can place code in the correct directories.

**Definition of Done**
- Directory structure exists.
- README explains project layout.

---

## T00.02 — Setup Local Docker Infrastructure

**Stack Context**
- Stack: infra
- Tools: Docker Compose, PostgreSQL, Redis, RabbitMQ, MinIO/S3

**Goal**
Provide local development services for API, worker, and web testing.

**Dependencies**
- T00.01

**Implementation Steps**
1. Create `docker-compose.yml`.
2. Add PostgreSQL service with persistent volume.
3. Add Redis service.
4. Add RabbitMQ service with management UI.
5. Add MinIO service as S3-compatible storage.
6. Add default bucket initialization script.
7. Add `.env.example` values for local services.
8. Add README section for `docker compose up -d`.

**Files/Output**
- `docker-compose.yml`
- `.env.example`
- `infra/docker/minio/create-bucket.sh`

**Acceptance Criteria**
- PostgreSQL reachable on local port.
- Redis responds to ping.
- RabbitMQ management UI is accessible.
- MinIO bucket exists.

**Output for Next Task**
API and worker can connect to local dependencies.

---

## T00.03 — Setup Environment Config Contract

**Stack Context**
- Stack: infra/backend/frontend/mobile
- Tools: dotenv, typed config

**Goal**
Define environment variables used across all services.

**Dependencies**
- T00.02

**Implementation Steps**
1. Create `.env.example` with API, DB, Redis, RabbitMQ, S3, JWT, FCM, LINE keys.
2. Document which services consume each variable.
3. Add validation expectation for backend startup.
4. Add placeholders for dev/staging/prod.

**Required Variables**
```txt
APP_ENV=
API_PORT=
DATABASE_URL=
REDIS_URL=
RABBITMQ_URL=
S3_ENDPOINT=
S3_BUCKET=
S3_ACCESS_KEY=
S3_SECRET_KEY=
JWT_ACCESS_SECRET=
JWT_REFRESH_SECRET=
FCM_SERVER_KEY=
LINE_CHANNEL_ACCESS_TOKEN=
WEB_PUBLIC_API_BASE_URL=
FLUTTER_API_BASE_URL=
```

**Output**
- `.env.example`
- `docs/ENVIRONMENT.md`

**Definition of Done**
- Every required env var is documented.
- Missing required backend env should fail fast.

---

## T00.04 — Setup Go Fiber API Bootstrap

**Stack Context**
- Stack: backend
- Language/Framework: Go, Fiber

**Goal**
Create the initial API service skeleton.

**Dependencies**
- T00.01, T00.03

**Implementation Steps**
1. Initialize Go module under `services/api`.
2. Add Fiber server.
3. Add config loader.
4. Add structured logger.
5. Add route group `/api`.
6. Add `GET /healthz` and `GET /readyz`.
7. Add graceful shutdown.

**Files/Output**
- `services/api/go.mod`
- `services/api/cmd/api/main.go`
- `services/api/internal/config`
- `services/api/internal/http`

**Acceptance Criteria**
- API starts locally.
- `/healthz` returns 200.
- Server exits gracefully.

**Output for Next Task**
Backend connection setup can be added.

---

## T00.05 — Setup Backend DB Connection

**Stack Context**
- Stack: backend/database
- Tools: PostgreSQL, pgx or sqlc-compatible access

**Goal**
Connect API to PostgreSQL and expose DB health checks.

**Dependencies**
- T00.02, T00.04

**Implementation Steps**
1. Add PostgreSQL client package.
2. Read `DATABASE_URL` from config.
3. Add connection pool.
4. Ping DB in `/readyz`.
5. Add basic migration command placeholder.

**Output**
- `services/api/internal/db`
- DB readiness included in `/readyz`

**Definition of Done**
- API fails fast if DB config missing.
- `/readyz` fails when DB is unavailable.

---

## T00.06 — Setup Redis Client

**Stack Context**
- Stack: backend/infra
- Tools: Redis

**Goal**
Prepare Redis for rate limiting, cache, and short-lived state.

**Dependencies**
- T00.02, T00.04

**Implementation Steps**
1. Add Redis client package.
2. Read `REDIS_URL` from env.
3. Ping Redis in `/readyz`.
4. Create helper for `SetNX`, counters, TTL.

**Output**
- `services/api/internal/redis`

**Definition of Done**
- API readiness fails if Redis unavailable.
- Helpers are testable.

---

## T00.07 — Setup RabbitMQ Publisher

**Stack Context**
- Stack: backend/queue
- Tools: RabbitMQ, AMQP

**Goal**
Allow API to publish async domain events.

**Dependencies**
- T00.02, T00.04

**Implementation Steps**
1. Add AMQP client package.
2. Create exchange `tipdrop.events`.
3. Define event envelope structure.
4. Add publisher helper.
5. Add connection health to `/readyz`.

**Event Envelope**
```json
{
  "event_id": "uuid",
  "event_type": "tip.slip_uploaded",
  "occurred_at": "iso8601",
  "payload": {}
}
```

**Output**
- `services/api/internal/queue`
- Event contract documented in `docs/EVENTS.md`

---

## T00.08 — Setup Database Migration Tool

**Stack Context**
- Stack: database/infra
- Tools: SQL migrations

**Goal**
Create repeatable database migrations.

**Dependencies**
- T00.02

**Implementation Steps**
1. Pick migration tool (`golang-migrate` recommended).
2. Create migration command script.
3. Add `000001_init.up.sql` and `000001_init.down.sql` placeholder.
4. Document migration commands.

**Output**
- `infra/migrations/000001_init.up.sql`
- `infra/migrations/000001_init.down.sql`
- `Makefile` or scripts

**Definition of Done**
- Migrate up/down works on local DB.

---

## T00.09 — Setup Next.js Web App

**Stack Context**
- Stack: frontend
- Framework: Next.js App Router, Tailwind CSS

**Goal**
Create web app used inside Flutter WebView and normal browser.

**Dependencies**
- T00.01, T00.03

**Implementation Steps**
1. Initialize Next.js app under `apps/web`.
2. Add Tailwind CSS.
3. Add API client wrapper.
4. Add environment config for API base URL.
5. Add base mobile-first layout.
6. Add placeholder routes.

**Routes to Stub**
- `/@username`
- `/tips/[id]/pay`
- `/tips/[id]/upload`
- `/tips/[id]/status`
- `/dashboard`

**Output**
- `apps/web`
- Stub pages compile.

---

## T00.10 — Setup Flutter App

**Stack Context**
- Stack: mobile
- Framework: Flutter

**Goal**
Create mobile shell app for auth, WebView, push, QR scanning.

**Dependencies**
- T00.01, T00.03

**Implementation Steps**
1. Initialize Flutter app under `apps/mobile`.
2. Add environment config support.
3. Add routing structure.
4. Add HTTP client.
5. Add secure storage dependency.
6. Add placeholder screens: Splash, Login, WebViewShell, Dashboard.

**Output**
- `apps/mobile`
- Flutter app runs on simulator/device.

---

# Phase 1 — Auth Vertical Slice

## T01.01 — Create Initial Database Schema

**Stack Context**
- Stack: database/backend
- Tools: PostgreSQL migrations

**Goal**
Create required tables for users, OTP, notification tokens, and audit logs.

**Dependencies**
- T00.08

**Implementation Steps**
1. Create enum `user_role` with `worker`, `donor`, `both`.
2. Create table `users`.
3. Create table `otp_codes`.
4. Create table `notification_tokens`.
5. Create table `audit_logs`.
6. Add indexes and unique constraints.

**Schema Outline**
```sql
users(id uuid primary key, phone text unique not null, role user_role not null, created_at timestamptz not null);
otp_codes(id uuid primary key, phone text not null, code_hash text not null, expires_at timestamptz not null, consumed_at timestamptz, created_at timestamptz not null);
notification_tokens(id uuid primary key, user_id uuid references users(id), provider text not null, token text not null, created_at timestamptz not null);
audit_logs(id uuid primary key, actor_id uuid, event_type text not null, metadata jsonb, created_at timestamptz not null);
```

**Output**
- Migration with auth tables.

**Definition of Done**
- Up/down migrations work.
- Unique phone constraint works.

---

## T01.02 — OTP Send API

**Stack Context**
- Stack: backend
- Framework: Go Fiber + PostgreSQL + Redis

**Goal**
Send/generate OTP for phone login.

**Dependencies**
- T01.01, T00.06

**Implementation Steps**
1. Add `POST /api/auth/otp/send`.
2. Validate Thai phone format.
3. Generate 6-digit OTP.
4. Hash OTP before saving.
5. Store OTP with 5-minute expiry.
6. Rate limit by phone using Redis.
7. For dev, log OTP; for production, integrate SMS provider later.

**API Contract**
```json
POST /api/auth/otp/send
{ "phone": "0812345678" }
```
Response:
```json
{ "success": true }
```

**Edge Cases**
- Invalid phone → 400.
- Too many requests → 429.
- DB failure → 500.

**Output for Next Task**
Valid OTP rows exist for verification.

---

## T01.03 — OTP Verify API

**Stack Context**
- Stack: backend
- Framework: Go Fiber + PostgreSQL

**Goal**
Verify OTP and create/login user.

**Dependencies**
- T01.02

**Implementation Steps**
1. Add `POST /api/auth/otp/verify`.
2. Find latest unconsumed OTP for phone.
3. Reject expired OTP.
4. Compare hash.
5. Mark OTP consumed.
6. Create user if phone not found.
7. Issue access and refresh tokens.
8. Write audit log `user.login`.

**API Contract**
```json
POST /api/auth/otp/verify
{ "phone": "0812345678", "otp": "123456" }
```
Response:
```json
{
  "access_token": "jwt",
  "refresh_token": "jwt",
  "user": { "id": "uuid", "phone": "0812345678", "role": "worker" }
}
```

**Output for Next Task**
Authenticated user session.

**Definition of Done**
- Wrong OTP rejected.
- Expired OTP rejected.
- New user is created once.

---

## T01.04 — JWT Access and Refresh Token System

**Stack Context**
- Stack: backend/security
- Tools: JWT

**Goal**
Implement secure token creation and validation.

**Dependencies**
- T01.03

**Implementation Steps**
1. Create JWT service.
2. Access token TTL: short-lived.
3. Refresh token TTL: longer-lived.
4. Include `sub`, `phone`, `role`, `iat`, `exp` claims.
5. Add token parsing and validation helpers.
6. Add token refresh endpoint.

**API Contract**
```json
POST /api/auth/refresh
{ "refresh_token": "jwt" }
```
Response:
```json
{ "access_token": "jwt" }
```

**Output**
- JWT service usable by middleware and Flutter.

---

## T01.05 — Auth Middleware

**Stack Context**
- Stack: backend/security
- Framework: Go Fiber

**Goal**
Protect worker-only and authenticated endpoints.

**Dependencies**
- T01.04

**Implementation Steps**
1. Create middleware `RequireAuth`.
2. Parse `Authorization: Bearer <token>`.
3. Load current user into request context.
4. Add `RequireWorker` helper.
5. Add ownership helper for worker-owned records.

**Output**
- Middleware available for profile/tip endpoints.

**Definition of Done**
- Missing token → 401.
- Invalid token → 401.
- Valid token attaches user context.

---

## T01.06 — Flutter Auth UI

**Stack Context**
- Stack: mobile
- Framework: Flutter

**Goal**
Allow user to login via phone OTP.

**Dependencies**
- T00.10, T01.02, T01.03

**Implementation Steps**
1. Create phone input screen.
2. Call OTP send API.
3. Create OTP input screen.
4. Call OTP verify API.
5. Store tokens securely.
6. Navigate to profile setup or WebView shell after login.

**Output**
- Flutter user can login.
- Tokens saved in secure storage.

---

## T01.07 — WebView Token Bridge Baseline

**Stack Context**
- Stack: mobile/frontend
- Framework: Flutter WebView + Next.js

**Goal**
Prepare token sharing between Flutter and Next.js WebView screens.

**Dependencies**
- T01.06, T00.09

**Implementation Steps**
1. Create WebView screen in Flutter.
2. Load Next.js URL.
3. Inject access token via JS bridge or header strategy.
4. Add Next.js client helper to read bridge token.
5. Add safe fallback for normal browser.

**Output**
- Next.js can call backend as authenticated user inside Flutter WebView.

**Security Notes**
- Do not persist token in unsafe browser localStorage unless explicitly accepted.
- Prefer short-lived bridge token or in-memory storage.

---

# Phase 2 — Worker Profile

## T02.01 — Worker Profile Database Schema

**Stack Context**
- Stack: database/backend

**Goal**
Add worker profile fields separately from core user identity.

**Dependencies**
- T01.01

**Implementation Steps**
1. Create `worker_profiles` table.
2. Add unique username.
3. Add display fields.
4. Add promptpay ID field.
5. Add avatar URL.
6. Add indexes on username.

**Schema Outline**
```sql
worker_profiles(
  user_id uuid primary key references users(id),
  username text unique not null,
  display_name text not null,
  avatar_url text,
  bio text,
  job_title text,
  promptpay_id text not null,
  created_at timestamptz not null,
  updated_at timestamptz not null
)
```

**Output**
- Profile schema ready.

---

## T02.02 — Create/Update Worker Profile API

**Stack Context**
- Stack: backend
- Framework: Go Fiber + PostgreSQL

**Goal**
Allow authenticated worker to create or update their public profile.

**Dependencies**
- T02.01, T01.05

**Implementation Steps**
1. Add `PATCH /api/profile`.
2. Validate username format.
3. Validate promptpay ID presence.
4. Upsert profile for current user.
5. Ensure username uniqueness.
6. Return normalized profile.

**API Contract**
```json
PATCH /api/profile
{
  "username": "nun",
  "display_name": "นุ่น สุภาพร",
  "bio": "บาริสต้าที่รักกาแฟ",
  "job_title": "บาริสต้า",
  "promptpay_id": "0812345678",
  "avatar_url": "https://..."
}
```

**Output**
- Worker can create/update profile.

---

## T02.03 — Get Current User/Profile API

**Stack Context**
- Stack: backend

**Goal**
Return current authenticated user and profile state.

**Dependencies**
- T02.02

**Implementation Steps**
1. Add `GET /api/me`.
2. Include user fields.
3. Include worker profile if exists.
4. Add `profile_completed` boolean.

**Output**
- Flutter and Next.js can decide next screen.

---

## T02.04 — Public Worker Profile API

**Stack Context**
- Stack: backend

**Goal**
Expose public worker profile for donor flow.

**Dependencies**
- T02.02

**Implementation Steps**
1. Add `GET /api/workers/:username`.
2. Return public fields only.
3. Include rating and stats placeholders.
4. Do not expose private user data.

**API Contract**
```json
GET /api/workers/nun
{
  "worker_id": "uuid",
  "username": "nun",
  "display_name": "นุ่น สุภาพร",
  "job_title": "บาริสต้า",
  "avatar_url": "https://...",
  "bio": "...",
  "rating": 4.9,
  "total_tips": 1240
}
```

**Output for Next Task**
Public web profile can render donor flow.

---

## T02.05 — Flutter Profile Setup Screen

**Stack Context**
- Stack: mobile
- Framework: Flutter

**Goal**
Allow worker to create profile after login.

**Dependencies**
- T02.02, T02.03

**Implementation Steps**
1. Create profile form.
2. Validate required fields locally.
3. Submit to `PATCH /api/profile`.
4. Show username conflict errors.
5. Navigate to WebView/dashboard after success.

**Output**
- Worker can complete onboarding in Flutter.

---

## T02.06 — Next.js Public Worker Profile Page

**Stack Context**
- Stack: frontend
- Framework: Next.js

**Goal**
Render donor-facing worker profile at `/@username`.

**Dependencies**
- T02.04, T00.09

**Implementation Steps**
1. Create route `/@username`.
2. Fetch public profile.
3. Render avatar, display name, job title, workplace placeholder, rating.
4. Render tip amount selector.
5. Add next button to initiate tip later.

**Output**
- Donor can view worker profile.

---

# Phase 3 — Tip Flow Core

## T03.01 — Tip Request Database Schema

**Stack Context**
- Stack: database/backend

**Goal**
Create database model for tip lifecycle.

**Dependencies**
- T02.01

**Implementation Steps**
1. Create enum `tip_status`.
2. Create table `tip_requests`.
3. Add worker foreign key.
4. Add amount, ref code, status, slip key/url.
5. Add donor optional fields.
6. Add confirmed/expired timestamps.
7. Add indexes by worker/status/created_at.

**Schema Outline**
```sql
tip_requests(
  id uuid primary key,
  worker_id uuid references users(id) not null,
  amount integer not null,
  ref_code text unique not null,
  status tip_status not null,
  slip_key text,
  donor_name text,
  donor_note text,
  rating integer,
  confirmed_at timestamptz,
  expires_at timestamptz not null,
  created_at timestamptz not null
)
```

**Output**
- Tip persistence ready.

---

## T03.02 — Generate Ref Code Service

**Stack Context**
- Stack: backend

**Goal**
Generate unique donor-visible reference code.

**Dependencies**
- T03.01

**Implementation Steps**
1. Create ref code generator format `TIP-XXXX` or safer extended format.
2. Check uniqueness against DB.
3. Retry on collision.
4. Unit test collision behavior.

**Output**
- Backend service can generate ref code for tip initiation.

---

## T03.03 — Initiate Tip API

**Stack Context**
- Stack: backend

**Goal**
Create pending tip request for donor.

**Dependencies**
- T03.01, T03.02, T02.04

**Implementation Steps**
1. Add `POST /api/tips/initiate`.
2. Validate worker exists and has profile.
3. Validate amount > 0 and within allowed range.
4. Generate ref code.
5. Create tip with status `pending`.
6. Set `expires_at = now + 24h`.
7. Return promptpay data.

**API Contract**
```json
POST /api/tips/initiate
{
  "worker_id": "uuid",
  "amount": 50,
  "donor_name": "Palm",
  "donor_note": "ขอบคุณครับ"
}
```
Response:
```json
{
  "tip_id": "uuid",
  "ref_code": "TIP-2847",
  "amount": 50,
  "promptpay_id": "0812345678",
  "expires_at": "iso8601"
}
```

**Output for Next Task**
QR page can display amount, ref code, and promptpay ID.

---

## T03.04 — PromptPay Payload/QR Helper

**Stack Context**
- Stack: backend/frontend

**Goal**
Generate or provide PromptPay payment payload metadata for QR rendering.

**Dependencies**
- T03.03

**Implementation Steps**
1. Decide QR generation location: backend payload vs frontend library.
2. Document PromptPay payload format.
3. Ensure amount is encoded.
4. Ensure ref code is shown as transfer note instruction.

**Output**
- QR page has enough data to render PromptPay QR.

---

## T03.05 — S3 Presigned Upload API

**Stack Context**
- Stack: backend/storage
- Tools: S3-compatible object storage

**Goal**
Create signed upload URL for slip image.

**Dependencies**
- T00.02, T00.03

**Implementation Steps**
1. Add S3 client package.
2. Add `POST /api/uploads/presign`.
3. Validate purpose = `tip_slip`.
4. Validate MIME type.
5. Generate random object key.
6. Return signed PUT URL and file key.
7. Store upload session if needed.

**API Contract**
```json
POST /api/uploads/presign
{ "file_type": "image/jpeg", "purpose": "tip_slip" }
```
Response:
```json
{
  "upload_url": "https://...",
  "file_key": "slips/2026/05/uuid.jpg"
}
```

**Definition of Done**
- Only image MIME types allowed.
- Upload works against MinIO locally.

---

## T03.06 — Attach Slip to Tip API

**Stack Context**
- Stack: backend

**Goal**
Mark tip as slip uploaded after donor uploads image.

**Dependencies**
- T03.03, T03.05, T00.07

**Implementation Steps**
1. Add `POST /api/tips/:id/slip`.
2. Validate tip exists.
3. Validate current status is `pending`.
4. Validate file key belongs to allowed slip path/session.
5. Set status to `slip_uploaded`.
6. Set `expires_at = now + 72h`.
7. Publish event `tip.slip_uploaded`.

**API Contract**
```json
POST /api/tips/:id/slip
{ "file_key": "slips/2026/05/uuid.jpg" }
```
Response:
```json
{ "tip_id": "uuid", "status": "slip_uploaded" }
```

**Output for Next Task**
Worker pending list and notification can consume uploaded slips.

---

## T03.07 — Worker Pending Tips API

**Stack Context**
- Stack: backend

**Goal**
Allow worker to view pending uploaded slips.

**Dependencies**
- T03.06, T01.05

**Implementation Steps**
1. Add `GET /api/tips?status=slip_uploaded`.
2. Require worker auth.
3. Return only tips belonging to current worker.
4. Include signed read URL or proxied slip URL.
5. Add pagination.

**Output**
- Worker dashboard can list pending slips.

---

## T03.08 — Confirm Tip API

**Stack Context**
- Stack: backend

**Goal**
Allow worker to confirm valid slip.

**Dependencies**
- T03.07, T00.07

**Implementation Steps**
1. Add `PATCH /api/tips/:id/confirm`.
2. Require auth.
3. Ensure current user owns the tip.
4. Ensure status is `slip_uploaded`.
5. Set status to `confirmed` and `confirmed_at`.
6. Publish `tip.confirmed`.
7. Add audit log.

**Output**
- Tip confirmed and leaderboard can update.

---

## T03.09 — Dispute Tip API

**Stack Context**
- Stack: backend

**Goal**
Allow worker to dispute invalid slip.

**Dependencies**
- T03.07

**Implementation Steps**
1. Add `PATCH /api/tips/:id/dispute`.
2. Require auth.
3. Check worker ownership.
4. Ensure status is `slip_uploaded`.
5. Set status to `disputed`.
6. Publish `tip.disputed`.
7. Add audit log.

**Output**
- Disputed tips are tracked.

---

## T03.10 — Donor Tip Status API

**Stack Context**
- Stack: backend

**Goal**
Allow donor to poll tip status without login.

**Dependencies**
- T03.03

**Implementation Steps**
1. Add `GET /api/tips/:id/status`.
2. Return safe public fields only.
3. Consider adding public status token later.
4. Do not expose worker private data.

**API Response**
```json
{
  "tip_id": "uuid",
  "status": "confirmed",
  "ref_code": "TIP-2847",
  "amount": 50
}
```

**Output**
- Donor status page can poll.

---

## T03.11 — Next.js Tip Amount Selector

**Stack Context**
- Stack: frontend

**Goal**
Enable donor to choose tip amount on worker profile.

**Dependencies**
- T02.06, T03.03

**Implementation Steps**
1. Add preset buttons: 20, 50, 100, 200.
2. Optional custom amount.
3. Call initiate API.
4. Navigate to `/tips/[id]/pay`.
5. Store returned promptpay metadata in route state or refetch by API.

**Output**
- Donor can start tip flow.

---

## T03.12 — Next.js PromptPay QR Page

**Stack Context**
- Stack: frontend

**Goal**
Show PromptPay QR, amount, and ref code instructions.

**Dependencies**
- T03.04, T03.11

**Implementation Steps**
1. Render amount.
2. Render QR.
3. Render ref code prominently.
4. Explain donor must add ref code in bank note.
5. Add CTA to upload slip.

**Output**
- Donor can transfer and proceed to upload.

---

## T03.13 — Next.js Slip Upload Page

**Stack Context**
- Stack: frontend

**Goal**
Upload slip image to S3 and attach it to tip.

**Dependencies**
- T03.05, T03.06, T03.12

**Implementation Steps**
1. Add file input/camera support.
2. Validate file type and size client-side.
3. Request presigned URL.
4. Upload file to S3.
5. Call attach slip API.
6. Navigate to status page.

**Output**
- Slip upload flow complete.

---

## T03.14 — Next.js Donor Status Page

**Stack Context**
- Stack: frontend

**Goal**
Show donor status after slip upload.

**Dependencies**
- T03.10, T03.13

**Implementation Steps**
1. Poll status every 5–10 seconds.
2. Render pending/slip_uploaded/confirmed/disputed/expired states.
3. Stop polling after terminal status.
4. Show thank-you screen on confirmed.

**Output**
- Donor has complete post-upload experience.

---

## T03.15 — Worker Dashboard Web MVP

**Stack Context**
- Stack: frontend

**Goal**
Provide worker pending slip list and confirm/dispute actions in WebView.

**Dependencies**
- T03.07, T03.08, T03.09, T01.07

**Implementation Steps**
1. Create `/dashboard` page.
2. Fetch pending tips.
3. Show donor name, ref code, amount, slip preview.
4. Add confirm button.
5. Add dispute button.
6. Refresh list after action.

**Output**
- Worker can process tips via WebView.

---

# Phase 4 — Queue and Background Worker

## T04.01 — Worker Service Bootstrap

**Stack Context**
- Stack: worker
- Language: Go
- Tools: RabbitMQ, PostgreSQL, Redis

**Goal**
Create background worker service.

**Dependencies**
- T00.02, T00.07

**Implementation Steps**
1. Initialize Go module under `services/worker` or share module.
2. Load env config.
3. Connect RabbitMQ.
4. Connect DB.
5. Connect Redis.
6. Add graceful shutdown.
7. Add test consumer.

**Output**
- Worker can consume messages.

---

## T04.02 — Domain Event Contracts

**Stack Context**
- Stack: backend/worker/docs

**Goal**
Document all events used by API and worker.

**Dependencies**
- T00.07

**Events**
- `tip.created`
- `tip.slip_uploaded`
- `tip.confirmed`
- `tip.disputed`
- `tip.expired`
- `notification.send`
- `leaderboard.refresh`

**Output**
- `docs/EVENTS.md`
- Shared event envelope structs.

---

## T04.03 — Notification Job Consumer

**Stack Context**
- Stack: worker/notification

**Goal**
Consume `tip.slip_uploaded` and trigger notification send.

**Dependencies**
- T04.01, T04.02, T05.01

**Implementation Steps**
1. Consume `tip.slip_uploaded`.
2. Load worker notification tokens.
3. Create notification job.
4. Send via FCM/LINE providers.
5. Ack on success.
6. Retry on transient failure.

**Output**
- Worker receives alert when donor uploads slip.

---

## T04.04 — Retry and Dead Letter Queue

**Stack Context**
- Stack: worker/queue

**Goal**
Handle failed background jobs safely.

**Dependencies**
- T04.01

**Implementation Steps**
1. Configure retry exchange/queue.
2. Add max retry count.
3. Send failed messages to DLQ.
4. Log DLQ metadata.
5. Document how to inspect failed jobs.

**Output**
- Failed jobs are not lost silently.

---

## T04.05 — Expire Tips Cron

**Stack Context**
- Stack: worker/backend

**Goal**
Auto-expire old pending or uploaded tips.

**Dependencies**
- T03.01, T04.01

**Rules**
- `pending` older than 24h → `expired`.
- `slip_uploaded` older than 72h → `expired`.

**Implementation Steps**
1. Add scheduled job every hour.
2. Query expired candidates.
3. Update status atomically.
4. Publish `tip.expired`.
5. Write audit log.

**Output**
- Tip lifecycle is automatically cleaned.

---

# Phase 5 — Notification

## T05.01 — Notification Token Registration API

**Stack Context**
- Stack: backend/mobile

**Goal**
Allow Flutter app to register FCM token.

**Dependencies**
- T01.05, T01.01

**Implementation Steps**
1. Add `POST /api/notification-tokens`.
2. Require auth.
3. Accept provider and token.
4. Upsert token by user/provider/token.
5. Add delete endpoint for logout later.

**API Contract**
```json
POST /api/notification-tokens
{ "provider": "fcm", "token": "..." }
```

**Output**
- Worker can receive push notifications.

---

## T05.02 — FCM Send Provider

**Stack Context**
- Stack: worker/notification

**Goal**
Send push notifications through FCM.

**Dependencies**
- T05.01, T04.03

**Implementation Steps**
1. Add FCM client.
2. Build notification payload.
3. Handle invalid tokens.
4. Log send result.

**Output**
- Push notification reaches Flutter device.

---

## T05.03 — LINE Messaging API Provider

**Stack Context**
- Stack: worker/notification

**Goal**
Send LINE messages for workers who connect LINE.

**Dependencies**
- T04.03

**Implementation Steps**
1. Add LINE Messaging API client.
2. Store user LINE identifier later if needed.
3. Send message template for pending slip.
4. Log result.

**Output**
- LINE notification provider is ready or stubbed behind feature flag.

---

## T05.04 — Flutter Push Notification Handling

**Stack Context**
- Stack: mobile

**Goal**
Receive push notification and route to pending slip screen.

**Dependencies**
- T05.01, T05.02

**Implementation Steps**
1. Setup Firebase Messaging.
2. Request notification permission.
3. Register FCM token with backend.
4. Handle foreground notification.
5. Handle notification tap.
6. Open WebView dashboard or native dashboard.

**Output**
- Worker receives and opens pending tip from push.

---

# Phase 6 — Leaderboard and Stats

## T06.01 — Worker Rating/Stats Schema

**Stack Context**
- Stack: database/backend

**Goal**
Store aggregate stats for worker dashboard and leaderboard.

**Dependencies**
- T03.01

**Implementation Steps**
1. Create `worker_ratings` or `worker_stats` table.
2. Track average rating.
3. Track total confirmed tips.
4. Track weekly amount.
5. Add `week_reset_at`.

**Output**
- Stats persistence ready.

---

## T06.02 — Update Stats on Tip Confirmed

**Stack Context**
- Stack: backend/worker

**Goal**
Update worker stats when tip is confirmed.

**Dependencies**
- T03.08, T06.01

**Implementation Steps**
1. Update stats inside confirm transaction or consume event.
2. Increment total tips.
3. Increment weekly amount.
4. Recalculate rating if rating exists.
5. Publish leaderboard refresh if async.

**Output**
- Confirmed tips affect dashboard and leaderboard.

---

## T06.03 — Redis Leaderboard Cache

**Stack Context**
- Stack: backend/worker/cache

**Goal**
Cache leaderboard results for fast read.

**Dependencies**
- T06.02, T00.06

**Implementation Steps**
1. Define Redis keys.
2. Build weekly leaderboard query.
3. Cache JSON result with TTL.
4. Invalidate/refresh on `tip.confirmed`.

**Keys**
```txt
leaderboard:weekly:global
leaderboard:weekly:zone:{zone}
worker:stats:{worker_id}
```

**Output**
- Leaderboard API can read cache.

---

## T06.04 — Leaderboard API

**Stack Context**
- Stack: backend

**Goal**
Expose leaderboard to web and mobile.

**Dependencies**
- T06.03

**Implementation Steps**
1. Add `GET /api/leaderboard?zone=&week=current`.
2. Read Redis cache first.
3. Fallback to DB query.
4. Return ranked workers.

**Output**
- Web discover/leaderboard can render.

---

## T06.05 — Leaderboard Page

**Stack Context**
- Stack: frontend

**Goal**
Render weekly leaderboard.

**Dependencies**
- T06.04

**Implementation Steps**
1. Create leaderboard section/page.
2. Fetch leaderboard API.
3. Render rank, worker, workplace placeholder, weekly amount.
4. Add loading/error states.

**Output**
- Donor can browse top workers.

---

# Phase 7 — Flutter Integration

## T07.01 — Flutter WebView Shell

**Stack Context**
- Stack: mobile
- Framework: Flutter WebView

**Goal**
Open Next.js screens inside Flutter app.

**Dependencies**
- T01.07, T00.10

**Implementation Steps**
1. Add WebView package.
2. Create `WebViewScreen(url)`.
3. Add loading/error UI.
4. Restrict allowed domains.
5. Add navigation handling.

**Output**
- Flutter can host web donor/worker screens.

---

## T07.02 — Flutter → Web Token Injection

**Stack Context**
- Stack: mobile/frontend/security

**Goal**
Allow authenticated WebView pages to call backend.

**Dependencies**
- T07.01, T01.06

**Implementation Steps**
1. Read access token from secure storage.
2. Inject token into WebView runtime.
3. Next.js stores token in memory.
4. API client uses token in Authorization header.
5. Handle token expiration.

**Output**
- Authenticated WebView dashboard works.

---

## T07.03 — Web → Flutter JS Bridge

**Stack Context**
- Stack: mobile/frontend

**Goal**
Allow web screens to notify Flutter about important events.

**Dependencies**
- T07.01

**Events**
- `AUTH_EXPIRED`
- `TIP_CONFIRMED`
- `OPEN_NATIVE_DASHBOARD`
- `CLOSE_WEBVIEW`

**Implementation Steps**
1. Define JS bridge contract.
2. Add Flutter message handler.
3. Add Next.js helper `postFlutterMessage`.
4. Add tests/manual QA checklist.

**Output**
- Web and Flutter can communicate safely.

---

## T07.04 — Flutter QR Scanner

**Stack Context**
- Stack: mobile

**Goal**
Scan worker QR and open public profile.

**Dependencies**
- T07.01, T02.06

**Implementation Steps**
1. Add QR scanner package.
2. Request camera permission.
3. Parse scanned URL.
4. Validate TipDrop domain/path.
5. Open profile in WebView.

**Output**
- User can scan worker QR to start tip flow.

---

## T07.05 — Native Worker Dashboard MVP Optional

**Stack Context**
- Stack: mobile

**Goal**
Build native dashboard for frequently used worker actions if WebView UX is insufficient.

**Dependencies**
- T03.07, T03.08, T03.09

**Implementation Steps**
1. Create pending tips list screen.
2. Fetch pending tips API.
3. Show slip preview.
4. Add confirm/dispute actions.
5. Add refresh state.

**Output**
- Worker can process tips natively.

**Decision Point**
This can be postponed if WebView dashboard is acceptable for MVP.

---

# Phase 8 — Workplace and Discover

## T08.01 — Workplace Database Schema

**Stack Context**
- Stack: database/backend

**Goal**
Support worker-created workplaces and active workplace rule.

**Dependencies**
- T02.01

**Implementation Steps**
1. Create `workplaces` table.
2. Create `worker_workplaces` junction table.
3. Add `is_active` boolean.
4. Enforce max 1 active workplace per worker.
5. Add invite code.

**Output**
- Workplace persistence ready.

---

## T08.02 — Create Workplace API

**Stack Context**
- Stack: backend

**Goal**
Allow worker to create workplace.

**Dependencies**
- T08.01, T01.05

**Implementation Steps**
1. Add `POST /api/workplaces`.
2. Require worker auth.
3. Validate name/category.
4. Generate invite code.
5. Add creator as member.

**Output**
- Worker can create workplace.

---

## T08.03 — Join Workplace API

**Stack Context**
- Stack: backend

**Goal**
Allow worker to join via invite code.

**Dependencies**
- T08.02

**Implementation Steps**
1. Add `POST /api/workplaces/join`.
2. Validate invite code.
3. Add membership.
4. Prevent duplicate membership.

**Output**
- Workers can join same workplace.

---

## T08.04 — Activate Workplace API

**Stack Context**
- Stack: backend

**Goal**
Allow only one active workplace per worker.

**Dependencies**
- T08.03

**Implementation Steps**
1. Add `PATCH /api/workplaces/:id/activate`.
2. Verify membership.
3. Deactivate previous active workplace in transaction.
4. Activate selected workplace.

**Output**
- Tip requests can capture current active workplace.

---

## T08.05 — Discover Page

**Stack Context**
- Stack: frontend/backend

**Goal**
Allow donors to discover workers by online/nearby/leaderboard sections.

**Dependencies**
- T06.04, T08.04

**Implementation Steps**
1. Add discover API or extend worker list API.
2. Add `/discover` page.
3. Render online workers placeholder.
4. Render leaderboard section.
5. Render nearby placeholder if geolocation not ready.

**Output**
- Donor can browse workers.

---

# Phase 9 — Testing and QA

## T09.01 — Backend Unit Tests

**Stack Context**
- Stack: backend/testing

**Goal**
Cover critical business logic with unit tests.

**Dependencies**
- T03.08, T03.09, T04.05

**Test Areas**
- OTP validation.
- JWT parsing.
- Tip status transitions.
- Worker ownership checks.
- Ref code uniqueness.
- Expiry logic.

**Output**
- `go test ./...` passes.

---

## T09.02 — Backend Integration Tests

**Stack Context**
- Stack: backend/testing

**Goal**
Validate full API flow against test DB.

**Dependencies**
- T03.14

**Scenario**
1. Send OTP.
2. Verify OTP.
3. Create profile.
4. Get public profile.
5. Initiate tip.
6. Upload slip.
7. Worker confirms.
8. Status becomes confirmed.

**Output**
- End-to-end backend flow tested.

---

## T09.03 — Web E2E Tests

**Stack Context**
- Stack: frontend/testing
- Tool: Playwright

**Goal**
Test donor flow in browser.

**Dependencies**
- T03.14

**Test Areas**
- Worker profile loads.
- Amount selector works.
- QR page renders.
- Slip upload flow completes with mock S3.
- Status page polls.

**Output**
- Web donor flow E2E test.

---

## T09.04 — Flutter Smoke Tests

**Stack Context**
- Stack: mobile/testing

**Goal**
Validate mobile app critical screens.

**Dependencies**
- T07.04

**Test Areas**
- App launches.
- Login screen renders.
- Token storage mock works.
- WebView screen opens.
- QR scanner route handles mock value.

**Output**
- Basic mobile test coverage.

---

# Phase 10 — Deployment

## T10.01 — Docker Build for API and Worker

**Stack Context**
- Stack: infra/backend/worker

**Goal**
Containerize backend services.

**Dependencies**
- T00.04, T04.01

**Implementation Steps**
1. Add `services/api/Dockerfile`.
2. Add `services/worker/Dockerfile`.
3. Use multi-stage builds.
4. Run as non-root user.
5. Add healthcheck command.

**Output**
- API and worker images build locally.

---

## T10.02 — CI Pipeline

**Stack Context**
- Stack: infra/ci
- Tool: GitHub Actions

**Goal**
Run checks on every PR.

**Dependencies**
- T09.01, T09.03, T09.04

**Implementation Steps**
1. Add workflow for backend tests.
2. Add workflow for web lint/build.
3. Add workflow for Flutter analyze/test.
4. Add Docker build check.

**Output**
- `.github/workflows/ci.yml`

---

## T10.03 — Staging Environment Plan

**Stack Context**
- Stack: infra/deploy

**Goal**
Document and prepare staging deployment.

**Dependencies**
- T10.01

**Implementation Steps**
1. Pick hosting for API/worker.
2. Pick managed PostgreSQL.
3. Pick Redis/RabbitMQ provider or VPS.
4. Pick S3 provider.
5. Define staging env variables.
6. Document deployment steps.

**Output**
- `docs/DEPLOYMENT.md`

---

## T10.04 — Production Readiness Checklist

**Stack Context**
- Stack: infra/security/ops

**Goal**
Define requirements before production launch.

**Checklist**
- DB backups.
- S3 lifecycle policy.
- Error tracking.
- Structured logs.
- Metrics.
- Rate limiting.
- Secret management.
- SSL/domain.
- Manual slip dispute procedure.

**Output**
- `docs/PRODUCTION_CHECKLIST.md`

---

# Recommended MVP Cut

If scope needs to be reduced, keep only:

1. Phase 0 foundation.
2. Phase 1 auth.
3. Phase 2 profile.
4. Phase 3 complete tip flow.
5. Phase 4 expiry cron.
6. Minimal notification via FCM only.
7. WebView shell.

Postpone:

- LINE Messaging API.
- Native worker dashboard.
- Workplace/discover.
- Advanced leaderboard by zone.
- Production-grade DLQ dashboard.

# Review Notes

Before creating GitHub Issues, review these decision points:

1. Should dashboard be WebView-only for MVP or native Flutter?
2. Should LINE Messaging API be included in MVP or Phase 2?
3. Should workplace be included before public launch?
4. Should PromptPay QR be generated backend-side or frontend-side?
5. Should donor status URL require a signed public token?
