SHELL := /bin/sh

ENV_FILE := be/.env
COMPOSE  := docker compose
API_MAIN := ./cmd/api/main.go
WORKER_MAIN := ./cmd/worker/main.go
SEED_MAIN := ./cmd/seed/main.go

.DEFAULT_GOAL := help

.PHONY: help dev up down stop restart logs ps run worker seed migrate psql redis-cli build test fmt vet tidy fe-install clean

help: ## Tampilkan daftar target
	@echo "Target yang tersedia:"
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | \
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

dev: ## Nyalakan Postgres + Dragonfly (temporary) lalu jalankan API lokal; container dihentikan saat API keluar
	@set -a; [ -f $(ENV_FILE) ] && . ./$(ENV_FILE); set +a; \
		trap 'echo "==> Menghentikan container Postgres + Dragonfly"; $(COMPOSE) stop' EXIT INT TERM; \
		$(COMPOSE) up -d; \
		echo "==> Menjalankan API lokal: go -C be run $(API_MAIN)"; \
		go -C be run $(API_MAIN)

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
	bun install

clean: ## Hentikan container dan hapus volume data
	$(COMPOSE) down -v
