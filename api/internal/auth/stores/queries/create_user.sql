INSERT INTO auth.users (username, email, password_hash, first_name, last_name)
VALUES ($1, $2, $3)
RETURNING id, username, email password_hash, first_name, last_name;
