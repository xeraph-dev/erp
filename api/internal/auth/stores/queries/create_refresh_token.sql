INSERT INTO auth.refresh_tokens (token_hash, user_id, family_id, expires_at)
VALUES (@token_hash, @user_id, @family_id, @expires_at)
RETURNING token_hash, user_id, family_id, expires_at, revoked_at;
