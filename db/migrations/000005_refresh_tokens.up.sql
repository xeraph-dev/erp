CREATE TABLE IF NOT EXISTS erp.refresh_tokens (
    token_hash TEXT PRIMARY KEY,
    user_id UUID NOT NULL
    REFERENCES erp.users_table (id) ON UPDATE CASCADE ON DELETE CASCADE,
    family_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,
    revoked_at TIMESTAMP WITH TIME ZONE
);
