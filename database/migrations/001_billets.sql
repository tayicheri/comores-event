-- Migration manuelle (volume Postgres déjà initialisé).
-- Appliquer :
--   docker compose exec -T db-postgres \
--     psql -U "$POSTGRES_USER" -d "$POSTGRES_DB" < database/migrations/001_billets.sql

DO $$ BEGIN
    CREATE TYPE billet_status AS ENUM ('valid', 'used');
EXCEPTION
    WHEN duplicate_object THEN NULL;
END $$;

CREATE TABLE IF NOT EXISTS billets (
    id                         UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    purchased_at               TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    first_name                 TEXT NOT NULL,
    last_name                  TEXT NOT NULL,
    email                      TEXT NOT NULL,
    city                       TEXT NOT NULL,
    quantity                   INTEGER NOT NULL CHECK (quantity >= 1),
    day_vendredi               BOOLEAN NOT NULL DEFAULT FALSE,
    day_samedi                 BOOLEAN NOT NULL DEFAULT FALSE,
    comments                   TEXT,
    status                     billet_status NOT NULL DEFAULT 'valid',
    amount_cents               INTEGER NOT NULL CHECK (amount_cents >= 0),
    stripe_payment_intent_id   TEXT NOT NULL UNIQUE,
    created_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at                 TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    CONSTRAINT billets_at_least_one_day CHECK (day_vendredi OR day_samedi)
);

CREATE INDEX IF NOT EXISTS idx_billets_purchased_at ON billets (purchased_at DESC);
CREATE INDEX IF NOT EXISTS idx_billets_status ON billets (status);
CREATE INDEX IF NOT EXISTS idx_billets_email ON billets (email);
