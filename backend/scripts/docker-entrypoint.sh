#!/bin/sh
set -e

wait_for_postgres() {
  if [ "${SKIP_DB_WAIT:-false}" = "true" ]; then
    return 0
  fi
  host="${DB_HOST_POSTGRES:-db-postgres}"
  port="${DB_PORT:-5432}"
  user="${DB_USER:-postgres}"

  echo "Attente PostgreSQL (${host}:${port})..."
  until pg_isready -h "$host" -p "$port" -U "$user" -q; do
    sleep 1
  done
  echo "PostgreSQL prêt"
}

wait_for_postgres

# Process principal : API (air en dev, /app/main en prod)
exec "$@"
