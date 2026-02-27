#!/bin/sh

# Exit on error
set -e

echo "Waiting for PostgreSQL to be ready..."
# Wait for postgres to be healthy (checking localhost first, then the host from env)
POSTGRES_HOST=$(echo $DATABASE_URL | sed -E 's/.*@([^:]+).*/\1/')
for i in {1..30}; do
  if pg_isready -h "$POSTGRES_HOST" 2>/dev/null; then
    break
  fi
  echo "PostgreSQL not ready, waiting... ($i/30)"
  sleep 1
done

echo "Running database migrations..."
./goose -dir ./internal/database/migrations postgres "$DATABASE_URL" up

echo "Migrations completed. Starting application..."
exec ./main
