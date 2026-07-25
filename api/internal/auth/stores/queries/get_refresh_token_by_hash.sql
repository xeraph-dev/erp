SELECT token_hash, user_id, family_id, expires_at, revoked_at
FROM auth.refresh_tokens
WHERE token_hash = @token_hash;
