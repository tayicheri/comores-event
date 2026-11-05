# Auth

Authentification email/password avec cookies.

## Cookies

| Cookie | HttpOnly | Rôle |
|--------|----------|------|
| `auth_token` | oui | JWT access (15 min) |
| `refresh_token` | oui | Opaque refresh (7 j), path `/api/auth` |
| `XSRF-TOKEN` | non | Double-submit CSRF (`X-XSRF-TOKEN`) |

## Endpoints

| Méthode | Path | Auth | Notes |
|---------|------|------|-------|
| POST | `/api/auth/register` | public | CSRF exempt, rate-limit |
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

Register / Login :

```json
{ "email": "user@example.com", "password": "********", "display_name": "Alice" }
```

`display_name` optionnel à l’inscription.

## Compte seed

Après init Postgres : `admin@example.com` / `change_me_admin_password` (rôle `admin`).
