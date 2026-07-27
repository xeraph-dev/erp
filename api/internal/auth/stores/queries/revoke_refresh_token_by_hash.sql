UPDATE auth.refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP
WHERE token_hash = @token_hash
RETURNING token_hash, user_id, family_id, expires_at, revoked_at;
