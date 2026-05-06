.PHONY: test test-go test-web build-web docker-up docker-down

test: test-go test-web

test-go:
	go test ./services/api/... ./services/worker/...

test-web:
	npm --workspace apps/web run lint

build-web:
	npm --workspace apps/web run build

docker-up:
	docker compose -f infra/docker/docker-compose.yml up -d

docker-down:
	docker compose -f infra/docker/docker-compose.yml down
