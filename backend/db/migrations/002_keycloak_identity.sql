ALTER TABLE users ADD COLUMN IF NOT EXISTS keycloak_issuer TEXT;
ALTER TABLE users ADD COLUMN IF NOT EXISTS keycloak_subject TEXT;

CREATE UNIQUE INDEX IF NOT EXISTS users_keycloak_identity_idx
    ON users (keycloak_issuer, keycloak_subject)
    WHERE keycloak_issuer IS NOT NULL AND keycloak_subject IS NOT NULL;
