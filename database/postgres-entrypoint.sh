#!/bin/sh
# Enveloppe l'entrypoint officiel Postgres : init habituelle + sync admin bootstrap.
set -eu

sync_admin_password() {
  i=0
  while [ "$i" -lt 60 ]; do
    if pg_isready -U "${POSTGRES_USER}" -d "${POSTGRES_DB}" >/dev/null 2>&1; then
      /bin/sh /usr/local/bin/comores-sync-admin-password.sh
      return 0
    fi
    i=$((i + 1))
    sleep 1
  done
  echo "comores: PostgreSQL n'est pas prêt, sync admin ignorée" >&2
  return 1
}

# Image CMD est null dès que compose override entrypoint : forcer postgres.
if [ $# -eq 0 ]; then
  set -- postgres
fi
if [ "${1#-}" != "$1" ]; then
  set -- postgres "$@"
fi

if [ "$1" = "postgres" ]; then
  sync_admin_password &
  exec docker-entrypoint.sh "$@"
fi

exec docker-entrypoint.sh "$@"
