# PostgreSQL untuk pengembangan lokal.
# Migrasi di folder ./migrations otomatis dijalankan saat volume data masih kosong.
FROM postgres:16-alpine

ENV POSTGRES_USER=postgres \
    POSTGRES_PASSWORD=postgres \
    POSTGRES_DB=lsp

COPY be/migrations/ /docker-entrypoint-initdb.d/

EXPOSE 5432
