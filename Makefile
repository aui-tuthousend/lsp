SHELL := /bin/sh

ENV_FILE := be/.env
COMPOSE  := docker compose
API_MAIN := ./cmd/api/main.go
WORKER_MAIN := ./cmd/worker/main.go
SEED_MAIN := ./cmd/seed/main.go

.DEFAULT_GOAL := help

.PHONY: help dev up down stop restart logs ps run worker seed migrate psql redis-cli build test fmt vet tidy fe-install fe-dev fe-build fe-typecheck r2-cors clean

help: ## Tampilkan daftar target
	@echo "Target yang tersedia:"
	@grep -E '^[a-zA-Z0-9_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
		awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-12s\033[0m %s\n", $$1, $$2}'

up: ## Jalankan Postgres + Dragonfly di Docker dan tunggu siap
	$(COMPOSE) up -d --build --wait

down: ## Hentikan Postgres + Dragonfly (data tetap tersimpan)
	$(COMPOSE) down

stop: down ## Alias untuk down

restart: ## Restart container Postgres + Dragonfly
	$(COMPOSE) restart

logs: ## Ikuti log container
	$(COMPOSE) logs -f

ps: ## Tampilkan status container
	$(COMPOSE) ps

dev: ## Jalankan API (Go) + dev server FE sekaligus (infra diasumsikan hidup)
	@set -a; [ -f $(ENV_FILE) ] && . ./$(ENV_FILE); set +a; \
	trap 'kill 0' INT TERM; \
	( go -C be run $(API_MAIN) ) & \
	( cd fe && ./node_modules/.bin/vite dev --host ) & \
	wait

run: ## Jalankan API lokal saja (tanpa menyentuh Docker)
	@set -a; [ -f $(ENV_FILE) ] && . ./$(ENV_FILE); set +a; go -C be run $(API_MAIN)

worker: ## Jalankan worker lokal saja
	@set -a; [ -f $(ENV_FILE) ] && . ./$(ENV_FILE); set +a; go -C be run $(WORKER_MAIN)

seed: ## Isi data user awal (admin, asesi, asesor, validator)
	@set -a; [ -f $(ENV_FILE) ] && . ./$(ENV_FILE); set +a; go -C be run $(SEED_MAIN)

migrate: ## Terapkan ulang seluruh file migrasi ke Postgres yang berjalan
	@for f in be/migrations/*.sql; do \
		echo "==> Menerapkan $$f"; \
		$(COMPOSE) exec -T postgres psql -v ON_ERROR_STOP=1 -U postgres -d lsp < "$$f"; \
	done

psql: ## Masuk ke shell psql
	$(COMPOSE) exec postgres psql -U postgres -d lsp

redis-cli: ## Masuk ke shell redis-cli pada Dragonfly
	$(COMPOSE) exec dragonfly redis-cli

build: ## Kompilasi seluruh paket backend
	go -C be build ./...

test: ## Jalankan seluruh test backend
	go -C be test ./...

fmt: ## Format kode Go backend
	gofmt -w be/cmd be/internal

vet: ## Jalankan go vet backend
	go -C be vet ./...

tidy: ## Rapikan dependency backend
	go -C be tidy

fe-install: ## Install dependency frontend (bun)
	cd fe && bun install

fe-dev: ## Jalankan dev server frontend (bisa diakses dari device lain di jaringan)
	cd fe && ./node_modules/.bin/vite dev --host

fe-build: ## Build frontend
	cd fe && ./node_modules/.bin/vite build

fe-typecheck: ## Typecheck frontend
	cd fe && ./node_modules/.bin/tsc --noEmit

r2-cors: ## Terapkan CORS policy bucket R2 dari be/cors.json (via bunx wrangler)
	@set -a; [ -f $(ENV_FILE) ] && . ./$(ENV_FILE); set +a; \
	test -n "$$R2_BUCKET" || { echo "R2_BUCKET belum diisi di $(ENV_FILE)"; exit 1; }; \
	bunx wrangler r2 bucket cors set "$$R2_BUCKET" --file be/cors.json --force

clean: ## Hentikan container dan hapus volume data
	$(COMPOSE) down -v
