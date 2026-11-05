#!/usr/bin/env sh
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

. "$ROOT/scripts/lib/env-utils.sh"

REGENERATE_AUTH_HASH=0
COMPOSE_ARGS=""

for arg in "$@"; do
  case "$arg" in
    --regenerate-auth-hash)
      REGENERATE_AUTH_HASH=1
      ;;
    *)
      COMPOSE_ARGS="$COMPOSE_ARGS $arg"
      ;;
  esac
done

if [ ! -f .env ]; then
  echo "Fichier .env racine manquant. Exécutez : ./scripts/setup-env.sh"
  exit 1
fi

load_root_env

if ! docker network inspect app_caddy_net >/dev/null 2>&1; then
  echo "Réseau app_caddy_net introuvable."
  echo "Démarrez d'abord la stack prod : ./scripts/docker-prod.sh"
  exit 1
fi

if [ -z "${CADDY_MONITORING_DOMAIN:-}" ]; then
  echo "Définissez CADDY_MONITORING_DOMAIN dans .env."
  exit 1
fi

if [ -z "${GRAFANA_ADMIN_PASSWORD:-}" ]; then
  echo "Définissez GRAFANA_ADMIN_PASSWORD dans .env."
  exit 1
fi

GRAFANA_BASIC_AUTH_USER="${GRAFANA_BASIC_AUTH_USER:-admin}"
set_env_var "GRAFANA_BASIC_AUTH_USER" "$GRAFANA_BASIC_AUTH_USER"

EXPECTED_ROOT_URL="https://${CADDY_MONITORING_DOMAIN}"
if [ -z "${GRAFANA_ROOT_URL:-}" ]; then
  GRAFANA_ROOT_URL="$EXPECTED_ROOT_URL"
  set_env_var "GRAFANA_ROOT_URL" "$GRAFANA_ROOT_URL"
  echo "→ GRAFANA_ROOT_URL défini sur $GRAFANA_ROOT_URL"
fi

generate_basic_auth_hash() {
  docker run --rm caddy:2-alpine caddy hash-password --plaintext "$GRAFANA_ADMIN_PASSWORD"
}

sync_grafana_admin_password() {
  if ! docker ps --format '{{.Names}}' | grep -qx 'app-grafana'; then
    echo "Conteneur app-grafana introuvable — mot de passe Grafana non synchronisé." >&2
    return 1
  fi

  echo "→ Réinitialisation du mot de passe admin Grafana…"
  docker exec app-grafana grafana cli admin reset-admin-password "$GRAFANA_ADMIN_PASSWORD"
  echo "→ Mot de passe Grafana synchronisé avec GRAFANA_ADMIN_PASSWORD."
}

if [ "$REGENERATE_AUTH_HASH" -eq 1 ] || [ -z "${GRAFANA_BASIC_AUTH_HASH:-}" ]; then
  echo "→ Génération du hash basic auth Caddy…"
  GRAFANA_BASIC_AUTH_HASH="$(generate_basic_auth_hash)"
  if [ -z "$GRAFANA_BASIC_AUTH_HASH" ]; then
    echo "Échec de la génération du hash basic auth." >&2
    exit 1
  fi
  set_env_var "GRAFANA_BASIC_AUTH_HASH" "$GRAFANA_BASIC_AUTH_HASH"
  echo "→ GRAFANA_BASIC_AUTH_HASH enregistré dans .env"
fi

load_root_env

if [ -z "${GRAFANA_BASIC_AUTH_HASH:-}" ]; then
  echo "GRAFANA_BASIC_AUTH_HASH manquant dans .env après génération." >&2
  exit 1
fi

# shellcheck disable=SC2086
docker compose -f docker-compose.observability.yml --env-file .env up -d --build $COMPOSE_ARGS

echo ""
echo "Stack observabilité démarrée."
echo "Grafana public : https://${CADDY_MONITORING_DOMAIN}"
echo "Grafana local  : http://127.0.0.1:3001"
echo ""

if [ "$REGENERATE_AUTH_HASH" -eq 1 ]; then
  sync_grafana_admin_password || true
  echo ""
fi

echo "→ Recréation de Caddy pour recharger le .env…"
docker compose -f docker-compose.prod.yml --env-file .env up -d --force-recreate caddy
echo "→ Caddy recréé."
