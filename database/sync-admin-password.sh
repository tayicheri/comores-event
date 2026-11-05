#!/bin/sh
# Aligne email et hash bcrypt de l'admin bootstrap (id fixe) sur ADMIN_EMAIL / ADMIN_PASSWORD.
# No-op si aucune variable n'est définie.
set -eu

ADMIN_ID='20000000-0000-4000-8000-000000000001'

if [ -z "${ADMIN_PASSWORD:-}" ] && [ -z "${ADMIN_EMAIL:-}" ]; then
  echo "comores: ADMIN_* non définis, compte admin inchangé" >&2
  exit 0
fi

if [ -n "${ADMIN_PASSWORD:-}" ]; then
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
fi

if [ -n "${ADMIN_EMAIL:-}" ]; then
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
fi

users_rel=$(psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" -tAc "SELECT to_regclass('public.users')")
if [ "$users_rel" != "users" ]; then
  echo "comores: table users absente, sync admin reportée" >&2
  exit 0
fi

if [ -n "${ADMIN_PASSWORD:-}" ]; then
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -v pw="$ADMIN_PASSWORD" \
    -v sid="$ADMIN_ID" <<'EOSQL'
UPDATE users
SET
    password_hash = crypt(:'pw', gen_salt('bf', 14)),
    token_version = token_version + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = :'sid'::uuid
  AND (
    password_hash IS NULL
    OR password_hash IS DISTINCT FROM crypt(:'pw', password_hash)
  );
EOSQL
fi

if [ -n "${ADMIN_EMAIL:-}" ]; then
  psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" \
    -v email="$ADMIN_EMAIL" \
    -v sid="$ADMIN_ID" <<'EOSQL'
UPDATE users
SET
    email = :'email',
    token_version = token_version + 1,
    updated_at = CURRENT_TIMESTAMP
WHERE id = :'sid'::uuid
  AND email IS DISTINCT FROM :'email';
EOSQL
fi
