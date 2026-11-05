-- Schéma minimal auth (users, roles, refresh_tokens)
-- 00_admin_password.sh pose les GUC comores.admin_* avant ce fichier.
CREATE EXTENSION IF NOT EXISTS pgcrypto;

CREATE TYPE account_status AS ENUM ('active', 'disabled', 'pending');

CREATE TABLE roles (
    id          UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name        TEXT NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE TABLE users (
    id                UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    email             TEXT NOT NULL UNIQUE,
    password_hash     TEXT NOT NULL,
    display_name      TEXT,
    account_status    account_status NOT NULL DEFAULT 'active',
    email_verified_at TIMESTAMPTZ,
    token_version     INTEGER NOT NULL DEFAULT 0,
    failed_login_count INTEGER NOT NULL DEFAULT 0,
    locked_until      TIMESTAMPTZ,
    created_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at        TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    deleted_at        TIMESTAMPTZ
);

CREATE TABLE user_roles (
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role_id UUID NOT NULL REFERENCES roles(id) ON DELETE CASCADE,
    PRIMARY KEY (user_id, role_id)
);

CREATE TABLE refresh_tokens (
    id            UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    user_id       UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    session_id    UUID NOT NULL DEFAULT gen_random_uuid(),
    token_hash    TEXT NOT NULL UNIQUE,
    device_name   TEXT,
    device_os     TEXT,
    ip_address    TEXT,
    user_agent    TEXT,
    last_seen_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at    TIMESTAMPTZ NOT NULL,
    revoked_at    TIMESTAMPTZ,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_refresh_tokens_user_id ON refresh_tokens(user_id);
CREATE INDEX idx_refresh_tokens_session_id ON refresh_tokens(session_id);
CREATE INDEX idx_users_email ON users(email);

INSERT INTO roles (name) VALUES ('user'), ('admin');

-- Admin bootstrap : id fixe, email/mot de passe depuis GUC (hash bcrypt cost 14).
INSERT INTO users (id, email, password_hash, display_name, account_status, email_verified_at)
VALUES (
    '20000000-0000-4000-8000-000000000001',
    current_setting('comores.admin_email'),
    crypt(current_setting('comores.admin_password'), gen_salt('bf', 14)),
    'Admin',
    'active',
    NOW()
)
ON CONFLICT (id) DO NOTHING;

INSERT INTO user_roles (user_id, role_id)
SELECT u.id, r.id
FROM users u
CROSS JOIN roles r
WHERE u.id = '20000000-0000-4000-8000-000000000001' AND r.name = 'admin'
ON CONFLICT DO NOTHING;
