#!/usr/bin/env sh
set -euo pipefail

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
DB_BACKUP_ROOT="$ROOT"
. "$ROOT/scripts/lib/db-backup-common.sh"

usage() {
  cat <<EOF
Usage: $0 [--dry-run] hourly|daily|monthly|yearly
EOF
}

BACKUP_TYPE=""
DRY_RUN=false

while [ $# -gt 0 ]; do
  case "$1" in
    --dry-run)
      DRY_RUN=true
      export DRY_RUN
      shift
      ;;
    hourly|daily|monthly|yearly)
      BACKUP_TYPE=$1
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

[ -n "$BACKUP_TYPE" ] || {
  usage
  exit 1
}

load_database_env
ensure_backup_dir

log "Démarrage sauvegarde $BACKUP_TYPE (TZ=$TZ, conteneur=$CONTAINER)"

case "$BACKUP_TYPE" in
  hourly)
    write_backup_file "$(hourly_filename)"
    ;;
  daily)
    _today=$(backup_date)
    write_backup_file "$(daily_filename "$_today")"
    purge_hourly_for_date "$(backup_yesterday_date)"
    ;;
  monthly)
    _month=$(backup_month_for_monthly)
    write_backup_file "$(monthly_filename "$_month")"
    purge_daily_for_month "$_month"
    ;;
  yearly)
    _year=$(backup_year_for_annual)
    write_backup_file "$(yearly_filename "$_year")"
    purge_all_monthly_backups
    ;;
esac

log "Sauvegarde $BACKUP_TYPE terminée"
