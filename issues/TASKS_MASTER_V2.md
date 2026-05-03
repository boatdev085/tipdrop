# TipDrop Task Master V2

## Direction Changes

This version replaces the previous web-heavy and OTP-login plan.

1. Flutter is the main app for core product features.
2. Next.js is a supporting web app for leaderboard, discover, public profile preview, SEO, and legal pages.
3. Core tip flow, slip upload, worker confirmation, and dashboard should be implemented in Flutter first.
4. Login uses Google and Facebook OAuth.
5. OTP is not used for login.
6. OTP is used to verify the worker payment profile before the worker can receive tips.
7. OAuth identity does not mean the worker can receive tips.
8. Worker must have a verified payment profile before tip initiation returns payment information.

---

## Architecture

```txt
Flutter App
  -> Go Fiber API
  -> PostgreSQL
  -> Redis
  -> RabbitMQ
  -> S3
  -> Go Worker Service
  -> FCM
```

```txt
Next.js Web
  -> Leaderboard
  -> Discover
  -> Public worker profile preview
  -> Legal pages
```

---

# Phase 0 — Foundation

## T00.01 — Setup Monorepo

**Stack:** infra/docs

**Goal:** Create repo structure.

**Steps:**
1. Create `apps/mobile`.
2. Create `apps/web`.
3. Create `services/api`.
4. Create `services/worker`.
5. Create `infra/migrations`.
6. Create `infra/docker`.
7. Create `docs` and `issues`.
8. Update README with Flutter-first direction.

**Output:** Monorepo structure ready.

---

## T00.02 — Setup Docker Infrastructure

**Stack:** infra

**Goal:** Run local dependencies.

**Services:**
- PostgreSQL
- Redis
- RabbitMQ
- MinIO or S3-compatible storage

**Output:** `docker-compose.yml` and local dev services.

---

## T00.03 — Setup Environment Contract

**Stack:** infra/backend/mobile/frontend

**Goal:** Document required env values for API, Flutter, Web, OAuth, queue, storage, and notifications.

**Output:** `.env.example` and `docs/ENVIRONMENT.md`.

**Notes:** Do not commit real secrets.

---

## T00.04 — Setup Go Fiber API

**Stack:** backend

**Goal:** Create API skeleton.

**Steps:**
1. Initialize Go module.
2. Add Fiber.
3. Add config loader.
4. Add logger.
5. Add `/healthz`.
6. Add `/readyz`.
7. Add graceful shutdown.

**Output:** API runs locally.

---

## T00.05 — Setup Backend Connections

**Stack:** backend/infra

**Goal:** Connect API to PostgreSQL, Redis, RabbitMQ, and S3.

**Output:** API readiness checks include all dependencies.

---

## T00.06 — Setup Migration Tool

**Stack:** database/infra

**Goal:** Add repeatable SQL migrations.

**Output:** migration command and first migration placeholder.

---

## T00.07 — Setup Next.js Supporting Web App

**Stack:** frontend

**Goal:** Create Next.js only for support surfaces, not core app flow.

**Routes to Stub:**
- `/`
- `/leaderboard`
- `/discover`
- `/@username`
- `/legal/privacy`
- `/legal/terms`

**Out of Scope for Web MVP:**
- Main worker dashboard
- Main donor tip flow
- Slip upload flow
- Tip confirmation flow

**Output:** `apps/web` builds with support pages.

---

## T00.08 — Setup Flutter App

**Stack:** mobile

**Goal:** Create primary app shell.

**Screens to Stub:**
- Splash
- OAuth login
- Home
- Public profile setup
- Payment profile verification
- Donor tip flow
- Worker dashboard
- Leaderboard
- Discover
- QR scanner

**Output:** Flutter app runs locally.

---

# Phase 1 — OAuth Login

## T01.01 — OAuth User Schema

**Stack:** database/backend

**Goal:** Store users and OAuth accounts.

**Tables:**
- `users`
- `oauth_accounts`
- `refresh_tokens` if token rotation is persisted
- `notification_tokens`
- `audit_logs`

**Important Rule:** No OTP login tables are part of auth login.

**Output:** OAuth-ready identity schema.

---

## T01.02 — Google OAuth Login API

**Stack:** backend/mobile

**Goal:** Flutter logs in with Google, backend verifies provider identity, then backend issues app session.

**Endpoint:** `POST /api/auth/oauth/google`

**Steps:**
1. Accept provider credential from Flutter.
2. Verify with Google.
3. Extract provider user id, email, display name, and avatar if available.
4. Find or create user.
5. Link OAuth account.
6. Issue app access and refresh token.
7. Write audit log.

**Output:** App session from Google login.

---

## T01.03 — Facebook OAuth Login API

**Stack:** backend/mobile

**Goal:** Flutter logs in with Facebook, backend verifies provider identity, then backend issues app session.

**Endpoint:** `POST /api/auth/oauth/facebook`

**Steps:**
1. Accept provider credential from Flutter.
2. Verify with Facebook.
3. Extract provider user id, email if available, display name, and avatar if available.
4. Find or create user.
5. Link OAuth account.
6. Issue app access and refresh token.
7. Write audit log.

**Output:** App session from Facebook login.

---

## T01.04 — App JWT Session System

**Stack:** backend/security

**Goal:** Create app-owned session after OAuth login.

**Steps:**
1. Create access token.
2. Create refresh token.
3. Add refresh endpoint.
4. Add token validation helpers.
5. Optional: persist refresh token hash.

**Output:** App session system independent from OAuth provider session.

---

## T01.05 — Auth Middleware

**Stack:** backend/security

**Goal:** Protect authenticated endpoints.

**Steps:**
1. Parse bearer token.
2. Validate JWT.
3. Attach user to request context.
4. Add `RequireAuth`.
5. Add `RequireWorker`.
6. Add ownership checks.

**Output:** Protected API routes are possible.

---

## T01.06 — Flutter Google and Facebook Login UI

**Stack:** mobile

**Goal:** Implement login buttons and token exchange.

**Steps:**
1. Add Google login button.
2. Add Facebook login button.
3. Integrate provider SDKs/packages.
4. Send provider credential to backend.
5. Store app tokens securely.
6. Route user based on onboarding state.

**Output:** User can login with Google or Facebook.

---

# Phase 2 — Worker Profile and Payment Profile Verification

## T02.01 — Profile and Payment Profile Schema

**Stack:** database/backend

**Goal:** Separate public profile from verified payment profile.

**Tables:**
- `worker_profiles`
- `payment_profiles`
- `payment_profile_otp_codes`

**Important Rule:** Public profile completion does not mean payment profile is verified.

**Output:** Schema can represent payment eligibility.

---

## T02.02 — Public Worker Profile API

**Stack:** backend/mobile/frontend

**Goal:** Create and update public worker profile.

**Endpoint:** `PATCH /api/profile`

**Fields:**
- username
- display_name
- avatar_url
- bio
- job_title

**Output:** Worker has public profile.

---

## T02.03 — Payment Profile API

**Stack:** backend/mobile

**Goal:** Worker enters PromptPay/payment information for receiving tips.

**Endpoint:** `PATCH /api/payment-profile`

**Fields:**
- promptpay_id
- promptpay_type
- phone if phone-based verification is used

**Behavior:**
- Changing payment information resets verification status.
- Worker cannot receive tips until verified.

**Output:** Payment profile exists but may be unverified.

---

## T02.04 — Send Payment Profile OTP API

**Stack:** backend/security

**Goal:** Send OTP for verifying payment profile ownership.

**Endpoint:** `POST /api/payment-profile/otp/send`

**Steps:**
1. Require auth.
2. Load payment profile.
3. Generate OTP.
4. Store hash with short expiry.
5. Rate limit by user and phone.
6. In dev mode, log OTP. In production, plug SMS provider.

**Output:** Payment profile OTP is ready for verification.

---

## T02.05 — Verify Payment Profile OTP API

**Stack:** backend/security

**Goal:** Verify OTP and mark payment profile as verified.

**Endpoint:** `POST /api/payment-profile/otp/verify`

**Steps:**
1. Require auth.
2. Find latest unconsumed OTP.
3. Reject expired OTP.
4. Compare hash.
5. Mark OTP consumed.
6. Set payment profile status to verified.
7. Write audit log.

**Output:** Worker becomes eligible to receive tips.

---

## T02.06 — Current User and Onboarding State API

**Stack:** backend/mobile

**Goal:** Tell Flutter what the user needs to do next.

**Endpoint:** `GET /api/me`

**Response Includes:**
- user identity
- public profile status
- payment profile status
- can_receive_tips

**Output:** Flutter can route onboarding correctly.

---

## T02.07 — Flutter Public Profile Setup

**Stack:** mobile

**Goal:** Worker creates public profile in Flutter.

**Output:** Public profile completed.

---

## T02.08 — Flutter Payment Profile Verification

**Stack:** mobile

**Goal:** Worker enters payment info and verifies OTP in Flutter.

**Steps:**
1. Enter payment info.
2. Save payment profile.
3. Request OTP.
4. Enter OTP.
5. Verify OTP.
6. Show verified state.

**Output:** Worker can receive tips.

---

# Phase 3 — Flutter Core Tip Flow

## T03.01 — Tip Request Schema

**Stack:** database/backend

**Goal:** Store tip lifecycle.

**Table:** `tip_requests`

**Statuses:**
- pending
- slip_uploaded
- confirmed
- disputed
- expired

**Output:** Tip persistence ready.

---

## T03.02 — Initiate Tip API

**Stack:** backend/mobile

**Goal:** Create pending tip request.

**Endpoint:** `POST /api/tips/initiate`

**Rules:**
1. Worker must exist.
2. Worker payment profile must be verified.
3. Amount must be valid.
4. Ref code must be unique.
5. Payment info is returned only after eligibility checks.

**Output:** Flutter can render PromptPay QR.

---

## T03.03 — PromptPay QR Helper

**Stack:** backend/mobile

**Goal:** Provide data needed for Flutter to render PromptPay QR.

**Open Decision:** Generate QR payload backend-side or Flutter-side.

**Output:** QR screen can display amount, payment destination, and ref code instruction.

---

## T03.04 — S3 Presigned Upload API

**Stack:** backend/mobile/storage

**Goal:** Allow Flutter to upload slip image directly to S3-compatible storage.

**Endpoint:** `POST /api/uploads/presign`

**Rules:**
- Image only.
- Max size policy.
- Random object key.
- No arbitrary external URL accepted.

**Output:** Flutter can upload slip.

---

## T03.05 — Attach Slip to Tip API

**Stack:** backend/mobile

**Goal:** Move tip from pending to slip_uploaded.

**Endpoint:** `POST /api/tips/:id/slip`

**Rules:**
- Tip must exist.
- Current status must be pending.
- File key must be valid.
- Publish `tip.slip_uploaded`.

**Output:** Worker can see pending slip.

---

## T03.06 — Worker Pending Tips API

**Stack:** backend/mobile

**Goal:** Return uploaded slips for Flutter worker dashboard.

**Endpoint:** `GET /api/tips?status=slip_uploaded`

**Rules:**
- Require auth.
- Return only current worker tips.
- Include slip preview URL.

**Output:** Flutter dashboard can list pending tips.

---

## T03.07 — Confirm Tip API

**Stack:** backend/mobile

**Goal:** Worker confirms tip.

**Endpoint:** `PATCH /api/tips/:id/confirm`

**Rules:**
- Require auth.
- Worker must own tip.
- Status must be slip_uploaded.
- Publish `tip.confirmed`.

**Output:** Tip becomes confirmed.

---

## T03.08 — Dispute Tip API

**Stack:** backend/mobile

**Goal:** Worker disputes invalid slip.

**Endpoint:** `PATCH /api/tips/:id/dispute`

**Output:** Tip becomes disputed.

---

## T03.09 — Tip Status API

**Stack:** backend/mobile

**Goal:** Flutter donor flow can poll status.

**Endpoint:** `GET /api/tips/:id/status`

**Output:** Donor sees pending, uploaded, confirmed, disputed, or expired.

---

## T03.10 — Flutter Donor Tip Flow

**Stack:** mobile

**Goal:** Implement core donor tip flow in Flutter.

**Steps:**
1. Open worker profile from QR/link.
2. Show public profile.
3. Select amount.
4. Initiate tip.
5. Show PromptPay QR.
6. Show ref code instruction.
7. Upload slip.
8. Attach slip to tip.
9. Show status page.
10. Poll until terminal status.

**Output:** Donor can tip from Flutter.

---

## T03.11 — Flutter Worker Dashboard MVP

**Stack:** mobile

**Goal:** Worker processes tips in Flutter.

**Steps:**
1. Fetch pending tips.
2. Show amount, donor, ref code, time.
3. Show slip preview.
4. Confirm.
5. Dispute.
6. Refresh state.

**Output:** Worker can confirm/dispute tips in Flutter.

---

# Phase 4 — Queue and Worker Service

## T04.01 — Worker Service Bootstrap

**Stack:** worker

**Goal:** Create background worker service.

**Responsibilities:**
- Consume RabbitMQ events.
- Send notifications.
- Expire old tips.
- Refresh leaderboard.

**Output:** Worker consumes test events.

---

## T04.02 — Event Contract

**Stack:** backend/worker/docs

**Events:**
- tip.created
- tip.slip_uploaded
- tip.confirmed
- tip.disputed
- tip.expired
- notification.send
- leaderboard.refresh
- payment_profile.verified

**Output:** `docs/EVENTS.md`.

---

## T04.03 — Expire Tips Cron

**Stack:** worker/backend

**Rules:**
- pending older than 24h becomes expired.
- slip_uploaded older than 72h becomes expired.

**Output:** Old tips expire automatically.

---

# Phase 5 — Notifications

## T05.01 — Register FCM Token API

**Stack:** backend/mobile

**Endpoint:** `POST /api/notification-tokens`

**Goal:** Flutter registers device token.

**Output:** Worker can receive push.

---

## T05.02 — Flutter Push Handling

**Stack:** mobile

**Goal:** Receive push and route to pending tips.

**Output:** Tap notification opens relevant screen.

---

## T05.03 — FCM Notification Worker

**Stack:** worker/notification

**Goal:** Send push when slip is uploaded.

**Output:** Worker receives notification.

---

## T05.04 — LINE Messaging Optional

**Stack:** worker/notification

**Goal:** Prepare optional LINE provider after MVP.

**Decision:** Not required for MVP unless moved forward.

---

# Phase 6 — Leaderboard and Web Support

## T06.01 — Worker Stats Schema

**Stack:** database/backend

**Goal:** Store aggregate stats.

**Output:** Stats can support dashboard and leaderboard.

---

## T06.02 — Update Stats on Confirmed Tip

**Stack:** backend/worker

**Goal:** Update stats when tip is confirmed.

**Output:** Leaderboard data stays current.

---

## T06.03 — Leaderboard API

**Stack:** backend/mobile/frontend

**Endpoint:** `GET /api/leaderboard`

**Goal:** Expose ranked workers.

**Output:** Flutter and Next.js can render leaderboard.

---

## T06.04 — Flutter Leaderboard Screen

**Stack:** mobile

**Goal:** Show leaderboard in Flutter.

**Output:** Native leaderboard exists.

---

## T06.05 — Next.js Leaderboard Page

**Stack:** frontend

**Goal:** Show public leaderboard on web.

**Route:** `/leaderboard`

**Output:** Supporting web leaderboard exists.

---

# Phase 7 — Discover and QR

## T07.01 — Flutter QR Scanner

**Stack:** mobile

**Goal:** Scan worker QR and start Flutter donor flow.

**Output:** Donor can scan and tip.

---

## T07.02 — Flutter Discover Screen

**Stack:** mobile

**Goal:** Browse workers inside Flutter.

**Output:** App discover screen exists.

---

## T07.03 — Next.js Discover Page Optional

**Stack:** frontend

**Goal:** Public SEO/shareable discover page.

**Decision:** Optional after Flutter discover.

---

# Phase 8 — Workplace Optional

## T08.01 — Workplace Schema

**Stack:** database/backend

**Goal:** Support workplace later.

**Decision:** Optional after MVP.

---

## T08.02 — Workplace APIs

**Stack:** backend/mobile

**APIs:**
- create workplace
- join workplace
- activate workplace

**Rule:** One worker can have only one active workplace.

**Decision:** Optional after MVP.

---

# Phase 9 — Testing

## T09.01 — Backend Unit Tests

**Areas:**
- OAuth linking
- JWT validation
- Payment profile OTP
- Tip status transitions
- Ownership checks
- Expiry logic

---

## T09.02 — Backend Integration Tests

**Scenario:**
1. OAuth login.
2. Create public profile.
3. Create payment profile.
4. Verify payment profile OTP.
5. Initiate tip.
6. Upload slip.
7. Confirm tip.
8. Leaderboard updates.

---

## T09.03 — Flutter Smoke Tests

**Areas:**
- OAuth UI
- Profile setup
- Payment profile OTP
- Donor tip flow
- Worker dashboard

---

## T09.04 — Web Smoke Tests

**Areas:**
- Leaderboard
- Discover
- Public worker profile preview

---

# Phase 10 — Deployment

## T10.01 — Docker Build API and Worker

**Output:** Production-ready API and worker images.

---

## T10.02 — CI Pipeline

**Checks:**
- Backend tests
- Flutter analyze/test
- Web build
- Docker build

---

## T10.03 — Staging Plan

**Output:** `docs/DEPLOYMENT.md`.

---

## T10.04 — Production Checklist

**Checklist:**
- DB backups
- S3 lifecycle
- Error tracking
- Logs and metrics
- Secret management
- Rate limits
- Manual dispute process

---

# Recommended MVP V2

Keep:
1. Flutter app foundation.
2. Google + Facebook OAuth.
3. Public worker profile.
4. Payment profile + OTP verification.
5. Flutter donor tip flow.
6. Flutter worker dashboard.
7. S3 slip upload.
8. FCM notification.
9. Basic leaderboard API.
10. Next.js leaderboard and public profile preview.

Postpone:
1. LINE Messaging API.
2. Workplace system.
3. Next.js authenticated dashboard.
4. Next.js core tip flow.
5. Advanced analytics.
6. Zone-based discover.

---

# Open Decisions

1. PromptPay QR payload: backend-side or Flutter-side?
2. Payment profile MVP: phone PromptPay only or more types?
3. Donor status page: public by ID or signed status token?
4. Web public worker page: allow tipping or redirect to app?
5. Facebook login: launch day requirement or feature flag?
