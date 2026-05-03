# TipDrop Task Master V2

## Purpose

This document is the implementation map for TipDrop V2. It replaces the earlier web-heavy and OTP-login plan with a Flutter-first product, a Go Fiber backend, and a small Next.js support site.

The goal of this update is to make each task clear enough for the next developer to pick it up without needing hidden product context. Each task now includes feasibility, dependencies, implementation detail, and acceptance criteria.

## Overall Feasibility Assessment

**Feasibility:** Medium-high for an MVP if scope stays disciplined. The core product is achievable with a small team because the payment flow is PromptPay/slip-confirmation based rather than direct card acquiring or wallet custody. The riskiest areas are OAuth app setup, OTP/SMS provider readiness, slip fraud handling, push notification configuration, and production operations.

**Recommended MVP cut:** Build Flutter login, worker profile, verified PromptPay payment profile, donor tip initiation, slip upload, worker confirmation, FCM notification, and basic leaderboard first. Keep workplace, LINE, advanced discover, and authenticated web dashboard out of MVP.

**Main technical risks:** OAuth provider review/configuration, mobile build setup, secure token storage, OTP rate limiting, S3 upload policy correctness, idempotent tip status transitions, and observability for money-adjacent flows.

**Important product rule:** OAuth verifies identity for app login only. A worker can receive tips only after payment profile verification succeeds.

---

## Direction Changes

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
  -> S3-compatible storage
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

## Task Detail Standard

Each implementation ticket should satisfy these fields before being considered ready:

- **Feasibility:** Easy, Medium, Hard, or Optional.
- **Depends on:** Previous tasks or external setup required.
- **Developer Detail:** What to build and what decisions to preserve.
- **Acceptance Criteria:** Observable checks that prove the task is done.
- **Out of Scope:** Explicit exclusions to prevent scope creep.

---

# Phase 0 — Foundation

## T00.01 — Setup Monorepo

**Stack:** infra/docs

**Feasibility:** Easy.

**Depends on:** Empty or near-empty repository.

**Developer Detail:** Create a monorepo structure with `apps/mobile`, `apps/web`, `services/api`, `services/worker`, `infra/migrations`, `infra/docker`, `docs`, and `issues`. Add placeholder README files where useful so empty directories are tracked. Update the root README to state that Flutter is the primary app and Next.js is supporting web only.

**Acceptance Criteria:** Repository has the agreed folder layout, README explains the V2 direction, and no app-specific dependency setup is mixed into this task.

**Out of Scope:** Implementing Flutter, Next.js, Go, Docker, or CI.

---

## T00.02 — Setup Docker Infrastructure

**Stack:** infra

**Feasibility:** Easy-medium.

**Depends on:** T00.01.

**Developer Detail:** Add local Docker Compose under `infra/docker` or root, with PostgreSQL, Redis, RabbitMQ management UI, and MinIO. Use stable service names that backend env values can reference. Add volumes for local persistence and health checks where supported.

**Acceptance Criteria:** `docker compose up` starts all local dependencies, PostgreSQL accepts connections, Redis responds to ping, RabbitMQ UI is reachable, and MinIO exposes S3-compatible credentials for local dev.

**Out of Scope:** Production infrastructure, cloud provisioning, and secret storage.

---

## T00.03 — Setup Environment Contract

**Stack:** infra/backend/mobile/frontend

**Feasibility:** Easy.

**Depends on:** T00.01, T00.02.

**Developer Detail:** Create `.env.example` and `docs/ENVIRONMENT.md`. Document API port, database URL, Redis URL, RabbitMQ URL, S3 endpoint/bucket/keys, JWT secrets, OAuth client IDs, OAuth redirect settings, FCM settings, public web URL, and mobile deep link scheme. Mark which values are required for local dev versus production.

**Acceptance Criteria:** A new developer can copy `.env.example`, fill local values, and understand which variables belong to API, worker, Flutter, and Next.js.

**Out of Scope:** Real secrets and provider account creation.

---

## T00.04 — Setup Go Fiber API

**Stack:** backend

**Feasibility:** Easy-medium.

**Depends on:** T00.01, T00.03.

**Developer Detail:** Initialize the Go module in `services/api`. Add Fiber, structured logging, config loading, request ID middleware, `/healthz`, `/readyz`, and graceful shutdown. Keep handlers thin and create clear package boundaries for config, HTTP routes, middleware, and dependencies.

**Acceptance Criteria:** API runs locally, `/healthz` returns success without dependencies, `/readyz` is wired for future dependency checks, and shutdown handles interrupt signals cleanly.

**Out of Scope:** Database schema, auth, and business routes.

---

## T00.05 — Setup Backend Connections

**Stack:** backend/infra

**Feasibility:** Medium.

**Depends on:** T00.02, T00.03, T00.04.

**Developer Detail:** Add connection clients for PostgreSQL, Redis, RabbitMQ, and S3-compatible storage. Centralize initialization and cleanup. Make `/readyz` fail if a required dependency is unavailable. Use short timeouts so readiness checks do not hang.

**Acceptance Criteria:** API can connect to all local Docker services, readiness reports dependency failures clearly, and connection cleanup runs during shutdown.

**Out of Scope:** Repository methods, queues consumers, and object upload routes.

---

## T00.06 — Setup Migration Tool

**Stack:** database/infra

**Feasibility:** Easy-medium.

**Depends on:** T00.02, T00.05.

**Developer Detail:** Choose a Go-friendly migration tool such as golang-migrate or goose. Add a migration command documented in README or `docs/ENVIRONMENT.md`. Add an initial placeholder or extension migration and define naming rules for future migrations.

**Acceptance Criteria:** Developer can run migrations up/down against local PostgreSQL, migration files live under `infra/migrations`, and failed migrations are visible in command output.

**Out of Scope:** Full product schema.

---

## T00.07 — Setup Next.js Supporting Web App

**Stack:** frontend

**Feasibility:** Easy.

**Depends on:** T00.01, T00.03.

**Developer Detail:** Create the Next.js app in `apps/web`. Stub `/`, `/leaderboard`, `/discover`, `/@username`, `/legal/privacy`, and `/legal/terms`. Keep the copy and routes aligned with support surfaces only. Add an API client placeholder that can later call the Go API.

**Acceptance Criteria:** Web app builds, all stub routes render, and the web app does not include authenticated dashboard or core tip flow UI.

**Out of Scope:** Production SEO polish and complete UI design.

---

## T00.08 — Setup Flutter App

**Stack:** mobile

**Feasibility:** Medium.

**Depends on:** T00.01, T00.03.

**Developer Detail:** Create the Flutter app in `apps/mobile`. Add an app shell with routing, theme, environment config, and stub screens for splash, OAuth login, home, public profile setup, payment profile verification, donor tip flow, worker dashboard, leaderboard, discover, and QR scanner.

**Acceptance Criteria:** Flutter app runs locally, navigation between stub screens works, environment values are accessible, and future feature screens have obvious locations.

**Out of Scope:** OAuth SDK integration and real API calls.

---

# Phase 1 — OAuth Login

## T01.01 — OAuth User Schema

**Stack:** database/backend

**Feasibility:** Medium.

**Depends on:** T00.06.

**Developer Detail:** Add `users`, `oauth_accounts`, `refresh_tokens`, `notification_tokens`, and `audit_logs`. `oauth_accounts` must enforce unique `(provider, provider_user_id)`. Users may have nullable email because Facebook may not return one. Store refresh token hashes, not raw refresh tokens.

**Acceptance Criteria:** Migrations apply cleanly, unique constraints prevent duplicate provider accounts, and no OTP-login tables are introduced.

**Out of Scope:** Worker profile and payment profile tables.

---

## T01.02 — Google OAuth Login API

**Stack:** backend/mobile

**Feasibility:** Medium.

**Depends on:** T01.01, T01.04.

**Developer Detail:** Implement `POST /api/auth/oauth/google`. Accept the credential generated by Flutter, verify it with Google, map provider identity to a local user, link or update the OAuth account, issue app-owned access/refresh tokens, and write an audit log. Handle duplicate email carefully: provider identity is the authority, not email alone.

**Acceptance Criteria:** Valid Google login creates or reuses a user, invalid credentials fail with 401, response includes app session tokens and user summary, and audit log records login success/failure.

**Out of Scope:** Flutter button implementation.

---

## T01.03 — Facebook OAuth Login API

**Stack:** backend/mobile

**Feasibility:** Medium-high because Facebook provider setup can be slower.

**Depends on:** T01.01, T01.04.

**Developer Detail:** Implement `POST /api/auth/oauth/facebook`. Verify the token with Facebook, fetch provider user ID, display name, avatar, and email when available. Account linking must work without email. Add feature flag support if Facebook is not launch-ready.

**Acceptance Criteria:** Valid Facebook login creates or reuses a user, missing email does not break account creation, invalid token fails safely, and response shape matches Google login.

**Out of Scope:** Facebook app review and mobile SDK setup.

---

## T01.04 — App JWT Session System

**Stack:** backend/security

**Feasibility:** Medium.

**Depends on:** T01.01.

**Developer Detail:** Add access token creation, refresh token creation, refresh endpoint, token validation helpers, and optional token rotation. Use short-lived access tokens and longer refresh tokens. Store refresh token hashes with expiry, revoked timestamp, device metadata, and user ID.

**Acceptance Criteria:** Access tokens authorize requests, refresh endpoint rotates or renews tokens according to documented behavior, revoked/expired refresh tokens fail, and token claims include user ID and session ID.

**Out of Scope:** OAuth provider verification.

---

## T01.05 — Auth Middleware

**Stack:** backend/security

**Feasibility:** Easy-medium.

**Depends on:** T01.04.

**Developer Detail:** Add `RequireAuth`, request context user loading, `RequireWorker` helper, and ownership guard patterns. Middleware should return consistent JSON errors and never leak token parsing internals.

**Acceptance Criteria:** Protected test route rejects missing/invalid tokens, accepts valid tokens, attaches user context, and ownership helper can be reused by tip/profile routes.

**Out of Scope:** Full role/permission system beyond MVP needs.

---

## T01.06 — Flutter Google and Facebook Login UI

**Stack:** mobile

**Feasibility:** Medium-high due to mobile provider config.

**Depends on:** T00.08, T01.02, T01.03, T02.06.

**Developer Detail:** Add Google and Facebook login buttons, provider SDK setup, backend token exchange, secure app token storage, logout, loading/error states, and onboarding routing based on `GET /api/me`. Keep provider sessions separate from app tokens.

**Acceptance Criteria:** User can sign in with Google, Facebook path is either functional or feature-flagged, tokens are stored securely, logout clears app tokens, and user lands on the correct onboarding screen.

**Out of Scope:** Payment profile verification UI.

---

# Phase 2 — Worker Profile and Payment Profile Verification

## T02.01 — Profile and Payment Profile Schema

**Stack:** database/backend

**Feasibility:** Medium.

**Depends on:** T01.01.

**Developer Detail:** Add `worker_profiles`, `payment_profiles`, and `payment_profile_otp_codes`. Keep public profile separate from payment eligibility. Payment profile should include verification status, verified timestamp, verification method, normalized PromptPay ID, promptpay type, and last change timestamp. OTP rows store hashes, expiry, attempts, consumed timestamp, and destination.

**Acceptance Criteria:** Schema supports unverified and verified payment profiles, changing payment data can reset verification, and public profile completion does not imply payment eligibility.

**Out of Scope:** Workplace membership and stats.

---

## T02.02 — Public Worker Profile API

**Stack:** backend/mobile/frontend

**Feasibility:** Medium.

**Depends on:** T01.05, T02.01.

**Developer Detail:** Implement authenticated create/update profile behavior at `PATCH /api/profile`. Validate username format and uniqueness, sanitize bio fields, store display name/avatar/job title, and expose a read endpoint for public profile lookup if needed by Flutter/web.

**Acceptance Criteria:** Authenticated user can create/update their profile, username conflicts return 409, public-safe profile data can be fetched, and payment details are never included in public responses.

**Out of Scope:** Avatar upload processing unless already available.

---

## T02.03 — Payment Profile API

**Stack:** backend/mobile

**Feasibility:** Medium.

**Depends on:** T01.05, T02.01.

**Developer Detail:** Implement `PATCH /api/payment-profile`. Accept PromptPay ID, type, and phone if phone verification is used. Normalize phone/PromptPay values. If payment data changes after verification, reset status to unverified and record an audit log. Return masked payment data in responses.

**Acceptance Criteria:** Worker can save payment profile, invalid PromptPay formats are rejected, verification resets on change, response masks sensitive fields, and `can_receive_tips` remains false until OTP verification.

**Out of Scope:** Bank account payout support.

---

## T02.04 — Send Payment Profile OTP API

**Stack:** backend/security

**Feasibility:** Medium-high because SMS provider choice affects production.

**Depends on:** T02.03.

**Developer Detail:** Implement `POST /api/payment-profile/otp/send`. Require auth, ensure payment profile exists, generate OTP, hash it, store expiry, rate limit by user and phone, invalidate older active OTPs if appropriate, and log OTP only in dev mode. Production should call an SMS provider behind an interface.

**Acceptance Criteria:** OTP can be requested for a saved payment profile, rate limits block abuse, dev mode exposes enough logs for testing, production path does not log raw OTP, and audit logs capture sends.

**Out of Scope:** Selecting final paid SMS vendor unless product owner decides now.

---

## T02.05 — Verify Payment Profile OTP API

**Stack:** backend/security

**Feasibility:** Medium.

**Depends on:** T02.04.

**Developer Detail:** Implement `POST /api/payment-profile/otp/verify`. Compare OTP against latest unconsumed hash, enforce expiry and attempt limits, mark OTP consumed, set payment profile verified, publish `payment_profile.verified` if events exist, and write an audit log.

**Acceptance Criteria:** Correct OTP verifies profile, expired/wrong/consumed OTP fails, repeated wrong attempts are limited, verification sets `can_receive_tips` true, and changing payment info afterward resets it.

**Out of Scope:** Identity document verification.

---

## T02.06 — Current User and Onboarding State API

**Stack:** backend/mobile

**Feasibility:** Easy-medium.

**Depends on:** T01.05, T02.01.

**Developer Detail:** Implement `GET /api/me` with user identity, OAuth providers, public profile status, payment profile status, and `can_receive_tips`. Keep this endpoint as the single source for Flutter routing.

**Acceptance Criteria:** Authenticated request returns complete onboarding state, missing profile/payment profile is represented explicitly, and no sensitive payment values are exposed.

**Out of Scope:** Admin user management.

---

## T02.07 — Flutter Public Profile Setup

**Stack:** mobile

**Feasibility:** Medium.

**Depends on:** T01.06, T02.02, T02.06.

**Developer Detail:** Build profile form with username, display name, avatar URL or placeholder, bio, and job title. Add validation, save state, conflict handling, loading/error states, and route onward based on `/api/me`.

**Acceptance Criteria:** User can create/update profile from Flutter, username conflict is understandable, successful save updates onboarding state, and profile data appears in donor-facing preview.

**Out of Scope:** Rich avatar upload unless upload APIs are ready.

---

## T02.08 — Flutter Payment Profile Verification

**Stack:** mobile

**Feasibility:** Medium-high due to OTP UX.

**Depends on:** T02.03, T02.04, T02.05, T02.06.

**Developer Detail:** Build payment profile form, save action, OTP request, OTP entry, resend countdown, verification result, and verified state. Show clear messaging that public profile and payment eligibility are separate.

**Acceptance Criteria:** Worker can enter PromptPay info, request OTP, verify OTP, see verified state, and app prevents receive-tip flow until verification succeeds.

**Out of Scope:** Multiple payout methods.

---

# Phase 3 — Flutter Core Tip Flow

## T03.01 — Tip Request Schema

**Stack:** database/backend

**Feasibility:** Medium.

**Depends on:** T02.01.

**Developer Detail:** Add `tip_requests` with worker user/profile reference, optional donor user reference, amount, currency, ref code, status, slip object key, timestamps, expiry timestamps, and audit fields. Enforce valid status values and indexes for worker/status queries.

**Acceptance Criteria:** Schema stores pending, slip_uploaded, confirmed, disputed, and expired states; ref code uniqueness is enforced; and status query indexes exist.

**Out of Scope:** Automatic bank statement reconciliation.

---

## T03.02 — Initiate Tip API

**Stack:** backend/mobile

**Feasibility:** Medium.

**Depends on:** T02.02, T02.03, T02.05, T03.01.

**Developer Detail:** Implement `POST /api/tips/initiate`. Accept worker identifier and amount. Validate worker exists, public profile is active, payment profile is verified, amount is within min/max, and ref code is unique. Only return PromptPay payment destination after eligibility passes.

**Acceptance Criteria:** Eligible worker returns tip ID, amount, ref code, expiry, and payment data; unverified worker returns a safe error without payment data; invalid amount fails; duplicate ref code cannot be created.

**Out of Scope:** Charging or collecting money directly.

---

## T03.03 — PromptPay QR Helper

**Stack:** backend/mobile

**Feasibility:** Medium.

**Depends on:** T03.02.

**Developer Detail:** Decide where QR payload is generated. Recommended MVP: backend returns normalized PromptPay payload fields plus generated EMV payload string; Flutter renders QR from that payload. This centralizes correctness and keeps payment destination logic server-side.

**Acceptance Criteria:** QR screen can render a scannable PromptPay QR with amount and destination, ref code instruction is visible, and backend tests cover payload generation for phone PromptPay.

**Out of Scope:** Dynamic QR integration with bank APIs.

---

## T03.04 — S3 Presigned Upload API

**Stack:** backend/mobile/storage

**Feasibility:** Medium.

**Depends on:** T00.05, T01.05.

**Developer Detail:** Implement `POST /api/uploads/presign`. Require auth where applicable, accept content type and size metadata, restrict to image MIME types, create random object keys, set short expiry, and return upload URL plus object key. Consider a `slips/` prefix and private bucket policy.

**Acceptance Criteria:** Flutter can upload a JPEG/PNG directly to local MinIO/S3, oversized or non-image uploads are rejected by policy/API, and arbitrary external URLs are never accepted as slip evidence.

**Out of Scope:** Image moderation and OCR.

---

## T03.05 — Attach Slip to Tip API

**Stack:** backend/mobile

**Feasibility:** Medium.

**Depends on:** T03.01, T03.04, T04.02 if events are ready.

**Developer Detail:** Implement `POST /api/tips/:id/slip`. Validate tip exists, status is pending, file key belongs to allowed upload prefix, and caller is allowed to attach the slip. Move status to `slip_uploaded`, store key/timestamp, and publish `tip.slip_uploaded` when RabbitMQ is available.

**Acceptance Criteria:** Pending tip accepts one valid slip, non-pending tips reject slip changes, invalid object keys fail, and worker dashboard can see the new pending confirmation item.

**Out of Scope:** Slip amount validation by OCR.

---

## T03.06 — Worker Pending Tips API

**Stack:** backend/mobile

**Feasibility:** Easy-medium.

**Depends on:** T01.05, T03.05.

**Developer Detail:** Implement `GET /api/tips?status=slip_uploaded`. Require auth, scope results to current worker, support pagination, sort newest first, and include short-lived slip preview URLs.

**Acceptance Criteria:** Worker sees only their own slip_uploaded tips, unrelated users cannot access them, response includes amount/ref code/timestamps/slip preview URL, and pagination works.

**Out of Scope:** Admin dispute queue.

---

## T03.07 — Confirm Tip API

**Stack:** backend/mobile

**Feasibility:** Medium.

**Depends on:** T03.05, T04.02 if events are ready.

**Developer Detail:** Implement `PATCH /api/tips/:id/confirm`. Require auth, verify worker ownership, enforce `slip_uploaded` as the only source status, update atomically, write audit log, and publish `tip.confirmed`.

**Acceptance Criteria:** Owner can confirm slip_uploaded tip, non-owner receives forbidden/not found, repeated confirm is idempotent or safely rejected, and stats update can later consume the event.

**Out of Scope:** Manual financial settlement.

---

## T03.08 — Dispute Tip API

**Stack:** backend/mobile

**Feasibility:** Medium.

**Depends on:** T03.05.

**Developer Detail:** Implement `PATCH /api/tips/:id/dispute` with optional dispute reason. Require worker ownership and `slip_uploaded` status. Store disputed timestamp/reason and publish `tip.disputed`.

**Acceptance Criteria:** Worker can dispute invalid slip, invalid status transitions fail, donor status endpoint shows disputed, and audit log records the action.

**Out of Scope:** Full support case workflow.

---

## T03.09 — Tip Status API

**Stack:** backend/mobile

**Feasibility:** Easy-medium.

**Depends on:** T03.01.

**Developer Detail:** Implement `GET /api/tips/:id/status`. Recommended MVP: use an opaque status token returned by initiate, not a public predictable ID alone. Return only donor-safe state: status, amount, worker display info, ref code, timestamps, and terminal outcome.

**Acceptance Criteria:** Donor flow can poll status, unauthorized users cannot enumerate private tip data, terminal states stop polling, and response excludes worker payment details after initiation.

**Out of Scope:** Real-time websocket updates.

---

## T03.10 — Flutter Donor Tip Flow

**Stack:** mobile

**Feasibility:** Hard but core MVP.

**Depends on:** T02.02, T03.02, T03.03, T03.04, T03.05, T03.09.

**Developer Detail:** Build donor flow: open worker profile from QR/link/discover, show profile, choose amount, initiate tip, render PromptPay QR, show ref code instruction, upload slip, attach slip, show status, and poll until confirmed/disputed/expired. Handle errors for unverified worker and expired tip.

**Acceptance Criteria:** A donor can complete the full happy path in Flutter against local API, upload slip successfully, see pending confirmation, and eventually see confirmed/disputed state.

**Out of Scope:** Web donor flow and direct in-app payments.

---

## T03.11 — Flutter Worker Dashboard MVP

**Stack:** mobile

**Feasibility:** Medium-high.

**Depends on:** T03.06, T03.07, T03.08.

**Developer Detail:** Build worker dashboard list for pending slip confirmations. Show amount, ref code, donor label if available, uploaded time, slip preview, confirm action, dispute action, empty state, and refresh behavior. Disable buttons while requests are in flight.

**Acceptance Criteria:** Worker can view pending slips, open preview, confirm/dispute, list refreshes after action, and dashboard never shows other workers' tips.

**Out of Scope:** Advanced analytics dashboard.

---

# Phase 4 — Queue and Worker Service

## T04.01 — Worker Service Bootstrap

**Stack:** worker

**Feasibility:** Medium.

**Depends on:** T00.05.

**Developer Detail:** Create Go worker service in `services/worker` with config loading, logging, RabbitMQ connection, graceful shutdown, and a sample consumer. Keep job handlers separated by domain: notifications, expiry, leaderboard.

**Acceptance Criteria:** Worker starts locally, connects to RabbitMQ, consumes a test message, logs processing result, and shuts down cleanly.

**Out of Scope:** Actual FCM send and leaderboard jobs.

---

## T04.02 — Event Contract

**Stack:** backend/worker/docs

**Feasibility:** Easy.

**Depends on:** T04.01.

**Developer Detail:** Create `docs/EVENTS.md`. Define event names, exchange/queue naming, routing keys, JSON payload schema, version field, idempotency key, timestamp, producer, consumer, and retry/dead-letter behavior.

**Acceptance Criteria:** API and worker agree on event names for `tip.created`, `tip.slip_uploaded`, `tip.confirmed`, `tip.disputed`, `tip.expired`, `notification.send`, `leaderboard.refresh`, and `payment_profile.verified`.

**Out of Scope:** Implementing every event producer/consumer.

---

## T04.03 — Expire Tips Cron

**Stack:** worker/backend

**Feasibility:** Medium.

**Depends on:** T03.01, T04.01, T04.02.

**Developer Detail:** Add scheduled worker job to expire old tips. Pending tips older than 24h become expired. Slip_uploaded tips older than 72h become expired. Use database-side conditional updates to avoid races with confirm/dispute actions.

**Acceptance Criteria:** Expiry job updates only eligible rows, publishes `tip.expired`, is idempotent across reruns, and tests cover boundary times.

**Out of Scope:** User-configurable expiry windows.

---

# Phase 5 — Notifications

## T05.01 — Register FCM Token API

**Stack:** backend/mobile

**Feasibility:** Easy-medium.

**Depends on:** T01.05, T01.01.

**Developer Detail:** Implement `POST /api/notification-tokens`. Store token, platform, app version, device ID if available, user ID, enabled status, and last seen timestamp. Upsert by token or device identifier.

**Acceptance Criteria:** Authenticated Flutter app can register/update token, duplicate registrations do not create noisy duplicates, and logout or disable path can deactivate token.

**Out of Scope:** Push sending.

---

## T05.02 — Flutter Push Handling

**Stack:** mobile

**Feasibility:** Medium-high due to platform setup.

**Depends on:** T05.01.

**Developer Detail:** Configure Firebase Messaging for Android/iOS, request permissions, retrieve FCM token, send token to API, handle foreground/background/tap events, and route slip-upload notifications to worker dashboard or tip detail.

**Acceptance Criteria:** Device token registers, foreground notification handling works, tapping notification opens relevant screen, and permission-denied state is handled gracefully.

**Out of Scope:** LINE notifications.

---

## T05.03 — FCM Notification Worker

**Stack:** worker/notification

**Feasibility:** Medium.

**Depends on:** T04.01, T04.02, T05.01.

**Developer Detail:** Consume `tip.slip_uploaded` or `notification.send`, load worker notification tokens, send FCM push, record send result, and handle invalid tokens by deactivating them. Use retry rules for transient provider errors.

**Acceptance Criteria:** Slip upload triggers worker push in local/staging config, invalid tokens are handled, failures are logged with enough context, and duplicate events do not spam users if idempotency key repeats.

**Out of Scope:** Notification preferences beyond MVP.

---

## T05.04 — LINE Messaging Optional

**Stack:** worker/notification

**Feasibility:** Optional-medium.

**Depends on:** T04.02 and product decision.

**Developer Detail:** Keep LINE behind a provider interface compatible with notification events. Do not block MVP on this. Document required LINE credentials and user linking approach only if moved forward.

**Acceptance Criteria:** If implemented, LINE provider can be enabled by config and disabled without affecting FCM.

**Out of Scope:** MVP launch requirement.

---

# Phase 6 — Leaderboard and Web Support

## T06.01 — Worker Stats Schema

**Stack:** database/backend

**Feasibility:** Medium.

**Depends on:** T03.01.

**Developer Detail:** Add aggregate table such as `worker_stats` with worker ID, confirmed tip count, confirmed amount total, last confirmed at, ranking period fields if needed, and updated timestamp. Decide whether MVP uses all-time ranking only.

**Acceptance Criteria:** Stats table can support worker dashboard and leaderboard queries, has unique row per worker/period, and can be rebuilt from confirmed tips if needed.

**Out of Scope:** Advanced analytics and fraud scoring.

---

## T06.02 — Update Stats on Confirmed Tip

**Stack:** backend/worker

**Feasibility:** Medium.

**Depends on:** T03.07, T06.01.

**Developer Detail:** On `tip.confirmed`, update aggregate stats idempotently. Either update in the confirm transaction or consume event with processed-event tracking. Avoid double-counting repeated events.

**Acceptance Criteria:** Confirming a tip increments count and amount exactly once, repeated confirm/event does not double count, and stats can be recalculated for repair.

**Out of Scope:** Real-time ranking cache.

---

## T06.03 — Leaderboard API

**Stack:** backend/mobile/frontend

**Feasibility:** Easy-medium.

**Depends on:** T06.01, T06.02.

**Developer Detail:** Implement `GET /api/leaderboard` with pagination, period parameter if supported, public worker profile fields, rank, total tips, total amount if product wants it public, and current user's rank if authenticated.

**Acceptance Criteria:** API returns ranked workers, supports limit/page, excludes unlisted/incomplete profiles, and response is safe for public web consumption.

**Out of Scope:** Personalized recommendations.

---

## T06.04 — Flutter Leaderboard Screen

**Stack:** mobile

**Feasibility:** Easy-medium.

**Depends on:** T06.03.

**Developer Detail:** Build leaderboard screen with ranked list, worker avatar/name/job title, count or score, loading/error/empty states, and tap into public profile/tip flow.

**Acceptance Criteria:** Flutter renders leaderboard from API, handles pagination or refresh, and tapping a worker can enter donor profile flow.

**Out of Scope:** Complex filters.

---

## T06.05 — Next.js Leaderboard Page

**Stack:** frontend

**Feasibility:** Easy-medium.

**Depends on:** T00.07, T06.03.

**Developer Detail:** Build `/leaderboard` as public support page backed by API. Use SEO-safe rendering where practical and link workers to `/@username`. Keep it unauthenticated.

**Acceptance Criteria:** Page builds, fetches or statically renders leaderboard, has loading/error fallback, and does not include core tip flow unless product decision changes.

**Out of Scope:** Authenticated dashboard.

---

# Phase 7 — Discover and QR

## T07.01 — Flutter QR Scanner

**Stack:** mobile

**Feasibility:** Medium.

**Depends on:** T03.10.

**Developer Detail:** Add camera permission flow, QR scanner package, supported QR format, parsing for worker profile/deep link, invalid QR state, and routing into donor tip flow.

**Acceptance Criteria:** User can scan a valid worker QR and land on donor profile flow, invalid QR shows safe error, and camera permissions are handled on Android/iOS.

**Out of Scope:** Generating printed QR assets unless requested.

---

## T07.02 — Flutter Discover Screen

**Stack:** mobile

**Feasibility:** Medium.

**Depends on:** T02.02, T03.10.

**Developer Detail:** Add discover API if not already available, then build Flutter browse/search screen for public worker profiles. MVP can use simple search by username/display name and recent/featured ordering.

**Acceptance Criteria:** User can browse/search workers, open worker profile, and start donor tip flow for eligible workers.

**Out of Scope:** Zone-based discover and complex ranking.

---

## T07.03 — Next.js Discover Page Optional

**Stack:** frontend

**Feasibility:** Optional-medium.

**Depends on:** T07.02 API or public profile API.

**Developer Detail:** Build public `/discover` for SEO/shareability only after Flutter discover is stable. Use same public profile fields and link to `/@username`.

**Acceptance Criteria:** Public discover page renders indexed worker profiles without exposing payment data.

**Out of Scope:** Web core tip flow.

---

# Phase 8 — Workplace Optional

## T08.01 — Workplace Schema

**Stack:** database/backend

**Feasibility:** Optional-medium.

**Depends on:** T02.01.

**Developer Detail:** Design `workplaces` and `worker_workplaces` with active membership constraint. Keep it optional and avoid coupling core tip flow to workplace in MVP.

**Acceptance Criteria:** Schema supports future workplace association and enforces one active workplace per worker.

**Out of Scope:** MVP launch.

---

## T08.02 — Workplace APIs

**Stack:** backend/mobile

**Feasibility:** Optional-medium.

**Depends on:** T08.01.

**Developer Detail:** Add APIs to create workplace, join workplace, and activate workplace. Include ownership/admin rules before production use.

**Acceptance Criteria:** Worker can join and activate one workplace, API rejects multiple active workplaces, and workplace changes do not break existing worker profiles.

**Out of Scope:** Workplace admin dashboard and payroll features.

---

# Phase 9 — Testing

## T09.01 — Backend Unit Tests

**Stack:** backend/testing

**Feasibility:** Medium.

**Depends on:** Backend features from Phases 1-6.

**Developer Detail:** Add focused tests for OAuth linking, JWT validation, payment profile OTP, tip status transitions, ownership checks, expiry logic, and stats idempotency. Use table-driven tests where possible.

**Acceptance Criteria:** Unit test suite runs locally and in CI, covers success and failure paths, and money-adjacent state transitions are tested.

**Out of Scope:** Full end-to-end mobile automation.

---

## T09.02 — Backend Integration Tests

**Stack:** backend/testing

**Feasibility:** Medium-high.

**Depends on:** Docker infra and completed backend MVP flow.

**Developer Detail:** Add integration test that runs against test PostgreSQL/Redis/RabbitMQ/S3-compatible storage. Scenario: OAuth login stub, create public profile, create payment profile, verify OTP, initiate tip, upload/attach slip, confirm tip, leaderboard updates.

**Acceptance Criteria:** One command runs full backend happy path, test data is isolated, and failures show which step broke.

**Out of Scope:** Real provider OAuth/SMS calls.

---

## T09.03 — Flutter Smoke Tests

**Stack:** mobile/testing

**Feasibility:** Medium.

**Depends on:** Flutter MVP screens.

**Developer Detail:** Add smoke/widget tests for login UI, profile setup, payment profile OTP screens, donor tip flow shell, and worker dashboard. Mock API client where needed.

**Acceptance Criteria:** `flutter analyze` and test command pass, critical screens render, and basic navigation does not regress.

**Out of Scope:** Device farm testing.

---

## T09.04 — Web Smoke Tests

**Stack:** frontend/testing

**Feasibility:** Easy-medium.

**Depends on:** Next.js support pages.

**Developer Detail:** Add build and simple route tests for leaderboard, discover, and public worker profile preview. Use mocked API responses if backend is unavailable.

**Acceptance Criteria:** Web build passes, support routes render, and public pages do not require authentication.

**Out of Scope:** Core tip flow tests on web.

---

# Phase 10 — Deployment

## T10.01 — Docker Build API and Worker

**Stack:** deployment/backend/worker

**Feasibility:** Medium.

**Depends on:** T00.04, T04.01.

**Developer Detail:** Add production Dockerfiles for API and worker using multi-stage builds. Run as non-root, include healthcheck where practical, and keep images small.

**Acceptance Criteria:** API and worker images build locally, start with env config, expose expected ports/commands, and do not include source-only dev artifacts unnecessarily.

**Out of Scope:** Cloud deployment.

---

## T10.02 — CI Pipeline

**Stack:** ci

**Feasibility:** Medium.

**Depends on:** Test commands and app skeletons.

**Developer Detail:** Add CI for backend tests, Flutter analyze/test, web build/test, and Docker build. Cache dependencies but keep secrets out of CI logs.

**Acceptance Criteria:** Pull requests run checks automatically, failing tests block merge by policy if branch protection is enabled, and CI documents required secrets.

**Out of Scope:** Release automation to app stores.

---

## T10.03 — Staging Plan

**Stack:** deployment/docs

**Feasibility:** Medium.

**Depends on:** T10.01, T10.02.

**Developer Detail:** Create `docs/DEPLOYMENT.md` covering staging topology, environment variables, database migrations, S3 bucket setup, RabbitMQ/Redis, FCM credentials, OAuth redirect/app config, rollback, and smoke test checklist.

**Acceptance Criteria:** A developer can deploy staging by following the doc, known manual provider setup steps are explicit, and rollback path is documented.

**Out of Scope:** Production incident runbook.

---

## T10.04 — Production Checklist

**Stack:** deployment/security/ops

**Feasibility:** Medium-high.

**Depends on:** Staging plan and MVP feature completion.

**Developer Detail:** Create production readiness checklist covering DB backups, S3 lifecycle, error tracking, structured logs, metrics, secret management, rate limits, manual dispute process, provider credentials, privacy/terms, data retention, and alerting.

**Acceptance Criteria:** Checklist exists, each item has owner/status, and production launch is blocked until critical security and operations items are complete.

**Out of Scope:** Implementing every ops integration in this ticket.

---

# Recommended MVP V2

## Keep

1. Flutter app foundation.
2. Google OAuth login.
3. Facebook OAuth login only if provider setup is ready; otherwise feature-flag it.
4. Public worker profile.
5. Payment profile plus OTP verification.
6. Flutter donor tip flow.
7. Flutter worker dashboard.
8. S3-compatible slip upload.
9. FCM notification for uploaded slips.
10. Basic leaderboard API.
11. Next.js leaderboard and public profile preview.

## Postpone

1. LINE Messaging API.
2. Workplace system.
3. Next.js authenticated dashboard.
4. Next.js core tip flow.
5. Advanced analytics.
6. Zone-based discover.
7. OCR-based slip validation.
8. Direct bank integration.

---

# Open Decisions

1. **PromptPay QR payload:** Recommendation for MVP is backend-side EMV payload generation, Flutter-side QR rendering.
2. **Payment profile MVP:** Recommendation is phone-based PromptPay first. Add national ID/e-wallet variants only after initial flow is stable.
3. **Donor status page:** Recommendation is signed or opaque status token, not public predictable tip ID.
4. **Web public worker page:** Recommendation is preview plus app/deep-link redirect for MVP. Do not build web tip flow yet.
5. **Facebook login:** Recommendation is feature flag until provider setup and mobile review are confirmed.
6. **Slip fraud handling:** Recommendation is manual worker confirmation for MVP, with audit logs and dispute state. OCR can come later.

---

# Suggested Build Order

1. Finish Phase 0 foundation.
2. Build Phase 1 auth with Google first, Facebook behind flag if needed.
3. Build Phase 2 profile/payment verification.
4. Build Phase 3 core tip flow end to end with local S3.
5. Add Phase 4 events and worker only where needed for notifications/expiry.
6. Add Phase 5 FCM notification.
7. Add basic Phase 6 leaderboard.
8. Add tests and deployment hardening before production.
