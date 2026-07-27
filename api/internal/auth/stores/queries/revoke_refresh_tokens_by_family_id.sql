UPDATE auth.refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP
WHERE family_id = @family_id
RETURNING token_hash, user_id, family_id, expires_at, revoked_at;
