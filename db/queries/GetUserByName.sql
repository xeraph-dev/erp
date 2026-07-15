SELECT id, username, password_hash, email, first_name, last_name
FROM erp.users
WHERE username = $1;
