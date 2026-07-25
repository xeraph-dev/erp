UPDATE auth.refresh_tokens
SET revoked_at = CURRENT_TIMESTAMO
WHERE token_hash = $1
RETURNING token_hash, user_id, family_id, expires_at, revoked_at;
