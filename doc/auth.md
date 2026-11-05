# Auth

Authentification email/password avec cookies. Côté front : login **admin uniquement** sur `/lekeza/login`.

## Routes front

| Path | Accès | Notes |
|------|--------|-------|
| `/` | public | One-page événement (maquette) |
| `/lekeza` | rôle `admin` | Dashboard billetterie (liste + stats) |
| `/lekeza/login` | public | Formulaire de connexion admin |
| `/login` | redirect | → `/lekeza/login` |

Pas d’inscription UI. Un utilisateur authentifié sans rôle `admin` est renvoyé vers `/`.

## Cookies

| Cookie | HttpOnly | Rôle |
|--------|----------|------|
| `auth_token` | oui | JWT access (15 min) |
| `refresh_token` | oui | Opaque refresh (7 j), path `/api/auth` |
| `XSRF-TOKEN` | non | Double-submit CSRF (`X-XSRF-TOKEN`) |

## Endpoints

| Méthode | Path | Auth | Notes |
|---------|------|------|-------|
| POST | `/api/auth/register` | public | CSRF exempt, rate-limit (API dormant côté UI) |
| POST | `/api/auth/login` | public | CSRF exempt, rate-limit |
| POST | `/api/auth/refresh` | refresh cookie | CSRF exempt |
| POST | `/api/auth/logout` | JWT + CSRF | |
| POST | `/api/auth/logout-all` | JWT + CSRF | bump `token_version` |
| GET | `/api/auth/me` | JWT | |
| PATCH | `/api/auth/me` | JWT + CSRF | `display_name` |
| GET | `/health` | public | liveness |
| GET | `/health/ready` | public | DB + Redis |
| GET | `/metrics` | public (bloqué via Caddy en prod) | Prometheus |

## Corps JSON

Login :

```json
{ "email": "admin@example.com", "password": "********" }
```

## Compte admin bootstrap (`ADMIN_*`)

Variables dans `database/.env` (créées par `./scripts/setup-env.sh`) :

| Variable | Rôle |
|----------|------|
| `ADMIN_EMAIL` | Email du compte admin |
| `ADMIN_PASSWORD` | Mot de passe en clair **uniquement** dans `.env` (≥ 8 car.) |

Flux (pattern foodyfly) :

1. **1er init volume** : `00_admin_password.sh` pose les GUC Postgres → `01_schema.sql` crée l’utilisateur id `20000000-0000-4000-8000-000000000001` avec `crypt(..., gen_salt('bf', 14))` + rôle `admin` → `02_reset_admin_guc.sh` retire les GUC.
2. **Chaque démarrage** : `sync-admin-password.sh` (via `postgres-entrypoint.sh`) réaligne email / hash bcrypt sur `ADMIN_*` si le volume existe déjà.

Le mot de passe n’est **jamais** stocké en clair en base — uniquement `password_hash` bcrypt (cost 14), compatible avec le backend Go.
