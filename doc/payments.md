# Paiements Stripe & billets

Encaissement carte via Stripe PaymentIntent, création d’un billet après confirmation, dashboard admin `/lekeza`.

## Flux public

1. Sélection formule / jour / quantité sur la landing.
2. Modale : récap → **Continuer** → coordonnées (nom, prénom, email, ville) → formulaire carte (Elements Stripe).
3. `POST /api/payments/intent` crée un PaymentIntent (montant **recalculé serveur**).
4. Le front confirme la carte (`confirmCardPayment`), puis `POST /api/payments/confirm`.
5. Un billet `status=valid` est inséré (idempotent sur `stripe_payment_intent_id`).
6. Backup : webhook `POST /api/webhooks/stripe` (`payment_intent.succeeded`).

Jours : `formula=pass` → `day_vendredi` et `day_samedi` à true ; `formula=day` → seul le jour choisi.

## Endpoints

| Méthode | Path | Auth | Notes |
|---------|------|------|-------|
| POST | `/api/payments/intent` | public + CSRF + rate-limit | body acheteur + formule |
| POST | `/api/payments/confirm` | public + CSRF + rate-limit | `payment_intent_id` |
| POST | `/api/webhooks/stripe` | signature Stripe | CSRF exempt |
| GET | `/api/admin/billets` | JWT admin | liste + stats |
| PATCH | `/api/admin/billets/{id}` | JWT admin | `status` / `comments` |

## Variables d’environnement

**Backend** (`backend/.env`) :

| Variable | Rôle |
|----------|------|
| `STRIPE_SECRET_KEY` | Clé secrète Stripe (obligatoire en prod) |
| `STRIPE_WEBHOOK_SECRET` | Secret endpoint webhook |
| `TICKET_PRICE_DAY_CENTS` | Défaut `2000` (20 €) |
| `TICKET_PRICE_PASS_CENTS` | Défaut `3000` (30 €) |
| `RATE_LIMIT_PAYMENTS` | Défaut `20` / min / IP |

**Frontend** :

| Variable | Rôle |
|----------|------|
| `VITE_STRIPE_PUBLISHABLE_KEY` | Clé publique Stripe |

## Schéma `billets`

Voir [`database/init/postgres/01_schema.sql`](../database/init/postgres/01_schema.sql).  
Volume déjà initialisé : appliquer [`database/migrations/001_billets.sql`](../database/migrations/001_billets.sql).

## CSRF anonyme

Les GET API posent le cookie `XSRF-TOKEN` s’il est absent, pour permettre le checkout sans session.

## Admin

Après connexion `/lekeza/login` : stats (total, CA, valid/used, places ven/sam) + table + action « Marquer utilisé ».
