# TipDrop Task Master List

This document defines the full task breakdown before execution. It is intended for review, refinement, and AI analysis before creating GitHub Issues.

---

## Phase 0 — Setup

- Setup Docker (PostgreSQL, Redis, RabbitMQ, S3)
- Setup Go Fiber API
- Setup DB connection
- Setup Redis client
- Setup RabbitMQ publisher
- Healthcheck endpoints
- Setup Next.js app
- Setup Flutter app

---

## Phase 1 — Auth

- OTP Send API
- OTP Verify API
- JWT Access Token
- Refresh Token
- Auth Middleware
- Create user on first login
- Flutter Auth UI
- Secure token storage
- WebView token bridge

---

## Phase 2 — Profile

- Create worker profile API
- Update profile API
- Get public profile API
- Worker profile page (Next.js)
- Profile setup screen (Flutter)

---

## Phase 3 — Tip Flow

- Initiate tip API
- Generate ref code
- PromptPay payload
- Tip amount selector UI
- QR display page
- S3 presigned upload API
- Upload slip UI
- Attach slip to tip
- Confirm tip API
- Dispute tip API
- Tip status API
- Status page UI

---

## Phase 4 — Queue

- Publish tip events
- Setup RabbitMQ consumer
- Notification job
- Retry & DLQ
- Expire tips cron

---

## Phase 5 — Notification

- Register device token API
- Send FCM
- Send LINE message
- Receive push (Flutter)

---

## Phase 6 — Leaderboard

- Aggregate weekly tips
- Redis cache leaderboard
- Leaderboard page

---

## Phase 7 — Mobile Integration

- WebView screen
- Inject token to JS
- JS bridge
- QR scanner

---

## Phase 8 — Workplace

- Create workplace
- Join workplace
- Activate workplace
- Discover page

---

## Phase 9 — Testing

- Backend unit tests
- Integration tests
- Web E2E tests
- Mobile smoke tests

---

## Phase 10 — Deployment

- Docker build
- CI pipeline
- Staging setup
- Production setup

---

## Notes

- This file is the source of truth for task planning
- Issues should be generated only after review
- Each task will later include detailed sub-steps, API contracts, and outputs
