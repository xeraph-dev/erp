UPDATE auth.refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP
WHERE revoked_at IS NULL
  AND family_id = @family_id
RETURNING token_hash, user_id, family_id, expires_at, revoked_at;
