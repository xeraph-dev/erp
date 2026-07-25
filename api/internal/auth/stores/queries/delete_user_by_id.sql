DELETE FROM auth.users
WHERE id = $1
RETURNING id, username, email password_hash, first_name, last_name;
