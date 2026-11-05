#!/bin/sh
# Prépare servers.json / pgpass puis délègue à l'entrypoint officiel pgAdmin.
set -eu

DIR=$(CDPATH= cd -- "$(dirname "$0")" && pwd)

sh "$DIR/generate-config.sh"

export PGADMIN_SERVER_JSON_FILE="${PGADMIN_SERVERS_JSON:-/tmp/app-servers.json}"
export PGPASS_FILE="${PGADMIN_PGPASS_FILE:-/tmp/app-pgpass}"

exec /entrypoint.sh "$@"
