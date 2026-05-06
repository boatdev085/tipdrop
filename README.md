# TipDrop

TipDrop is a Flutter-first tipping app with a Go Fiber API, PostgreSQL, Redis, RabbitMQ, S3-compatible storage, a Go background worker, and a small Next.js support web app.

## Repository Layout

- `apps/mobile` - Flutter primary app
- `apps/web` - Next.js support web app
- `services/api` - Go Fiber API
- `services/worker` - Go background worker
- `infra/docker` - local Docker Compose infrastructure
- `infra/migrations` - SQL database migrations
- `docs` - developer documentation
- `issues` - planning and task documents

## Local Development

1. Copy `.env.example` to `.env` and fill local values.
2. Start local dependencies:

```powershell
docker compose -f infra/docker/docker-compose.yml up -d
```

3. Run the API once Go is installed:

```powershell
cd services/api
go run ./cmd/api
```

4. Run the worker once Go is installed:

```powershell
cd services/worker
go run ./cmd/worker
```

5. Run the web app:

```powershell
cd apps/web
npm install
npm run dev
```

6. Generate Flutter platform folders and run the mobile app once Flutter is installed:

```powershell
cd apps/mobile
flutter create .
flutter pub get
flutter run
```

Flutter is the core product surface. Next.js is only for public/support pages such as leaderboard, discover, public profile preview, SEO, privacy, and terms.
