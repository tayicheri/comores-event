#!/usr/bin/env sh
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

echo "Validation docker-compose.yml"
docker compose -f docker-compose.yml config >/dev/null

if [ -f .env ]; then
  echo "Validation docker-compose.prod.yml"
  docker compose -f docker-compose.prod.yml --env-file .env config >/dev/null

  echo "Validation docker-compose.observability.yml"
  docker compose -f docker-compose.observability.yml --env-file .env config >/dev/null
else
  echo "Fichier .env absent : validation prod/observability ignorée."
fi

echo "Validation OK."
