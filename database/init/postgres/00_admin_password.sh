#!/bin/sh
# Pose les GUC lus par 01_schema.sql (email + crypt du mot de passe admin bootstrap).
set -eu

if [ -z "${ADMIN_PASSWORD:-}" ]; then
  echo "ADMIN_PASSWORD est requis (database/.env)" >&2
  exit 1
fi

pw_len=$(printf '%s' "$ADMIN_PASSWORD" | wc -c)
if [ "$pw_len" -lt 8 ]; then
  echo "ADMIN_PASSWORD doit faire au moins 8 caractères" >&2
  exit 1
fi

pw_lines=$(printf '%s\n' "$ADMIN_PASSWORD" | wc -l)
if [ "$pw_lines" -ne 1 ]; then
  echo "ADMIN_PASSWORD ne doit pas contenir de retour à la ligne" >&2
  exit 1
fi

if [ -z "${ADMIN_EMAIL:-}" ]; then
  echo "ADMIN_EMAIL est requis (database/.env)" >&2
  exit 1
fi

email_lines=$(printf '%s\n' "$ADMIN_EMAIL" | wc -l)
if [ "$email_lines" -ne 1 ]; then
  echo "ADMIN_EMAIL ne doit pas contenir de retour à la ligne" >&2
  exit 1
fi

case "$ADMIN_EMAIL" in
  *[[:space:]]*)
    echo "ADMIN_EMAIL ne doit pas contenir d'espace" >&2
    exit 1
    ;;
  *@*.*)
    ;;
  *)
    echo "ADMIN_EMAIL invalide" >&2
    exit 1
    ;;
esac

case "${POSTGRES_DB:-}" in
  ''|*[!A-Za-z0-9_]*)
    echo "POSTGRES_DB invalide pour ALTER DATABASE" >&2
    exit 1
    ;;
esac

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
  -v pw="$ADMIN_PASSWORD" \
  -v email="$ADMIN_EMAIL" <<EOSQL
ALTER DATABASE ${POSTGRES_DB} SET comores.admin_password = :'pw';
ALTER DATABASE ${POSTGRES_DB} SET comores.admin_email = :'email';
EOSQL
