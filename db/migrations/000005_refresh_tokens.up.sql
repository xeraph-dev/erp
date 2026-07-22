CREATE TABLE IF NOT EXISTS erp.refresh_tokens (
    user_id UUID NOT NULL
    REFERENCES erp.users (id) ON UPDATE CASCADE ON DELETE CASCADE,
    token_hash TEXT NOT NULL,
    family_id UUID NOT NULL,
    expires_at TIMESTAMP WITH TIME ZONE NOT NULL,

    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,

    PRIMARY KEY (user_id, token_hash)
);
