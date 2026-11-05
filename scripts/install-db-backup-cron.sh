#!/usr/bin/env sh
set -eu

ROOT="$(CDPATH= cd -- "$(dirname "$0")/.." && pwd)"
DB_ENV_FILE="$ROOT/database/.env"
BACKUP_DIR="${BACKUP_DIR:-/var/backups/app}"
BACKUP_LOG="${BACKUP_LOG:-/var/log/app-db-backup.log}"
MARKER="# app-db-backup"

usage() {
  cat <<EOF
Usage: $0 [--uninstall]
EOF
}

ensure_backup_dirs() {
  mkdir -p "$BACKUP_DIR"
  touch "$BACKUP_DIR/.write_test"
  rm -f "$BACKUP_DIR/.write_test"

  _log_dir=$(dirname "$BACKUP_LOG")
  mkdir -p "$_log_dir"
  touch "$BACKUP_LOG"
}

read_tz_from_env() {
  [ -f "$DB_ENV_FILE" ] || {
    echo "ERREUR: $DB_ENV_FILE introuvable" >&2
    exit 1
  }

  _tz=$(awk -F= '/^TZ=/{print $2}' "$DB_ENV_FILE" | tail -n 1 | tr -d '"' | tr -d "'")
  if [ -z "$_tz" ]; then
    _tz=UTC
  fi
  printf '%s' "$_tz"
}

UNINSTALL=false
while [ $# -gt 0 ]; do
  case "$1" in
    --uninstall)
      UNINSTALL=true
      shift
      ;;
    -h|--help)
      usage
      exit 0
      ;;
    *)
      echo "Argument inconnu : $1" >&2
      usage
      exit 1
      ;;
  esac
done

if [ "$UNINSTALL" = "true" ]; then
  _current=$(crontab -l 2>/dev/null || true)
  [ -n "$_current" ] || {
    echo "Crontab vide — rien à retirer."
    exit 0
  }
  printf '%s\n' "$_current" | grep -v "$MARKER" | grep -v 'scripts/db-backup.sh' | crontab -
  echo "Entrées app-db-backup retirées du crontab."
  exit 0
fi

TZ_VALUE=$(read_tz_from_env)
ensure_backup_dirs
BACKUP_SH="$ROOT/scripts/db-backup.sh"

if [ ! -x "$BACKUP_SH" ]; then
  chmod +x "$BACKUP_SH" "$ROOT/scripts/db-restore.sh" "$ROOT/scripts/lib/db-backup-common.sh" 2>/dev/null || true
fi

_new_block=$(cat <<EOF
$MARKER
TZ=$TZ_VALUE
0 * * * * cd $ROOT && $BACKUP_SH hourly >> $BACKUP_LOG 2>&1 $MARKER
5 2 * * * cd $ROOT && $BACKUP_SH daily >> $BACKUP_LOG 2>&1 $MARKER
10 3 1 * * cd $ROOT && $BACKUP_SH monthly >> $BACKUP_LOG 2>&1 $MARKER
15 4 1 1 * cd $ROOT && $BACKUP_SH yearly >> $BACKUP_LOG 2>&1 $MARKER
EOF
)

_current=$(crontab -l 2>/dev/null || true)
_filtered=$(printf '%s\n' "$_current" | grep -v "$MARKER" | grep -v 'scripts/db-backup.sh' || true)

{
  if [ -n "$_filtered" ]; then
    printf '%s\n' "$_filtered"
  fi
  printf '%s\n' "$_new_block"
} | crontab -

echo "Crontab installé (TZ=$TZ_VALUE)."
echo "Logs : $BACKUP_LOG"
