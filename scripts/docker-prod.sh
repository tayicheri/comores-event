#!/usr/bin/env sh
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

if [ ! -f database/.env ] || [ ! -f backend/.env ]; then
  echo "Fichiers .env manquants. Exécutez : ./scripts/setup-env.sh"
  exit 1
fi

if [ ! -f .env ]; then
  echo "Fichier .env racine manquant. Exécutez : ./scripts/setup-env.sh"
  exit 1
fi

exec docker compose -f docker-compose.prod.yml --env-file .env up -d --build "$@"
