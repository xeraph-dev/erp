DELETE FROM auth.users
WHERE id = @id
RETURNING id, username, email password_hash, first_name, last_name;
