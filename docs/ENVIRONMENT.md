# Environment Contract

This file documents the environment variables used by TipDrop local development.

## Shared

- `APP_ENV` - `development`, `staging`, or `production`.
- `PUBLIC_WEB_URL` - public web URL for links and profile previews.
- `API_BASE_URL` - API URL used by local clients.

## API

- `API_HOST` - bind host for the Go Fiber API.
- `API_PORT` - bind port for the Go Fiber API.
- `JWT_ACCESS_SECRET` - app access-token signing secret.
- `JWT_REFRESH_SECRET` - app refresh-token signing secret.

## Data Services

- `DATABASE_URL` - PostgreSQL connection string.
- `REDIS_URL` - Redis connection string.
- `RABBITMQ_URL` - RabbitMQ AMQP connection string.

## Storage

- `S3_ENDPOINT` - S3-compatible endpoint. Local default is MinIO.
- `S3_REGION` - S3 region.
- `S3_BUCKET` - bucket for private slip uploads.
- `S3_ACCESS_KEY_ID` - S3 access key.
- `S3_SECRET_ACCESS_KEY` - S3 secret key.
- `S3_FORCE_PATH_STYLE` - set `true` for MinIO.

## OAuth

- `GOOGLE_OAUTH_CLIENT_ID` - Google OAuth client ID.
- `FACEBOOK_OAUTH_APP_ID` - Facebook app ID.
- `FACEBOOK_OAUTH_APP_SECRET` - Facebook app secret.
- `FACEBOOK_LOGIN_ENABLED` - feature flag for Facebook login.

## OTP/SMS

- `OTP_DEV_LOG` - when true, development OTP codes may be logged.
- `SMS_PROVIDER` - SMS provider identifier. Use `dev` locally.
- `SMS_API_KEY` - SMS provider API key.

## Notifications

- `FCM_PROJECT_ID` - Firebase project ID.
- `FCM_CREDENTIALS_JSON` - service account JSON or path, depending on deployment.

## Mobile

- `FLUTTER_APP_SCHEME` - deep link scheme for the Flutter app.

## Web

- `NEXT_PUBLIC_API_BASE_URL` - public API base URL for Next.js.

Never commit real secrets.
