#!/usr/bin/env sh
set -euo pipefail

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
DB_BACKUP_ROOT="$ROOT"
. "$ROOT/scripts/lib/db-backup-common.sh"

usage() {
  cat <<EOF
Usage: $0 [--file CHEMIN] [--yes] [--dry-run]
EOF
}

BACKUP_FILE=""
ASSUME_YES=false
DRY_RUN=false

while [ $# -gt 0 ]; do
  case "$1" in
    --file)
      shift
      BACKUP_FILE="${1:?--file requiert un chemin}"
      shift
      ;;
    --yes)
      ASSUME_YES=true
      shift
      ;;
    --dry-run)
      DRY_RUN=true
      export DRY_RUN
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      log_error "Argument inconnu : $1"
      usage
      exit 1
      ;;
  esac
done

load_database_env

if [ -z "$BACKUP_FILE" ]; then
  BACKUP_FILE=$(find_latest_backup) || {
    log_error "Aucune archive .sql.gz dans $BACKUP_DIR"
    exit 1
  }
fi

[ -f "$BACKUP_FILE" ] || {
  log_error "Fichier introuvable : $BACKUP_FILE"
  exit 1
}

log "Archive sélectionnée : $BACKUP_FILE"

if [ "$ASSUME_YES" != "true" ] && [ "$DRY_RUN" != "true" ]; then
  printf 'Confirmer la restauration en tapant RESTORE : '
  read -r _confirm
  if [ "$_confirm" != "RESTORE" ]; then
    log "Restauration annulée"
    exit 1
  fi
fi

stop_dependent_services
if run_psql_restore "$BACKUP_FILE"; then
  log "Restauration réussie depuis $BACKUP_FILE"
else
  log_error "Échec de la restauration"
  start_dependent_services
  exit 1
fi
start_dependent_services
