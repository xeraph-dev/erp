INSERT INTO erp.refresh_tokens (user_id, family_id, token_hash, expires_at)
VALUES ($1, $2, $3, $4);
