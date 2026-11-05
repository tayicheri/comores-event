#!/bin/sh
# Génère servers.json + pgpass pour pré-enregistrer PostgreSQL dans pgAdmin (dev).
set -eu

PGADMIN_SERVERS_JSON="${PGADMIN_SERVERS_JSON:-/tmp/app-servers.json}"
PGADMIN_PGPASS_FILE="${PGADMIN_PGPASS_FILE:-/tmp/app-pgpass}"
PGADMIN_DB_HOST="${PGADMIN_DB_HOST:-db-postgres}"
PGADMIN_DB_PORT="${PGADMIN_DB_PORT:-5432}"
PGADMIN_SERVER_NAME="${PGADMIN_SERVER_NAME:-App (dev)}"

require_env() {
  eval "val=\${$1:-}"
  if [ -z "$val" ]; then
    echo "pgadmin: variable $1 manquante" >&2
    exit 1
  fi
}

is_safe_ident() {
  case "$1" in
    ''|*[!A-Za-z0-9._-]*) return 1 ;;
    *) return 0 ;;
  esac
}

pgpass_escape() {
  in=$1
  out=
  while [ -n "$in" ]; do
    rest=${in#?}
    c=${in%"$rest"}
    in=$rest
    case "$c" in
      \\) out=$out\\\\ ;;
      :) out=$out\\: ;;
      *) out=$out$c ;;
    esac
  done
  printf '%s' "$out"
}

require_env POSTGRES_USER
require_env POSTGRES_PASSWORD
require_env POSTGRES_DB

if ! is_safe_ident "$POSTGRES_USER"; then
  echo "pgadmin: POSTGRES_USER contient des caractères non supportés" >&2
  exit 1
fi
if ! is_safe_ident "$POSTGRES_DB"; then
  echo "pgadmin: POSTGRES_DB contient des caractères non supportés" >&2
  exit 1
fi
if ! is_safe_ident "$PGADMIN_DB_HOST"; then
  echo "pgadmin: PGADMIN_DB_HOST contient des caractères non supportés" >&2
  exit 1
fi
case "$PGADMIN_DB_PORT" in
  ''|*[!0-9]*)
    echo "pgadmin: PGADMIN_DB_PORT invalide" >&2
    exit 1
    ;;
esac

pw_lines=$(printf '%s\n' "$POSTGRES_PASSWORD" | wc -l)
if [ "$pw_lines" -ne 1 ]; then
  echo "pgadmin: POSTGRES_PASSWORD ne doit pas contenir de retour à la ligne" >&2
  exit 1
fi

user_esc=$(pgpass_escape "$POSTGRES_USER")
pass_esc=$(pgpass_escape "$POSTGRES_PASSWORD")

umask 077
cat > "$PGADMIN_SERVERS_JSON" <<EOF
{
  "Servers": {
    "1": {
      "Name": "${PGADMIN_SERVER_NAME}",
      "Group": "Servers",
      "Host": "${PGADMIN_DB_HOST}",
      "Port": ${PGADMIN_DB_PORT},
      "MaintenanceDB": "${POSTGRES_DB}",
      "Username": "${POSTGRES_USER}",
      "SSLMode": "prefer",
      "PassFile": "${PGADMIN_PGPASS_FILE}"
    }
  }
}
EOF

printf '%s\n' "${PGADMIN_DB_HOST}:${PGADMIN_DB_PORT}:*:${user_esc}:${pass_esc}" > "$PGADMIN_PGPASS_FILE"
chmod 600 "$PGADMIN_PGPASS_FILE"
