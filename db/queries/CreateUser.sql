INSERT INTO erp.users (username, password_hash, email)
VALUES ($1, $2, $3)
RETURNING id, username, password_hash, email, first_name, last_name;
