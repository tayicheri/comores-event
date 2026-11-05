# shellcheck shell=sh
# Helpers pour lire/écrire le .env racine sans casser les valeurs contenant des $.

ROOT_ENV_FILE="${ROOT_ENV_FILE:-.env}"

load_root_env() {
  if [ ! -f "$ROOT_ENV_FILE" ]; then
    echo "Fichier $ROOT_ENV_FILE manquant." >&2
    return 1
  fi

  set -a
  while IFS= read -r line || [ -n "$line" ]; do
    line=$(printf '%s' "$line" | sed 's/^[[:space:]]*//')
    case "$line" in
      ''|\#*) continue ;;
      *=*) ;;
      *) continue ;;
    esac

    key=${line%%=*}
    val=${line#*=}
    key=$(printf '%s' "$key" | sed 's/[[:space:]]*$//')
    val=$(printf '%s' "$val" | sed 's/^[[:space:]]*//')

    case "$val" in
      \"*\") val=${val#\"}; val=${val%\"} ;;
      \'*\') val=${val#\'}; val=${val%\'} ;;
    esac

    export "$key"="$val"
  done < "$ROOT_ENV_FILE"
  set +a
}

set_env_var() {
  key="$1"
  value="$2"
  file="$ROOT_ENV_FILE"
  tmp="${file}.tmp.$$"

  if [ ! -f "$file" ]; then
    printf "%s='%s'\n" "$key" "$value" >"$file"
    return 0
  fi

  if grep -q "^${key}=" "$file"; then
    awk -v k="$key" -v v="$value" '
      BEGIN { replaced = 0 }
      $0 ~ "^" k "=" {
        print k "=\047" v "\047"
        replaced = 1
        next
      }
      { print }
      END {
        if (!replaced) print k "=\047" v "\047"
      }
    ' "$file" >"$tmp"
  else
    cp "$file" "$tmp"
    printf "%s='%s'\n" "$key" "$value" >>"$tmp"
  fi

  mv "$tmp" "$file"
}
