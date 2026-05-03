# AGENTS.md

## Project Overview
TipDrop is a Flutter + Next.js (WebView) + Go Fiber platform for direct PromptPay tipping.

## Core Principles
1. Workers receive 100% of tips.
2. Platform never holds money.
3. Backend is the source of truth.
4. Tip state transitions must follow defined status machine.
5. All slip images must go through signed S3 upload.

## Stack
- Flutter (mobile shell + native)
- Next.js (WebView UI)
- Go Fiber (API)
- PostgreSQL (DB)
- Redis (cache)
- RabbitMQ (queue)
- S3 (storage)

## Task Requirements
Every task MUST include:
- Clear goal
- API contract (if applicable)
- DB changes (if applicable)
- Files changed
- Test coverage
- Output used by next task

## Forbidden
- Do not implement wallet system
- Do not bypass backend validation
- Do not store sensitive data without encryption

## Definition of Done
- Code builds
- Tests pass
- API documented
- No security violation
