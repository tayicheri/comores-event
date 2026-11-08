# comores-event — monorepo template

Scaffold réutilisable calqué sur le pattern FoodyFly : **Go API** + **Vue 3** + Postgres + Redis + Docker + observabilité.

`mobile/` et `workers/` sont réservés (vides).

## Démarrage rapide (dev)

```bash
./scripts/setup-env.sh
# Éditer JWT_SECRET, mots de passe DB si besoin
./scripts/docker-dev.sh
```

| Service | URL |
|---------|-----|
| Frontend | http://localhost:3000 |
| API | http://localhost:8080 |
| Health | http://localhost:8080/health/ready |
| pgAdmin | http://localhost:5050 |
| Redis | localhost:6379 |

Compte admin seed : `admin@example.com` / `change_me_admin_password`

## Structure

```
backend/          API Go (auth + middlewares)
frontend/         Vue 3 + Vite + Tailwind
database/         Postgres init SQL
observability/    Prometheus, Grafana, Loki, …
scripts/          setup-env, docker-*, backup DB
mobile/           (vide)
workers/          (vide)
```

## Scripts

| Script | Rôle |
|--------|------|
| `./scripts/setup-env.sh` | Copie les `.env.example` → `.env` |
| `./scripts/docker-dev.sh` | Stack dev |
| `./scripts/docker-prod.sh` | Stack prod + Caddy |
| `./scripts/docker-observability.sh` | Monitoring (nécessite le réseau prod) |
| `./scripts/compose-validate.sh` | Valide les compose |
| `./scripts/db-backup.sh` / `db-restore.sh` | Sauvegardes Postgres |

## Réutiliser comme template

1. Copier le repo.
2. Renommer le module Go dans `backend/go.mod` (`comores-event-api` → `votre-api`).
3. Adapter `VITE_APP_NAME`, domaines `CADDY_*`, préfixes conteneurs `app-*` si besoin.
4. Changer tous les secrets (`JWT_SECRET`, mots de passe DB/admin/Grafana).

## Prod & monitoring

```bash
./scripts/docker-prod.sh
./scripts/docker-observability.sh
```

Grafana local : http://127.0.0.1:3001  
Grafana public : domaine `CADDY_MONITORING_DOMAIN` (basic auth Caddy).

## Auth & billetterie

- Auth : [`doc/auth.md`](doc/auth.md)
- Paiements Stripe / billets / admin : [`doc/payments.md`](doc/payments.md)
- Frontend SSG / SEO : [`doc/frontend-ssg.md`](doc/frontend-ssg.md)

Middlewares : Recovery, Logger, Metrics, CORS, CSRF, SecurityHeaders, rate-limit Redis.
