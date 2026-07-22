SELECT id, username, email, password_hash, first_name, last_name
FROM erp.users
WHERE username = $1;
