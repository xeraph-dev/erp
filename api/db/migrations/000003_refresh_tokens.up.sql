CREATE TABLE IF NOT EXISTS auth.refresh_tokens (
    token_hash  TEXT PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    family_id   UUID NOT NULL,
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at  TIMESTAMP WITH TIME ZONE
);
