SELECT user_id, family_id, token_hash, expires_at, revoked_at
FROM erp.refresh_tokens
WHERE token_hash = $1
LIMIT 1;
