# AGENTS.md - Frontend Web

## Scope
These instructions apply to the Next.js support web app in `apps/web`.

## Goal
Provide lightweight support pages and WebView-compatible UI that consumes backend contracts without becoming the source of truth.

## Stack and Libraries
- Framework/runtime: Next.js 16.2.4 with React 19.2.4 and React DOM 19.2.4.
- Language/tooling: TypeScript 6.0.3 with `@types/node` 25.2.3, `@types/react` 19.2.14, and `@types/react-dom` 19.2.3.
- Package manager/workspace: npm workspace `apps/web`; run scripts through `npm --workspace apps/web ...` from the repository root.
- API access: keep backend calls in `lib/api.ts` or clearly named helper modules and use the browser/Next.js `fetch` stack unless a task explicitly adds another client library.

## API Contract
- Use documented backend API responses from `services/api`; do not duplicate business rules in the UI.
- Keep API client changes in `lib/api.ts` or clearly named helpers.
- Handle loading, empty, and error states for every backend request.
- Never collect or persist sensitive payment data beyond what the backend contract explicitly requires.
- Slip uploads must follow the signed S3 upload flow provided by the backend.

## DB Changes
Frontend tasks should not make direct database changes. If a UI feature requires data that is unavailable, document the backend API and DB changes needed instead of bypassing the API.

## Files Changed
For each task, summarize changed pages, components, styles, API client code, and any backend contract dependency.

## Test Coverage
- Run `npm --workspace apps/web run lint` for web changes when dependencies are installed.
- Run `npm --workspace apps/web run build` before completion when practical.
- Add component or integration coverage when a test framework is present for the affected area.

## Output Used By Next Task
Leave route names, required backend fields, and UX state notes so Flutter or backend follow-up work can reuse the same contract.
