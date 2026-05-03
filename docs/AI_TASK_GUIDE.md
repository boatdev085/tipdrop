# AI Task Guide

## How AI agents should work on TipDrop tasks

Every implementation task must include stack-specific context so an AI agent understands the required tools, boundaries, and expected output.

## Required issue sections

Each issue should contain:

1. Goal
2. Stack Context
3. Dependencies
4. Implementation Steps
5. API Contract, if applicable
6. Database Changes, if applicable
7. Files to Change
8. Edge Cases
9. Output for Next Task
10. Definition of Done

## Stack Context Format

Use this format inside every issue:

```md
## Stack Context
- Stack: backend | frontend | mobile | infra | worker
- Language/Framework: Go Fiber | Next.js | Flutter | PostgreSQL | Redis | RabbitMQ | S3
- Required knowledge:
  - Authentication/session handling
  - API contract discipline
  - Error handling
  - Testing expectations
- Must follow:
  - AGENTS.md
  - docs/ARCHITECTURE.md
```

## Status Comment Convention

Use comments on each issue to track progress:

```txt
[STARTED]
Owner:
Branch:
Current focus:

[BLOCKED]
Reason:
Needed decision:

[READY FOR REVIEW]
Summary:
Test evidence:

[DONE]
Merged PR:
Output for next task:
```

## AI Guardrails

- Do not invent APIs that conflict with docs/API_CONTRACT.md.
- Do not add a wallet table.
- Do not let frontend confirm tip state without backend authorization.
- Do not store slip images outside signed S3 upload flow.
- Do not merge changes without updating docs when behavior changes.
