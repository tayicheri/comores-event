#!/bin/sh
# Retire email et mot de passe en clair des GUC une fois le compte écrit dans users.
set -eu

case "${POSTGRES_DB:-}" in
  ''|*[!A-Za-z0-9_]*)
    echo "POSTGRES_DB invalide pour ALTER DATABASE" >&2
    exit 1
    ;;
esac

psql -v ON_ERROR_STOP=1 --username "$POSTGRES_USER" --dbname "$POSTGRES_DB" <<EOSQL
ALTER DATABASE ${POSTGRES_DB} RESET comores.admin_password;
ALTER DATABASE ${POSTGRES_DB} RESET comores.admin_email;
EOSQL
