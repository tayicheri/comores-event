#!/usr/bin/env sh
# Bibliothèque partagée backup / restore PostgreSQL App.

set -eu

BACKUP_DIR="${BACKUP_DIR:-/var/backups/app}"
CONTAINER="${CONTAINER:-app-db-postgres}"
DRY_RUN="${DRY_RUN:-false}"

log() {
  _ts=$(TZ="${TZ:-UTC}" date '+%Y-%m-%d %H:%M:%S %Z' 2>/dev/null || date '+%Y-%m-%d %H:%M:%S')
  printf '[%s] %s\n' "$_ts" "$1"
}

log_warn() {
  log "AVERTISSEMENT: $1"
}

log_error() {
  log "ERREUR: $1" >&2
}

if [ -z "${DB_BACKUP_ROOT:-}" ]; then
  log_error "DB_BACKUP_ROOT doit être défini avant de sourcer db-backup-common.sh"
  exit 1
fi

DB_ENV_FILE="${DB_ENV_FILE:-$DB_BACKUP_ROOT/database/.env}"

load_database_env() {
  if [ ! -f "$DB_ENV_FILE" ]; then
    log_error "Fichier introuvable : $DB_ENV_FILE"
    exit 1
  fi

  set -a
  . "$DB_ENV_FILE"
  set +a

  if [ -z "${POSTGRES_USER:-}" ] || [ -z "${POSTGRES_DB:-}" ]; then
    log_error "POSTGRES_USER et POSTGRES_DB sont requis dans $DB_ENV_FILE"
    exit 1
  fi

  if [ -z "${TZ:-}" ]; then
    TZ=UTC
    log_warn "TZ absent dans database/.env — utilisation de UTC"
  fi
  export TZ

  POSTGRES_PASSWORD="${POSTGRES_PASSWORD:-}"
}

backup_date() {
  TZ="$TZ" date +%Y-%m-%d
}

backup_hour() {
  TZ="$TZ" date +%H
}

backup_month_for_monthly() {
  if date -d 'last month' +%Y-%m >/dev/null 2>&1; then
    TZ="$TZ" date -d 'last month' +%Y-%m
  else
    TZ="$TZ" date -v-1m +%Y-%m
  fi
}

backup_year_for_annual() {
  _y=$(TZ="$TZ" date +%Y)
  printf '%s' $((_y - 1))
}

backup_yesterday_date() {
  if date -d 'yesterday' +%Y-%m-%d >/dev/null 2>&1; then
    TZ="$TZ" date -d 'yesterday' +%Y-%m-%d
  else
    TZ="$TZ" date -v-1d +%Y-%m-%d
  fi
}

hourly_filename() {
  _date="${1:-$(backup_date)}"
  _hour="${2:-$(backup_hour)}"
  printf '%s/app_%s_%s-hourly.sql.gz' "$BACKUP_DIR" "$_date" "$_hour"
}

daily_filename() {
  _date="${1:-$(backup_date)}"
  printf '%s/app_%s_daily.sql.gz' "$BACKUP_DIR" "$_date"
}

monthly_filename() {
  _month="${1:-$(backup_month_for_monthly)}"
  printf '%s/app_%s_monthly.sql.gz' "$BACKUP_DIR" "$_month"
}

yearly_filename() {
  _year="${1:-$(backup_year_for_annual)}"
  printf '%s/app_%s_yearly.sql.gz' "$BACKUP_DIR" "$_year"
}

ensure_backup_dir() {
  if [ "$DRY_RUN" = "true" ]; then
    log "DRY-RUN : mkdir -p $BACKUP_DIR"
    return 0
  fi

  mkdir -p "$BACKUP_DIR"
  touch "$BACKUP_DIR/.write_test"
  rm -f "$BACKUP_DIR/.write_test"
}

ensure_container_running() {
  if ! docker inspect -f '{{.State.Running}}' "$CONTAINER" 2>/dev/null | grep -q true; then
    log_error "Conteneur $CONTAINER absent ou arrêté"
    exit 1
  fi
}

run_pg_dump() {
  ensure_container_running
  docker exec -e PGPASSWORD="${POSTGRES_PASSWORD}" "$CONTAINER" \
    pg_dump --clean --if-exists -U "$POSTGRES_USER" -d "$POSTGRES_DB"
}

write_backup_file() {
  _dest="$1"
  if [ "$DRY_RUN" = "true" ]; then
    log "DRY-RUN : pg_dump | gzip > $_dest"
    return 0
  fi

  if ! run_pg_dump | gzip -c > "$_dest"; then
    log_error "Échec du dump vers $_dest"
    rm -f "$_dest"
    exit 1
  fi

  log "Sauvegarde créée : $_dest"
}

purge_hourly_for_date() {
  _date="$1"
  _pattern="$BACKUP_DIR/app_${_date}_*-hourly.sql.gz"
  _found=false

  for _f in $_pattern; do
    [ -e "$_f" ] || continue
    _found=true
    if [ "$DRY_RUN" = "true" ]; then
      log "DRY-RUN : rm -f $_f"
    else
      rm -f "$_f"
      log "Supprimé (horaire) : $_f"
    fi
  done

  if [ "$_found" = "false" ]; then
    log "Aucune sauvegarde horaire à purger pour $_date"
  fi
}

purge_daily_for_month() {
  _month="$1"
  _pattern="$BACKUP_DIR/app_${_month}-*_daily.sql.gz"
  _found=false

  for _f in $_pattern; do
    [ -e "$_f" ] || continue
    _found=true
    if [ "$DRY_RUN" = "true" ]; then
      log "DRY-RUN : rm -f $_f"
    else
      rm -f "$_f"
      log "Supprimé (journalier) : $_f"
    fi
  done

  if [ "$_found" = "false" ]; then
    log "Aucune sauvegarde journalière à purger pour $_month"
  fi
}

purge_all_monthly_backups() {
  _pattern="$BACKUP_DIR/app_*_monthly.sql.gz"
  _found=false

  for _f in $_pattern; do
    [ -e "$_f" ] || continue
    _found=true
    if [ "$DRY_RUN" = "true" ]; then
      log "DRY-RUN : rm -f $_f"
    else
      rm -f "$_f"
      log "Supprimé (mensuel) : $_f"
    fi
  done

  if [ "$_found" = "false" ]; then
    log "Aucune sauvegarde mensuelle à purger"
  fi
}

find_latest_backup() {
  _latest=""
  _latest_mtime=0

  [ -d "$BACKUP_DIR" ] || return 1

  for _f in "$BACKUP_DIR"/*.sql.gz; do
    [ -f "$_f" ] || continue
    if stat -f %m "$_f" >/dev/null 2>&1; then
      _mtime=$(stat -f %m "$_f")
    else
      _mtime=$(stat -c %Y "$_f")
    fi
    if [ "$_mtime" -gt "$_latest_mtime" ]; then
      _latest_mtime=$_mtime
      _latest=$_f
    fi
  done

  [ -n "$_latest" ] || return 1
  printf '%s\n' "$_latest"
}

stop_dependent_services() {
  for _svc in app-api app-caddy; do
    if docker inspect -f '{{.State.Running}}' "$_svc" 2>/dev/null | grep -q true; then
      if [ "$DRY_RUN" = "true" ]; then
        log "DRY-RUN : docker stop $_svc"
      else
        log "Arrêt de $_svc…"
        docker stop "$_svc" >/dev/null
      fi
    fi
  done
}

start_dependent_services() {
  for _svc in app-api app-caddy; do
    if docker inspect "$_svc" >/dev/null 2>&1; then
      if [ "$DRY_RUN" = "true" ]; then
        log "DRY-RUN : docker start $_svc"
      else
        docker start "$_svc" >/dev/null 2>&1 || true
        log "Redémarrage de $_svc"
      fi
    fi
  done
}

run_psql_restore() {
  _file="$1"
  ensure_container_running

  if [ "$DRY_RUN" = "true" ]; then
    log "DRY-RUN : gunzip -c $_file | docker exec -i $CONTAINER psql …"
    return 0
  fi

  gunzip -c "$_file" | docker exec -i -e PGPASSWORD="${POSTGRES_PASSWORD}" "$CONTAINER" \
    psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" -v ON_ERROR_STOP=1
}
