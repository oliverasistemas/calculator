.PHONY: dev-api dev-web test-api test-web test vet lint check build up down

dev-api:
	cd services/api && go run ./cmd

dev-web:
	cd app-web && npm run dev

test-api:
	cd services/api && go test ./...

test-web:
	cd app-web && npm test

test: test-api test-web

vet:
	cd services/api && go vet ./...

lint:
	cd services/api && golangci-lint run

check: vet lint test

build:
	docker compose build

up:
	docker compose up

down:
	docker compose down
