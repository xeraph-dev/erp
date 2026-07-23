UPDATE erp.refresh_tokens
SET revoked_at = CURRENT_TIMESTAMP
WHERE family_id = $1
RETURNING user_id, family_id, token_hash, expires_at, revoked_at;
