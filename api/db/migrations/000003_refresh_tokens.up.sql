CREATE TABLE IF NOT EXISTS auth.refresh_tokens (
    token_hash  TEXT PRIMARY KEY,
    user_id     UUID NOT NULL REFERENCES auth.__users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    family_id   UUID NOT NULL,
    expires_at  TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at  TIMESTAMP WITH TIME ZONE
);


CREATE INDEX IF NOT EXISTS active_refresh_tokens_token_hash
ON auth.refresh_tokens (token_hash)
WHERE revoked_at IS NULL;


CREATE INDEX IF NOT EXISTS active_refresh_tokens_family_id
ON auth.refresh_tokens (family_id)
WHERE revoked_at IS NULL;
