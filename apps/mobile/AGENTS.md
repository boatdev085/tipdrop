# AGENTS.md - Flutter Mobile

## Scope
These instructions apply to the Flutter mobile app in `apps/mobile`.

## Goal
Keep Flutter as the primary TipDrop product surface while relying on the backend for authentication, tip state, signed uploads, and payment-profile validation.

## Stack and Libraries
- Framework/runtime: Flutter app using Dart SDK `>=3.9.0 <4.0.0`.
- UI/platform: Flutter Material (`uses-material-design: true`) and `cupertino_icons` 1.0.9.
- Routing and networking: `go_router` 17.2.3 for navigation and `http` 1.6.0 for API calls.
- Secure local storage: `flutter_secure_storage` 10.0.0; never use it to create a client-side wallet or source-of-truth payment state.
- Mobile capabilities: `mobile_scanner` 7.2.0 for QR scanning, `qr_flutter` 4.1.0 for QR rendering, and `image_picker` 1.2.2 for user-selected slip images before backend-signed S3 upload.
- Push/Firebase: `firebase_core` 4.7.0 and `firebase_messaging` 16.2.0.
- Quality tooling: `flutter_lints` 6.0.0 and `flutter_test`.

## API Contract
- Consume documented backend API contracts; do not invent client-only source-of-truth state for tips or worker payment profiles.
- Keep environment and endpoint configuration explicit in app config files.
- Handle loading, offline, empty, success, and error states for network-backed screens.
- Use backend-provided signed S3 upload flows for slip images.
- Do not implement wallet balances, escrow, custody, or platform-held funds in the app.

## DB Changes
Flutter tasks should not make direct database changes. If app work needs new persisted data, document the required backend API and migration instead.

## Files Changed
For each task, summarize changed screens, widgets, config, assets, and any backend contract dependency.

## Test Coverage
- Run `flutter analyze` from `apps/mobile` for Flutter changes when Flutter is available.
- Run `flutter test` from `apps/mobile` when tests exist or are added.
- Add widget or unit tests for non-trivial UI state, validation, and API mapping logic when applicable.

## Output Used By Next Task
Leave screen names, expected route/deep-link behavior, required backend fields, and unresolved API gaps for the next agent.
