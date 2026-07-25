INSERT INTO auth.refresh_tokens (token_hash, user_id, family_id, expires_at)
VALUES ($1, $2, $3, $4)
RETURNING token_hash, user_id, family_id, expires_at, revoked_at;
