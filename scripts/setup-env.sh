#!/usr/bin/env sh
set -e
ROOT="$(cd "$(dirname "$0")/.." && pwd)"
cd "$ROOT"

copy_if_missing() {
  src="$1"
  dest="$2"
  if [ ! -f "$dest" ]; then
    cp "$src" "$dest"
    echo "Créé : $dest"
  else
    echo "Conservé (existe déjà) : $dest"
  fi
}

copy_if_missing "database/.env.example" "database/.env"
copy_if_missing "backend/.env.example" "backend/.env"
copy_if_missing ".env.example" ".env"
copy_if_missing "frontend/.env.development.example" "frontend/.env.development"
copy_if_missing "frontend/.env.production.example" "frontend/.env.production"

ensure_env_key() {
  file="$1"
  key="$2"
  default="$3"
  if [ -f "$file" ] && ! grep -q "^${key}=" "$file"; then
    printf '\n%s=%s\n' "$key" "$default" >>"$file"
    echo "Ajouté : $key dans $file"
  fi
}

ensure_env_key "database/.env" "ADMIN_EMAIL" "admin@example.com"
ensure_env_key ".env" "GRAFANA_ALERT_EMAIL" "ops@example.com"

echo ""
echo "Éditez database/.env, backend/.env et .env avant de lancer Docker."
echo "En dev : alignez DB_PASS / POSTGRES_PASSWORD entre database/.env et backend/.env."
echo "Changez JWT_SECRET (min 32 caractères) et ADMIN_PASSWORD."
echo "En prod / observabilité : renseignez CADDY_* et GRAFANA_* dans .env."
